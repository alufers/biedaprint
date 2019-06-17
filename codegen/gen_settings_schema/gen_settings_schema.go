package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/iancoleman/strcase"
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

func processPages(schema *ast.Schema, result *SettingsSchema) {
	pagesEnum := schema.Types["SettingsPage"]

	if pagesEnum == nil || pagesEnum.Kind != ast.Enum {
		panic("SettingsPage must not be nil and must be an enum")
	}

	pages := []*SettingsPage{}

	for _, enumValue := range pagesEnum.EnumValues {
		if directive := enumValue.Directives.ForName("settingsPageDesc"); directive != nil {
			pages = append(pages, &SettingsPage{
				ParamName:   strcase.ToKebab(enumValue.Name),
				EnumName:    enumValue.Name,
				Name:        directive.Arguments.ForName("name").Value.Raw,
				Description: directive.Arguments.ForName("description").Value.Raw,
			})
		}
	}

	result.Pages = pages
}

func isTypeEnum(t *ast.Definition) bool {

	if t == nil || t.Kind != ast.Enum {
		return false
	}
	return true
}

func isTypeObject(t *ast.Definition) bool {

	if t == nil || t.Kind != ast.Object {
		return false
	}
	return true
}

func getStringArgument(l ast.ArgumentList, name string) string {
	arg := l.ForName(name)
	if arg == nil {
		return ""
	}
	return arg.Value.Raw
}

func processEnumType(schema *ast.Schema, result *SettingsSchema, enumName string) {
	enumType := schema.Types[enumName]

	if !isTypeEnum(enumType) {
		return
	}

	for _, se := range result.Enums {
		if se.Name == enumName {
			return
		}
	}
	settingsEnum := &SettingsEnum{
		Name:   enumName,
		Values: []*SettingsEnumValue{},
	}

	for _, val := range enumType.EnumValues {
		if valDirective := val.Directives.ForName("enumValueDesc"); valDirective != nil {
			settingsEnum.Values = append(settingsEnum.Values, &SettingsEnumValue{
				Value: val.Name,
				Label: valDirective.Arguments.ForName("label").Value.Raw,
			})
		}
	}

	result.Enums = append(result.Enums, settingsEnum)
}

func processObjectType(schema *ast.Schema, result *SettingsSchema, keyPrefix string, typeName string) {
	settingsType := schema.Types[typeName]
	if settingsType == nil || settingsType.Kind != ast.Object {
		panic("type must not be nil and must be an object")
	}

	for _, field := range settingsType.Fields {
		if directive := field.Directives.ForName("settingsField"); directive != nil {
			typeName := field.Type.String()

			var description string
			pageNameArg := directive.Arguments.ForName("page")
			if pageNameArg == nil {
				panic("page name missing")
			}
			if descriptionArg := directive.Arguments.ForName("description"); descriptionArg != nil {
				description = descriptionArg.Value.Raw
			}
			var editComponent string
			if editComponentArg := directive.Arguments.ForName("editComponent"); editComponentArg != nil {
				editComponent = editComponentArg.Value.Raw
			} else {
				switch typeName {
				case "Int!":
					editComponent = "IntField"
				case "Float!":
					editComponent = "FloatField"
				case "String!":
					editComponent = "TextField"
				case "Boolean!":
					editComponent = "CheckboxField"
				}
				if isTypeEnum(schema.Types[field.Type.Name()]) {
					editComponent = "EnumSelect"
				}

				if isTypeObject(schema.Types[field.Type.Name()]) {
					processObjectType(schema, result, keyPrefix+field.Name+".", field.Type.Name())
					continue
				}

			}
			processEnumType(schema, result, field.Type.Name())
			result.Fields = append(result.Fields, &SettingsField{
				Key:           keyPrefix + field.Name,
				ParamName:     strcase.ToKebab(field.Name),
				PageEnumName:  pageNameArg.Value.Raw,
				GraphQLType:   typeName,
				Label:         directive.Arguments.ForName("label").Value.Raw,
				Description:   description,
				EditComponent: editComponent,
				EnumTypeName:  field.Type.Name(),
				Unit:          getStringArgument(directive.Arguments, "unit"),
			})
		}
	}
}

func processFields(schema *ast.Schema, result *SettingsSchema) {
	processObjectType(schema, result, "", "Settings")
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

	setSchema := &SettingsSchema{
		Fields: []*SettingsField{},
	}
	processFields(schema, setSchema)
	processPages(schema, setSchema)
	data, _ := json.MarshalIndent(setSchema, "", "  ")

	err = ioutil.WriteFile("frontend/src/assets/settings-schema.json", data, 0666)
	if err != nil {
		panic(err)
	}

}
