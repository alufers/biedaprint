<template>
  <div>
    <h2 class="subtitle">Connect to printer</h2>
    <div v-if="serialStatus === 'disconnected' || serialStatus == 'error'">
      <button
        class="button is-success"
        :class="isLoadingClass"
        @click="connectToSerial"
        v-if="serialStatus === 'disconnected' || serialStatus == 'error'"
      >Connect to printer</button>
    </div>
    <button
      class="button is-danger is-outlined"
      :class="isLoadingClass"
      @click="disconnectFromSerial"
      v-if="serialStatus === 'connected'"
    >Disconnect from printer</button>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import gql from "graphql-tag";
import TrackedValueSubscription from "../decorators/TrackedValueSubscription";
import { TrackedValue } from "../graphql-models-gen";
import LoadableMixin from "../LoadableMixin";

@Component({})
export default class PrinterConnectionWidget extends mixins(LoadableMixin) {
  @TrackedValueSubscription("serialStatus")
  serialStatus: string = "disconnected";

  connectToSerial() {
    this.withLoader(async () => {
      await this.$apollo.mutate({
        mutation: gql`
          mutation {
            connectToSerial(void: null)
          }
        `
      });
    });
  }
  async disconnectFromSerial() {
    this.withLoader(async () => {
      await this.$apollo.mutate({
        mutation: gql`
          mutation {
            disconnectFromSerial(void: null)
          }
        `
      });
    });
  }
}
</script>

<style scoped>
.button {
  width: 100%;
}
</style>