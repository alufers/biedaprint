package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
	"github.com/pkg/errors"
)

func (r *queryResolver) AvailableUpdates(ctx context.Context) ([]*AvailableUpdate, error) {
	resp, err := http.Get(AppReleasesEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get releases from github")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("github returned status code %v when requesting releases", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	releasesData := []*struct {
		TagName   string `json:"tag_name"`
		Name      string `json:"name"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
		Assets    []*struct {
			Name               string  `json:"name"`
			BrowserDownloadURL string  `json:"browser_download_url"`
			Size               float64 `json:"size"`
		} `json:"assets"`
	}{}
	err = decoder.Decode(&releasesData)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse releases json")
	}

	updates := []*AvailableUpdate{}
	for _, r := range releasesData {
		if AppVersion != "development" {
			currentVersion, err := semver.Make(AppVersion)
			if err == nil {
				releaseVersion, err := semver.Make(r.TagName)
				if err == nil {
					if releaseVersion.LT(currentVersion) {
						continue // skip old versions
					}
				}
			}
		}
		upd := &AvailableUpdate{
			TagName:   r.TagName,
			Title:     r.Name,
			Body:      r.Body,
			CreatedAt: r.CreatedAt,
		}
		for _, asset := range r.Assets {
			if asset.Name == AppReleaseExecutableName {
				upd.ExecutableURL = asset.BrowserDownloadURL
				upd.Size = byteCountBinary(int64(asset.Size))
				updates = append(updates, upd)
				break
			}
		}
	}
	return updates, nil

}

func (r *mutationResolver) DownloadUpdate(ctx context.Context, tagName string) (*bool, error) {
	metaResp, err := http.Get(fmt.Sprintf(AppSingleReleaseEndpoint, tagName))

	if err != nil {
		return nil, errors.Wrap(err, "failed to get release from github")
	}
	defer metaResp.Body.Close()
	if metaResp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("github returned status code %v when requesting release", metaResp.StatusCode)
	}

	decoder := json.NewDecoder(metaResp.Body)
	releaseData := &struct {
		Assets []*struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}{}
	err = decoder.Decode(&releaseData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse releases json")
	}

	var downloadURL string
	for _, asset := range releaseData.Assets {
		if asset.Name == AppReleaseExecutableName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return nil, errors.Wrap(err, "failed to find suitable executable in the release")
	}

	dataPath := r.App.GetSettings().DataPath
	out, err := os.Create(filepath.Join(dataPath, fmt.Sprintf("biedaprint-update-%v", tagName)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create update temporary file")
	}
	defer out.Close()
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request the update file")
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save the update file")
	}
	return nil, nil
}

func (r *mutationResolver) PerformUpdate(ctx context.Context, tagName string) (*bool, error) {
	dataPath := r.App.GetSettings().DataPath
	updateFilePath := filepath.Join(dataPath, fmt.Sprintf("biedaprint-update-%v", tagName))
	if _, err := os.Stat(updateFilePath); os.IsNotExist(err) {
		return nil, errors.New("update file for this tagName does not exist")
	}
	f, err := os.Open(updateFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the update file")
	}
	err = update.Apply(f, update.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply the update")
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		bin, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to get current executable: %v", err))
		}
		err = syscall.Exec(bin, append([]string{bin}, os.Args[1:]...), os.Environ())
		if err != nil {
			panic(fmt.Sprintf("cannot restart after update: %v", err))
		}
	}()
	return nil, nil
}
