<template>
  <div>
    <div class="field">
      <label class="label">{{name}}</label>

      <div class="print-stat">
        <div class="controls tags has-addons are-medium">
          <span class="tag">Actual</span>
          <span class="tag">{{temperature.toFixed(2)}}</span>
        </div>

        <div class="controls tags has-addons are-medium">
          <span class="tag">Target</span>
          <span class="tag" :class="{'is-danger': target > 0,}">
            <input
              class="temperature-input"
              type="number"
              min="0"
              max="300"
              v-model.number="targetEdit"
              @keyup.enter="setTarget"
            >
          </span>
        </div>

        <div class="field has-addons">
          <p class="control">
            <a
              class="button"
              @click="setTarget"
              :class="{'is-primary': targetEdit !== target, ...isLoadingClass}"
            >SET</a>
          </p>
          <p class="control">
            <a class="button" @click="heaterOff" :class="isLoadingClass">OFF</a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import TrackedValueSubscription from "../TrackedValueSubscription";
import LoadableMixin from "../LoadableMixin";
import { sendGcode } from "../../../queries/sendGcode.graphql";
import {
  SendGcodeMutation,
  SendGcodeMutationVariables
} from "../graphql-models-gen";

@Component({})
export default class HeaterControl extends mixins(LoadableMixin) {
  @Prop({ type: String })
  name!: string;
  @Prop({ type: String })
  temperatureTrackedValueName!: string;
  @Prop({ type: String })
  targetTrackedValueName!: string;
  @Prop({ type: String })
  temperatureGcode!: string;

  @TrackedValueSubscription(function(this: HeaterControl) {
    return this.targetTrackedValueName;
  })
  target = 0;
  @TrackedValueSubscription(function(this: HeaterControl) {
    return this.temperatureTrackedValueName;
  })
  temperature = 0;

  targetEdit = 0;

  heaterOff() {
    //this.connection.sendMessage("sendGCODE", { data: `${this.temperatureGcode} S0` });
    this.withLoader(async () => {
      this.targetEdit = 0;
      await this.$apollo.mutate<SendGcodeMutation>({
        mutation: sendGcode,
        variables: <SendGcodeMutationVariables>{
          cmd: `${this.temperatureGcode} S0`
        }
      });
    });
  }
  setTarget() {
    this.withLoader(async () => {
      await this.$apollo.mutate<SendGcodeMutation>({
        mutation: sendGcode,
        variables: <SendGcodeMutationVariables>{
          cmd: `${this.temperatureGcode} S${this.targetEdit}`
        }
      });
    });
  }
  @Watch("target")
  targetWatch(newTarget: number, oldTarget: number) {
    if (this.targetEdit === oldTarget) {
      this.targetEdit = newTarget;
    }
  }
}
</script>


<style scoped>
.temperature-input {
  font-size: 1rem;
  color: #555;
  border: none transparent;
  background: transparent;
  outline: none;
  width: auto;
}

.is-primary .temperatur-input {
  color: white;
}

.print-stat {
  height: 35px;
  display: flex;
  align-items: flex-start;
  align-content: center;
  justify-content: flex-start;
}
.print-stat .value {
  font-weight: bold;
}

.controls {
  padding-top: 2px;
  margin-right: 8px;
  height: 30px;
}
</style>
