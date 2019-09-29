<template>
  <div>
    <LoaderGuard>
      <h2 class="title">Updates</h2>
      <button class="button is-primary" @click="loadData">
        <i class="fas fa-sync"></i>
      </button>
      <div v-if="availableUpdates">
        <table class="table">
          <thead>
            <tr>
              <th>Title</th>
              <th>Version</th>
              <th>Size</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(update, key) in availableUpdates" :key="key">
              <td>{{update.title}}</td>
              <td>{{update.tagName}}</td>
              <td>{{update.size}}</td>
              <td>
                <button class="button is-primary">
                  <i class="fas fa-download" @click="downloadAndPerformUpdate(update.tagName)"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </LoaderGuard>
    <div class="modal is-active" v-if="showProgressModal">
      <div class="modal-background"></div>
      <div class="modal-content"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">Installing update...</p>
        </header>
        <section class="modal-card-body">
          <p>{{updateStatus}}</p>
          <progress class="progress" max="100">15%</progress>
        </section>
      </div>
      <button class="modal-close is-large" aria-label="close" @click="gcodeFileToDelete = null"></button>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import gql from "graphql-tag";
import getAvailableUpdates from "../../../../graphql/queries/getAvailableUpdates.graphql";
import performUpdate from "../../../../graphql/queries/performUpdate.graphql";
import downloadUpdate from "../../../../graphql/queries/downloadUpdate.graphql";
import {
  GetAvailableUpdatesQuery,
  AvailableUpdate,
  DownloadUpdateMutationVariables,
  DownloadUpdateMutation,
  PerformUpdateMutation,
  PerformUpdateMutationVariables
} from "../../graphql-models-gen";
import LoaderGuard from "../../components/LoaderGuard.vue";
import { Watch } from "vue-property-decorator";
import sleep from "../../util/sleep";

@Component({
  components: { LoaderGuard }
})
export default class Updates extends mixins(LoadableMixin) {
  availableUpdates: AvailableUpdate[] = null;
  showProgressModal = false;
  updateStatus = "";
  created() {
    this.loadData();
  }
  loadData() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<GetAvailableUpdatesQuery>({
        query: getAvailableUpdates,
        fetchPolicy: "network-only"
      });
      this.availableUpdates = data.availableUpdates;
    });
  }
  async downloadAndPerformUpdate(tagName: string) {
    this.showProgressModal = true;
    await this.withLoader(async () => {
      this.updateStatus = "Downloading update...";
      await this.$apollo.mutate<DownloadUpdateMutation>({
        mutation: downloadUpdate,
        variables: <DownloadUpdateMutationVariables>{
          tagName
        }
      });
      this.updateStatus = "Applying update...";
      await this.$apollo.mutate<PerformUpdateMutation>({
        mutation: performUpdate,
        variables: <PerformUpdateMutationVariables>{
          tagName
        }
      });
      this.updateStatus = "Refreshing...";
      await sleep(5000);
      (<any>document.location) = "/system/system-info?_rnd=" + Math.random();
    });
    this.showProgressModal = false;
  }
}
</script>

<style>
</style>
