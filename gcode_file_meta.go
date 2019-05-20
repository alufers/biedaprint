package biedaprint

import (
	"encoding/gob"
	"os"
)

func init() {
	gob.Register(GcodeFileMeta{})
	gob.Register(GcodeLayerIndex{})
}

func loadGcodeFileMeta(path string) (*GcodeFileMeta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	deco := gob.NewDecoder(f)
	var out *GcodeFileMeta
	err = deco.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
