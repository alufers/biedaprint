<!--
FieldsList renders field editComponents based on a list of field descriptors. It provides them with the data from the settings and handles the input events.
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
import settingsSchema from "../../assets/settings-schema.json";
import SettingsFieldDescriptor from "../../types/SettingsFieldDescriptor";
import SettingsPageDescriptor from "../../types/SettingsPageDescriptor";
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

type Settings = any

@Component({
  components: {
    LoaderGuard
  }
})
export default class FieldsList extends mixins(LoadableMixin) {
  @ApolloQuery<GetSettingsQuery>({
    query: getSettings
  })
  settings: Settings = null;

  @Prop({ type: Array })
  fields: SettingsFieldDescriptor[];

  fieldComponent(field: SettingsFieldDescriptor) {
    switch (field.editComponent) {
      case "EnumSelect":
        return EnumSelect;
      case "TemperaturePresetsTable":
        return TemperaturePresetsTable;
      case "TextField":
      case "IntField":
      default:
        return GenericInputField;
    }
  }

  fieldValue(field: SettingsFieldDescriptor) {
    return (this.settings as any)[field.key];
  }

  onFieldInputEvent(field: SettingsFieldDescriptor, newValue: any) {
    let cache = this.$apollo.provider.defaultClient.cache;
    let { settings } = cache.readQuery<GetSettingsQuery>({
      query: getSettings
    });

    (settings as any)[field.key] = newValue;
    cache.writeQuery<GetSettingsQuery>({
      query: getSettings,
      data: { settings }
    });
  }
}
</script>

<style>
</style>
