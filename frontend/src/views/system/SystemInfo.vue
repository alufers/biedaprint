<template>
  <div>
    <h2 class="title">System information</h2>
    <button class="button is-primary" @click="loadData">
      <i class="fas fa-sync"></i>
    </button>
    <div v-if="systemInfo">
      <table class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(value, key) in systemInfo" :key="key">
            <td>{{key}}</td>
            <td>{{value}}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import gql from "graphql-tag";
import getSystemInformation from "../../../../queries/getSystemInformation.graphql";
import { GetSystemInformationQuery } from "../../graphql-models-gen";

@Component({})
export default class SystemInfo extends mixins(LoadableMixin) {
  systemInfo: any = null;
  created() {
    this.loadData();
  }
  loadData() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<GetSystemInformationQuery>({
        query: getSystemInformation,
        fetchPolicy: "network-only"
      });
      this.systemInfo = data.systemInformation;
    });
  }
}
</script>

<style>
</style>
