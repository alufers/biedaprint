package main

type SettingsSchema struct {
	Pages  []*SettingsPage  `json:"pages"`
	Fields []*SettingsField `json:"fields"`
	Enums  []*SettingsEnum  `json:"enums"`
}

type SettingsPage struct {
	EnumName    string `json:"enumName"`
	ParamName   string `json:"paramName"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SettingsField struct {
	Key           string `json:"key"`
	ParamName     string `json:"paramName"`
	PageEnumName  string `json:"pageEnumName"`
	GraphQLType   string `json:"graphQlType"`
	EditComponent string `json:"editComponent"`
	Label         string `json:"label"`
	Description   string `json:"description"`
	EnumTypeName  string `json:"enumTypeName"`
}

type SettingsEnum struct {
	Name   string               `json:"name"`
	Values []*SettingsEnumValue `json:"values"`
}

type SettingsEnumValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
