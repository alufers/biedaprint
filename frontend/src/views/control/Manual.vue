<template>
  <div>
    <h3 class="title">Manual control</h3>
    <label class="label">Homing</label>
    <div class="field has-addons">
      <p class="control">
        <!-- home all axes -->
        <button
          class="button is-dark is-centered is-outlined"
          @click="sendGCODE('G28')"
          title="Home all axes"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <!-- home X -->
        <button
          class="button is-danger is-centered is-outlined"
          @click="sendGCODE('G28 X')"
          title="Home X"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <button
          class="button is-success is-centered is-outlined"
          @click="sendGCODE('G28 Y')"
          title="Home Y"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <!-- home Z -->
        <button
          class="button is-primary is-centered is-outlined"
          @click="sendGCODE('G28 Z')"
          title="Home Z"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import { sendGcode } from "../../../../graphql/queries/sendGcode.graphql";
import {
  SendGcodeMutation,
  SendGcodeMutationVariables
} from "../../graphql-models-gen";

@Component({})
export default class Manual extends mixins(LoadableMixin) {
  sendGCODE(gcode: string) {
    this.withLoader(async () => {
      await this.$apollo.mutate<SendGcodeMutation>({
        mutation: sendGcode,
        variables: <SendGcodeMutationVariables>{
          cmd: gcode
        }
      });
    });
  }
}
</script>
<style scoped>
td {
  text-align: center;
}
button {
  width: 50px;
  height: 50px;
}
</style>
