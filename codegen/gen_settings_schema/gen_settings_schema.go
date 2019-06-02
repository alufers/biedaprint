package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func loadSourceFiles() (sources []*ast.Source, err error) {
	filenames, err := filepath.Glob("./graphql/schema/*.graphql")
	if err != nil {
		err = errors.Wrap(err, "failed to glob graphql schema files")
		return
	}
	sources = make([]*ast.Source, 0, len(filenames))
	for _, fname := range filenames {
		rawData, readFileErr := ioutil.ReadFile(fname)
		if readFileErr != nil {
			readFileErr = errors.Wrapf(readFileErr, "failed to read graphql schema file %s", fname)
			return
		}
		strData := string(rawData)
		src := &ast.Source{
			Name:  fname,
			Input: strData,
		}
		sources = append(sources, src)

	}
	return
}

func processGraphqlType(def *ast.Definition) {

	for _, f := range def.Fields {
		if f.Directives.ForName("settingsField") != nil {
			if f.Type.String() == "Int!" || f.Type.String() == "String!" || f.Type.String() == "Float!" || f.Type.String() == "Boolean!" {
				pageArg := f.Directives.ForName("settingsField").Arguments.ForName("page")
				if pageArg != nil {
					fmt.Printf("%#v\n", pageArg.Value)
				}
			}
		}
	}
}

func processPages(schema *ast.Schema) []*SettingsPage {
	pagesEnum := schema.Types["SettingsPage"]

	if pagesEnum == nil || pagesEnum.Kind != ast.Enum {
		panic("SettingsPage must not be nil and must be an enum")
	}

	pages := []*SettingsPage{}

	for _, enumValue := range pagesEnum.EnumValues {
		if directive := enumValue.Directives.ForName("settingsPageDesc"); directive != nil {
			pages = append(pages, &SettingsPage{
				EnumName:    enumValue.Name,
				Name:        directive.Arguments.ForName("name").Value.Raw,
				Description: directive.Arguments.ForName("description").Value.Raw,
			})
		}
	}

	return pages
}

func main() {
	sources, err := loadSourceFiles()
	if err != nil {
		panic(err)
	}
	schema, parseError := gqlparser.LoadSchema(sources...)
	if parseError != nil {
		log.Fatalf("failed to parse schema: %v", parseError.Message)
	}
	ps := processPages(schema)

	setSchema := &SettingsSchema{
		Pages: ps,
	}

	data, _ := json.MarshalIndent(setSchema, "", "  ")

	err = ioutil.WriteFile("frontend/src/assets/settings-schema.json", data, 0666)
	if err != nil {
		panic(err)
	}

}
