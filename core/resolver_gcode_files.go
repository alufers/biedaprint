package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

func (r *mutationResolver) UploadGcode(ctx context.Context, file graphql.Upload) (*GcodeFileMeta, error) {
	dataPath := r.App.GetDataPath()
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
	err = r.App.GcodeFileMetaRepositoryService.Save(meta)
	if err != nil {
		return nil, fmt.Errorf("failed to save gcode file meta: %w", err)
	}
	return meta, nil
}

func (r *mutationResolver) DeleteGcodeFile(ctx context.Context, id int) (*bool, error) {
	file, err := r.App.GcodeFileMetaRepositoryService.GetOneByID(id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, fmt.Errorf("gcode file meta not found")
	}
	gcodeName := filepath.Join(r.App.GetDataPath(), "gcode_files/", file.GcodeFileName)

	err = os.Remove(gcodeName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove gcode file")
	}
	err = r.App.GcodeFileMetaRepositoryService.Delete(file)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *queryResolver) GcodeFileMetas(ctx context.Context) ([]*GcodeFileMeta, error) {
	return r.App.GcodeFileMetaRepositoryService.GetAll()
}
