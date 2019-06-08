<template>
  <div>
    <h3 class="subtitle">{{pageData.name}}</h3>
    <p class="has-text-grey-light">{{pageData.description}}</p>
    <br>

    <LoaderGuard>
      <component
        v-for="field of fields"
        :key="field.paramName"
        :is="fieldComponent(field)"
        :fieldDescriptor="field"
        :value="fieldValue(field)"
        @input="onFieldInputEvent(field, $event)"
      />
      <br>
      <br>
      <button class="button is-primary" @click="saveSettings">Save</button>
    </LoaderGuard>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import settingsSchema from "../../../assets/settings-schema.json";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import SettingsPageDescriptor from "../../../types/SettingsPageDescriptor";
import GenericInputField from "../../../components/settings/fields/GenericInputField.vue";
import EnumSelect from "../../../components/settings/fields/EnumSelect.vue";
import TemperaturePresetsTable from "../../../components/settings/fields/TemperaturePresetsTable.vue";
import {
  Settings,
  GetSettingsQuery,
  UpdateSettingsMutationVariables
} from "../../../graphql-models-gen";
import LoadableMixin from "../../../LoadableMixin";
import ApolloQuery from "../../../decorators/ApolloQuery";
import { getSettings } from "../../../../../graphql/queries/getSettings.graphql";
import LoaderGuard from "../../../components/LoaderGuard.vue";
import { updateSettings } from "../../../../../graphql/queries/updateSettings.graphql";
import omitTypename from "../../../util/omitTypename";

@Component({
  components: {
    LoaderGuard
  }
})
export default class SettingsPage extends mixins(LoadableMixin) {
  @ApolloQuery<GetSettingsQuery>({
    query: getSettings
  })
  settings: Settings = null;

  get fields(): SettingsFieldDescriptor[] {
    return settingsSchema.fields.filter(
      (f: SettingsFieldDescriptor) => f.pageEnumName === this.pageData.enumName
    );
  }

  get pageData(): SettingsPageDescriptor {
    return settingsSchema.pages.find(
      p => p.paramName === this.$route.params.pageName
    );
  }

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

  saveSettings() {
    this.withLoader(async () => {
      let cache = this.$apollo.provider.defaultClient.cache;
      let { settings } = cache.readQuery<GetSettingsQuery>({
        query: getSettings
      });

      await this.$apollo.mutate({
        mutation: updateSettings,
        variables: <UpdateSettingsMutationVariables>{
          newSettings: omitTypename(settings)
        }
      });
    });
  }
}
</script>

<style>
</style>
