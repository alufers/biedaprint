package main

type SettingsSchema struct {
	Pages []*SettingsPage `json:"pages"`
}

type SettingsPage struct {
	EnumName    string           `json:"enumName"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Fields      []*SettingsField `json:"fields"`
}

type SettingsField struct {
	Key           string `json:"key"`
	GraphQLType   string `json:"graphQlType"`
	EditComponent string `json:"editComponent"`
	Label         string `json:"label"`
	Description   string `json:"description"`
}
