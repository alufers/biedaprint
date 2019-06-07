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
import TextField from "../../../components/settings/fields/TextField.vue";
import IntField from "../../../components/settings/fields/IntField.vue";
import {
  Settings,
  GetSettingsQuery,
  UpdateSettingsMutationVariables
} from "../../../graphql-models-gen";
import LoadableMixin from "../../../LoadableMixin";
import ApolloQuery from "../../../ApolloQuery";
import { getSettings } from "../../../../../graphql/queries/getSettings.graphql";
import LoaderGuard from "../../../components/LoaderGuard.vue";
import { updateSettings } from "../../../../../graphql/queries/updateSettings.graphql";
import omitTypename from "../../../omitTypename";

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
      case "IntField":
        return IntField;
      case "TextField":
      default:
        return TextField;
    }
  }

  fieldValue(field: SettingsFieldDescriptor) {
    return this.settings[field.key];
  }

  onFieldInputEvent(field: SettingsFieldDescriptor, newValue: any) {
    let cache = this.$apollo.provider.defaultClient.cache;
    let { settings } = cache.readQuery<GetSettingsQuery>({
      query: getSettings
    });

    settings[field.key] = newValue;
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
