package main

type CuraDefinition struct {
	Name     string                                     `json:"name"`
	Version  float64                                    `json:"version"`
	Metadata CuraDefinitionMetadata                     `json:"metadata"`
	Settings map[string]*CuraDefinitionSettingsCategory `json:"settings"`
}

type CuraDefinitionMetadata struct {
	Type                      string            `json:"type"`
	Author                    string            `json:"author"`
	Category                  string            `json:"category"`
	Manufacturer              string            `json:"manufacturer"`
	SettingVersion            int64             `json:"setting_version"`
	FileFormats               string            `json:"file_formats"`
	Visible                   bool              `json:"visible"`
	HasMaterials              bool              `json:"has_materials"`
	PreferredMaterial         string            `json:"preferred_material"`
	PreferredQualityType      string            `json:"preferred_quality_type"`
	MachineExtruderTrains     map[string]string `json:"machine_extruder_trains"`
	SupportsUSBConnection     bool              `json:"supports_usb_connection"`
	SupportsNetworkConnection bool              `json:"supports_network_connection"`
}

type CuraDefinitionSettingsCategory struct {
	Label       string                            `json:"label"`
	Type        string                            `json:"type"`
	Icon        string                            `json:"icon"`
	Description string                            `json:"description"`
	Children    map[string]*CuraDefinitionSetting `json:"children"`
}

type CuraDefinitionSetting struct {
	Label       string `json:"label"`
	Description string `json:"description"`

	// °C
	// °C/s
	// °
	// mm
	// mm/s
	// mm³
	// s
	// 1/mm
	// %
	// mm²
	// mm/s²
	Unit string `json:"unit"`

	Type                 string                            `json:"type"`
	MinimumValue         string                            `json:"minimum_value"`
	MinimumValueWarning  string                            `json:"minimum_value_warning"`
	MaximumValueWarning  string                            `json:"maximum_value_warning"`
	MaximumValue         string                            `json:"maximum_value"`
	DefaultValue         interface{}                       `json:"default_value"`
	SettablePerMesh      bool                              `json:"settable_per_mesh"`
	Enabled              interface{}                       `json:"enabled"` // string with expression or bool
	Value                *string                           `json:"value"`
	SettablePerExtruder  *bool                             `json:"settable_per_extruder"`
	Children             map[string]*CuraDefinitionSetting `json:"children"`
	Resolve              *string                           `json:"resolve"`
	SettablePerMeshgroup *bool                             `json:"settable_per_meshgroup"`
	LimitToExtruder      *string                           `json:"limit_to_extruder"`
}
