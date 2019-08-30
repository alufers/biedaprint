<template>
  <LoaderGuard>
    <h2 class="subtitle">Settings</h2>
    <progress class="progress is-large is-primary" max="100" v-if="loading">15%</progress>
    <div class="columns">
      <div class="column is-one-fifth">
        <!-- <SettingsMenu/> -->
      </div>
      <div class="column box">
        <!-- <router-view/> -->
      </div>
    </div>
  </LoaderGuard>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import gql from "graphql-tag";
import getSettingsAndSerialPorts from "../../../../graphql/queries/getSettingsAndSerialPorts.graphql";
import updateSettings from "../../../../graphql/queries/updateSettings.graphql";
import LoaderGuard from "../../components/LoaderGuard.vue";
import SettingsMenu from "../../components/settings/SettingsMenu.vue";

@Component({
  components: {
    LoaderGuard,
    SettingsMenu
  }
})
export default class SettingsPage extends mixins(LoadableMixin) {}
</script>