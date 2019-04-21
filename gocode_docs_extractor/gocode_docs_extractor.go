package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gernest/front"
)

func processJSON(d interface{}) interface{} {
	switch typed := d.(type) {
	case map[interface{}]interface{}:
		newMap := map[string]interface{}{}
		for k, v := range typed {
			if strKey, ok := k.(string); ok {
				newMap[strKey] = processJSON(v)
			}
		}
		return newMap
	case map[string]interface{}:
		newMap := map[string]interface{}{}
		for k, v := range typed {
			newMap[k] = processJSON(v)
		}
		return newMap
	case []interface{}:
		newList := []interface{}{}
		for _, e := range typed {
			newList = append(newList, processJSON(e))
		}
		return newList
	default:

		return d
	}
}

func main() {
	files, err := filepath.Glob("_gcode/*.md")
	if err != nil {
		log.Fatal(err)
		return
	}
	outData := map[string]interface{}{}
	for _, fp := range files {
		f, err := os.Open(fp)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		m := front.NewMatter()
		m.Handle("---", front.YAMLHandler)
		frnt, _, err := m.Parse(f)
		if err != nil {
			log.Fatal(err)
			return
		}
		if codes, ok := frnt["codes"].([]interface{}); ok {
			for _, code := range codes {
				outData[code.(string)] = frnt
				break // we anly really want one code
			}
		}
		extension := filepath.Ext(fp)
		base := filepath.Base(fp)
		frnt["base"] = base[0 : len(base)-len(extension)]

	}
	data, err := json.Marshal(processJSON(outData))
	if err != nil {
		log.Fatal(err)
		return
	}
	err = ioutil.WriteFile("frontend/src/assets/gcode-docs.json", data, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
}
