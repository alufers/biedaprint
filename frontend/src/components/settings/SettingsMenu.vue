<!-- 
SettingsMenu geneartes a settings menu from the settings schema.
-->
<template>
  <aside class="menu">
    <div class="field">
      <p class="control has-icons-right">
        <input
          class="input"
          type="text"
          placeholder="Search settings"
          :value="searchQuery"
          @input="onSearched"
        />

        <span class="icon is-small is-right">
          <i class="fas fa-search"></i>
        </span>
      </p>
    </div>
    <p class="menu-label">Pages</p>
    <ul class="menu-list">
      <li v-for="page in pages" :key="page.enumName">
        <router-link :to="linkToPage(page)" :title="page.description">{{page.title}}</router-link>
      </li>
    </ul>
  </aside>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
// @ts-ignore
import settingsSchema from "../../../../graphql/schema/settings.schema.json";
import { JsonSchema, pages } from "../../util/settingsSchema";

@Component({})
export default class SettingsMenu extends Vue {
  settingsSchema = settingsSchema;
  pages = pages;
  linkToPage(page: JsonSchema) {
    return `/system/settings/${page.urlParamName}`;
  }
  get searchQuery() {
    let matched = this.$route.matched[this.$route.matched.length - 1];
    if (!matched || matched.name !== "settings-search") {
      return "";
    }
    return this.$route.query.query;
  }
  onSearched(ev: Event) {
    this.$router.push(
      "/system/settings/search?query=" +
        encodeURIComponent((ev.target as HTMLInputElement).value)
    );
  }
}
</script>

<style>
</style>
