<template>
  <aside class="menu">
    <p class="menu-label">Settings</p>
    <ul class="menu-list">
      <li v-for="page in settingsSchema.pages" :key="page.enumName">
        <a>{{page.name}}</a>
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
}
</script>

<style>
</style>
