package main

import "github.com/vektah/gqlparser/ast"

type TypeVisitor interface {
	VisitObjectType(def *ast.Definition)
	VisitEnumType(def *ast.Definition)
}

func TraverseGraphqlSchema(schema *ast.Schema, visitor TypeVisitor) {
	for _, def := range schema.Types {
		if isTypeEnum(def) {
			visitor.VisitEnumType(def)
		}
		if isTypeEnum(def) {
			visitor.VisitEnumType(def)
		}
	}
}
