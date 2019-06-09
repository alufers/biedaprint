<template>
  <div>
    <h3 class="title">Manual control</h3>
    <label class="label">Homing</label>
    <div class="field has-addons">
      <p class="control">
        <!-- home all axes -->
        <button
          :class="isLoadingClass"
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
          :class="isLoadingClass"
          class="button is-danger is-centered is-outlined"
          @click="sendGCODE('G28 X')"
          title="Home X"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <button
          :class="isLoadingClass"
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
          :class="isLoadingClass"
          class="button is-primary is-centered is-outlined"
          @click="sendGCODE('G28 Z')"
          title="Home Z"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
    </div>
    <label class="label">Movement</label>
    <table>
      <tbody>
        <tr>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-danger is-centered"
              @click="moveXPositive"
              title
            >X+</button>
          </td>
          <td></td>
        </tr>
        <tr>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-success is-centered"
              @click="moveYNegative"
              title
            >Y-</button>
          </td>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-success is-centered"
              @click="moveYPositive"
              title
            >Y+</button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-danger is-centered"
              @click="moveXNegative"
              title
            >X-</button>
          </td>
          <td></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import { sendGcode } from "../../../../graphql/queries/sendGcode.graphql";
import {
  SendGcodeMutation,
  SendGcodeMutationVariables,
  PerformManualMovementMutation,
  PerformManualMovementMutationVariables,
  ManualMovementPositionVector
} from "../../graphql-models-gen";
import { performManualMovement } from "../../../../graphql/queries/performManualMovement.graphql";

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
  performManualMovement(vec: Partial<ManualMovementPositionVector>) {
    vec = {
      X: vec.X || 0,
      Y: vec.Y || 0,
      Z: vec.Z || 0,
      E: vec.E || 0
    };
    this.withLoader(async () => {
      await this.$apollo.mutate<PerformManualMovementMutation>({
        mutation: performManualMovement,
        variables: <PerformManualMovementMutationVariables>{
          vec
        }
      });
    });
  }
  moveXPositive() {
    this.performManualMovement({ X: 10 });
  }
  moveXNegative() {
    this.performManualMovement({ X: -10 });
  }
  moveYNegative() {
    this.performManualMovement({ Y: -10 });
  }
  moveYPositive() {
    this.performManualMovement({ Y: 10 });
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
table .button {
  margin: 5px;
}
</style>
