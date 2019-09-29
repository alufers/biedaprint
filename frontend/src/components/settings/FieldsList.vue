<!--
FieldsList renders field editComponents based on a list of field schemas. It provides them with the data from the settings and handles the input events.
 -->
<template>
  <LoaderGuard>
    <component
      v-for="field of fields"
      :key="field.paramName"
      :is="fieldComponent(field)"
      :fieldDescriptor="field"
      :value="fieldValue(field)"
      @input="onFieldInputEvent(field, $event)"
    />
  </LoaderGuard>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import GenericInputField from "./fields/GenericInputField.vue";
import EnumSelect from "./fields/EnumSelect.vue";
import TemperaturePresetsTable from "./fields/TemperaturePresetsTable.vue";
import {
  // Settings,
  GetSettingsQuery,
  UpdateSettingsMutationVariables
} from "../../graphql-models-gen";
import LoadableMixin from "../../LoadableMixin";
import ApolloQuery from "../../decorators/ApolloQuery";
import LoaderGuard from "../LoaderGuard.vue";
import { getSettings } from "../../../../graphql/queries/getSettings.graphql";
import { updateSettings } from "../../../../graphql/queries/updateSettings.graphql";
import omitTypename from "../../util/omitTypename";
import { Prop } from "vue-property-decorator";
import { JsonSchema } from "../../util/settingsSchema";
import { get, set } from "lodash-es";

@Component({
  components: {
    LoaderGuard
  }
})
export default class FieldsList extends mixins(LoadableMixin) {
  @ApolloQuery<GetSettingsQuery>({
    query: getSettings
  })
  settings: any = null;

  @Prop({ type: Array })
  fields: JsonSchema[];

  fieldComponent(field: JsonSchema) {
    if (field.enum) {
      return EnumSelect;
    }
    switch (field.settingsField || field.type) {
      case "TemperaturePresetsTable":
        return TemperaturePresetsTable;

      case "string":
      case "integer":
      default:
        return GenericInputField;
    }
  }

  fieldValue(field: JsonSchema) {
    return get(this.settings, field.fullPath);
  }

  onFieldInputEvent(field: JsonSchema, newValue: any) {
    let cache = this.$apollo.provider.defaultClient.cache;
    let { settings } = cache.readQuery<GetSettingsQuery>({
      query: getSettings
    });

    set(settings, field.fullPath, newValue);
    cache.writeQuery<GetSettingsQuery>({
      query: getSettings,
      data: { settings }
    });
  }
}
</script>

<style>
</style>
