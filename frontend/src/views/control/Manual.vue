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
            >X+</button>
          </td>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-primary is-centered"
              @click="moveZPositive"
            >Z+</button>
          </td>
        </tr>
        <tr>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-success is-centered"
              @click="moveYNegative"
            >Y-</button>
          </td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-centered is-outlined"
              @click="centerXY"
            >
              <i class="fas fa-crosshairs"></i>
            </button>
          </td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-success is-centered"
              @click="moveYPositive"
            >Y+</button>
          </td>
          <td></td>
        </tr>
        <tr>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-danger is-centered"
              @click="moveXNegative"
            >X-</button>
          </td>
          <td></td>
          <td>
            <button
              :class="isLoadingClass"
              class="button is-primary is-centered"
              @click="moveZNegative"
            >Z-</button>
          </td>
        </tr>
        <tr>
          <td colspan="2">
            <button
              :class="isLoadingClass"
              class="button wide is-primary is-centered"
              @click="extrude"
            >Extrude</button>
          </td>
          <td colspan="2">
            <button
              :class="isLoadingClass"
              class="button wide is-primary is-centered"
              @click="retract"
            >Retract</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="tabs is-toggle">
      <ul>
        <li
          v-for="opt in allowedMovementAmounts"
          :key="opt"
          :class="opt === movementAmount && 'is-active'"
        >
          <a @click.prevent="movementAmount = opt">{{opt}}</a>
        </li>
      </ul>
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
  SendGcodeMutationVariables,
  PerformManualMovementMutation,
  PerformManualMovementMutationVariables,
  ManualMovementPositionVector
} from "../../graphql-models-gen";
import { performManualMovement } from "../../../../graphql/queries/performManualMovement.graphql";

@Component({})
export default class Manual extends mixins(LoadableMixin) {
  allowedMovementAmounts = [0.1, 1, 10, 100];
  movementAmount = 10;

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
    this.performManualMovement({ X: this.movementAmount });
  }
  moveXNegative() {
    this.performManualMovement({ X: -this.movementAmount });
  }
  moveYNegative() {
    this.performManualMovement({ Y: -this.movementAmount });
  }
  moveYPositive() {
    this.performManualMovement({ Y: this.movementAmount });
  }

  moveZPositive() {
    this.performManualMovement({ Z: this.movementAmount });
  }
  moveZNegative() {
    this.performManualMovement({ Z: -this.movementAmount });
  }
  extrude() {
    this.performManualMovement({ E: this.movementAmount });
  }
  retract() {
    this.performManualMovement({ E: -this.movementAmount });
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
button.wide {
  width: 110px;
}
table .button {
  margin: 5px;
}
</style>
