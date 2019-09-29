import settingsSchemaOriginal from "../../../graphql/schema/settings.schema.json";

export interface JsonSchema {
  title?: string;
  type?: string;
  isSettingsPage?: boolean;
  properties?: { [x: string]: JsonSchema };
  urlParamName?: string;
  settingsField?: string | boolean;
  fullPath?: string;
  enum?: Array<any>;
}

export const settingsSchema = settingsSchemaOriginal;

export let pages: JsonSchema[] = [];

const traverseSchema = (schema: JsonSchema, path: string = "") => {
  if (typeof schema !== "object") return;
  if (schema.isSettingsPage) {
    pages.push({ ...schema, fullPath: path });
  }
  if (typeof schema.properties === "object") {
    for (const key of Object.keys(schema.properties)) {
      traverseSchema(schema.properties[key], (path ? path + "." : "") + key);
    }
  }
};

traverseSchema(settingsSchema);

export function getNormalizedFields(
  schema: JsonSchema,
  path: string = ""
): JsonSchema[] {
  let fields: JsonSchema[] = [];
  if (schema.settingsField === false) {
    return fields;
  }
  if (schema.settingsField) {
    fields.push({ ...schema, fullPath: path });
  } else if (
    schema.type === "string" ||
    schema.type === "integer" ||
    schema.enum
  ) {
    fields.push({ ...schema, fullPath: path });
  } else if (typeof schema.properties === "object") {
    for (const key of Object.keys(schema.properties)) {
      fields = [
        ...fields,
        ...getNormalizedFields(
          schema.properties[key],
          (path ? path + "." : "") + key
        )
      ];
    }
  }
  return fields;
}
