# gen_settings_schema

This codegen generates a json file (`frontend/src/assets/settings-schema.json`) which contains information about the settings so that the frontend can generate the settings page.

Additionally it generates the graphql input types for the various settings types and a go file which converts the input types to the concrete ones.