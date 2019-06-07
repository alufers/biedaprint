<template>
  <aside class="menu">
    <div class="field">
      <p class="control has-icons-right">
        <input class="input" type="text" placeholder="Search settings">

        <span class="icon is-small is-right">
          <i class="fas fa-search"></i>
        </span>
      </p>
    </div>
    <p class="menu-label">Pages</p>
    <ul class="menu-list">
      <li v-for="page in settingsSchema.pages" :key="page.enumName">
        <router-link :to="linkToPage(page)" :title="page.description">{{page.name}}</router-link>
      </li>
    </ul>
  </aside>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import settingsSchema from "../../assets/settings-schema.json";

@Component({})
export default class SettingsMenu extends Vue {
  settingsSchema = settingsSchema;

  linkToPage(page: any) {
    return `/system/settings/${page.paramName}`;
  }
}
</script>

<style>
</style>
