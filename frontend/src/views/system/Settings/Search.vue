<template>
  <div>
    <h3 class="subtitle">Settings search</h3>
    <LoaderGuard>
      <HighlightableTextZone :highlights="highlightTokens">
        <FieldsList :fields="fields" />
      </HighlightableTextZone>
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
import LoaderGuard from "../../../components/LoaderGuard.vue";
import HighlightableTextZone from "../../../components/HighlightableTextZone.vue";
import FieldsList from "../../../components/settings/FieldsList.vue";
import Fuse from "fuse.js";
import {
  getNormalizedFields,
  settingsSchema,
  JsonSchema
} from "../../../util/settingsSchema";

type Settings = any;

@Component({
  components: {
    LoaderGuard,
    FieldsList,
    HighlightableTextZone
  }
})
export default class SettingsPage extends mixins(LoadableMixin) {
  fuse: Fuse<JsonSchema> = new Fuse(getNormalizedFields(settingsSchema), {
    keys: ["title", "description"],
    caseSensitive: false,
    tokenize: true,
    treshold: 0.3
  });
  get query() {
    return this.$route.query["query"] as string;
  }
  get fields(): JsonSchema[] {
    return this.fuse.search(this.query, {
      limit: 10
    });
  }
  get highlightTokens() {
    return this.query
      .split(" ")
      .map(q => q.trim())
      .filter(q => !!q);
  }
}
</script>

<style>
</style>
