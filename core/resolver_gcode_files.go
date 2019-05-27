package core

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

func (r *mutationResolver) UploadGcode(ctx context.Context, file graphql.Upload) (*GcodeFileMeta, error) {
	dataPath := r.App.GetSettings().DataPath
	gcodeFilename := RandStringRunes(8) + ".gcode"
	// copy example
	f, err := os.OpenFile(filepath.Join(dataPath, "gcode_files/"+gcodeFilename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gcode file for writing")
	}
	defer f.Close()
	_, err = io.Copy(f, file.File)
	if err != nil {
		return nil, errors.Wrap(err, "failed to copy file from request")
	}
	meta := &GcodeFileMeta{
		OriginalName:  file.Filename,
		UploadDate:    time.Now(),
		GcodeFileName: gcodeFilename,
	}
	err = meta.AnalyzeGcodeFile(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to analyze gcode file")
	}
	err = meta.Save(dataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save gcode file meta")
	}
	return meta, nil
}

func (r *mutationResolver) DeleteGcodeFile(ctx context.Context, gcodeFilename string) (*bool, error) {
	dataPath := r.App.GetSettings().DataPath
	gcodeName := filepath.Join(dataPath, "gcode_files/", gcodeFilename)
	gcodeMetaName := filepath.Join(dataPath, "gcode_files/", gcodeFilename+".meta")
	err := os.Remove(gcodeName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove gcode file")
	}
	err = os.Remove(gcodeMetaName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove gcode file meta")
	}
	return nil, nil
}

func (r *queryResolver) GcodeFileMetas(ctx context.Context) ([]*GcodeFileMeta, error) {
	metas := []*GcodeFileMeta{}
	metafilePaths, _ := filepath.Glob(filepath.Join(r.App.GetSettings().DataPath, "gcode_files/", "*.gcode.meta"))
	for _, fp := range metafilePaths {
		meta, err := loadGcodeFileMeta(fp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read gcode file meta")
		}
		metas = append(metas, meta)

	}
	sort.Slice(metas, func(i int, j int) bool {
		return !metas[i].UploadDate.Before(metas[j].UploadDate)
	})
	return metas, nil
}
