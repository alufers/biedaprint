package main

import "github.com/vektah/gqlparser/ast"

type SettingsSchemaVisitor struct {
	resultSchema *SettingsSchema
}

func NewSettingsSchemaVisitor() *SettingsSchemaVisitor {
	return &SettingsSchemaVisitor{
		resultSchema: &SettingsSchema{
			Pages:  []*SettingsPage{},
			Fields: []*SettingsField{},
			Enums:  []*SettingsEnum{},
		},
	}
}

func (ssv *SettingsSchemaVisitor) VisitObjectType(def *ast.Definition) {

}
func (ssv *SettingsSchemaVisitor) VisitEnumType(def *ast.Definition) {

}
