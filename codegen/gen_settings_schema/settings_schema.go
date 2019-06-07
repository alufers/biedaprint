package main

type SettingsSchema struct {
	Pages  []*SettingsPage  `json:"pages"`
	Fields []*SettingsField `json:"fields"`
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
}
