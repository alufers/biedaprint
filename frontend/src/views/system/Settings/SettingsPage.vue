<template>
  <div v-if="this.pageData">
    <h3 class="subtitle">{{pageData.title}}</h3>
    <p class="has-text-grey-light">{{pageData.description}}</p>
    <br />

    <LoaderGuard>
      <FieldsList :fields="fields" />
      <button class="button is-primary" @click="saveSettings">Save</button>
    </LoaderGuard>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import SettingsPageDescriptor from "../../../types/SettingsPageDescriptor";
import {
  // Settings,
  GetSettingsQuery,
  UpdateSettingsMutationVariables
} from "../../../graphql-models-gen";
import LoadableMixin from "../../../LoadableMixin";
import ApolloQuery from "../../../decorators/ApolloQuery";
import { getSettings } from "../../../../../graphql/queries/getSettings.graphql";
import LoaderGuard from "../../../components/LoaderGuard.vue";
import { updateSettings } from "../../../../../graphql/queries/updateSettings.graphql";
import omitTypename from "../../../util/omitTypename";
import FieldsList from "../../../components/settings/FieldsList.vue";
import {
  pages,
  getNormalizedFields,
  JsonSchema
} from "../../../util/settingsSchema";

type Settings = any;

@Component({
  components: {
    LoaderGuard,
    FieldsList
  }
})
export default class SettingsPage extends mixins(LoadableMixin) {
  get fields(): JsonSchema[] {
    return getNormalizedFields(this.pageData, this.pageData.fullPath);
  }

  get pageData() {
    return pages.find(p => p.urlParamName === this.$route.params.pageName);
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
