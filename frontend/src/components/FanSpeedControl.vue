<!--
HeaterControl is the UI element which shows the current temperature of a heater and allows it to pe controller by the user. 
It supports reading temperature presets from the settings as well as entering custom temps.
-->
<template>
  <LoaderGuard>
    <div class="field">
      <label class="label"> {{ name }}</label>

      <div class="print-stat">
        <div class="controls tags has-addons are-medium">
          <span class="tag tag-label">
            <span
              class="orb"
              :style="{
                background: trackedValueMeta && trackedValueMeta.plotColor,
              }"
            ></span>
            Actual</span
          >
          <span class="tag temp-value">{{ fanSpeed.toFixed(0) }}</span>
        </div>

        <div class="controls tags has-addons are-medium">
          <span class="tag tag-label">Target</span>
          <span class="tag" :class="{ 'is-warning': fanSpeed <= 10 }">
            <div class="dropdown" :class="{ 'is-active': showPresetsDropdown }">
              <div class="dropdown-trigger">
                <input
                  ref="temperatureInput"
                  class="slider"
                  type="range"
                  min="0"
                  max="255"
                  step="1"
                  v-model.number="targetEdit"
                  @keyup.enter="setTarget"
                  @focus="showPresetsDropdown = true"
                  @blur="hidePresetsDropdown"
                />
                <span
                  class="icon is-small"
                  @click="$refs.temperatureInput.focus()"
                >
                  <i class="fas fa-angle-down" aria-hidden="true"></i>
                </span>
              </div>
              <div class="dropdown-menu" id="dropdown-menu" role="menu">
                <div class="dropdown-content">
                  <a
                    href="#"
                    class="dropdown-item"
                    v-for="(tp, i) in temperaturePresets"
                    :key="i"
                    @click.prevent="selectPreset(i)"
                    >{{ tp.name }} ({{ tp.value }} °C)</a
                  >
                </div>
              </div>
            </div>
          </span>
        </div>
        <div class="field has-addons">
          <p class="control">
            <a
              class="button"
              @click="setTarget"
              :class="{
                'is-primary': targetEdit !== fanSpeed,
                ...isLoadingClass,
              }"
              >SET</a
            >
          </p>
          <p class="control">
            <a class="button" @click="heaterOff" :class="isLoadingClass">OFF</a>
          </p>
        </div>
      </div>
    </div>
  </LoaderGuard>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import TrackedValueSubscription from "../decorators/TrackedValueSubscription";
import LoadableMixin from "../LoadableMixin";
import { sendGcode } from "../../../graphql/queries/sendGcode.graphql";
import {
  SendGcodeMutation,
  SendGcodeMutationVariables,
  GetTemperaturePresetsQuery,
  TrackedValue,
} from "../graphql-models-gen";
import ApolloQuery from "../decorators/ApolloQuery";
import { getTemperaturePresets } from "../../../graphql/queries/getTemperaturePresets.graphql";
import { setTimeout } from "timers";
import { Presets } from "../types/settings";
import TrackedValueMeta from "../decorators/TrackedValueMeta";
import LoaderGuard from "./LoaderGuard.vue";

@Component({
  components: {
    LoaderGuard,
  },
})
export default class HeaterControl extends mixins(LoadableMixin) {
  @Prop({ type: String })
  name!: string;
  @Prop({ type: String })
  fanSpeedTrackedValue!: string;

  @Prop({ type: String })
  temperatureGcode!: string;
  @Prop({ type: String })
  temperaturePresetKey: keyof Presets;

  temperaturePresetsRaw: Presets[] = null;

  showPresetsDropdown = false;

  @TrackedValueSubscription(function(this: HeaterControl) {
    return this.fanSpeedTrackedValue;
  })
  fanSpeed = 0;

  @TrackedValueMeta(function(this: HeaterControl) {
    return this.fanSpeedTrackedValue;
  })
  trackedValueMeta: TrackedValue = null;

  targetEdit = 0;

  created() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<GetTemperaturePresetsQuery>({
        query: getTemperaturePresets,
      });
      this.temperaturePresetsRaw = data.settings;
    });
  }

  heaterOff() {
    //this.connection.sendMessage("sendGCODE", { data: `${this.temperatureGcode} S0` });
    this.withLoader(async () => {
      this.targetEdit = 0;
      await this.$apollo.mutate<SendGcodeMutation>({
        mutation: sendGcode,
        variables: <SendGcodeMutationVariables>{
          cmd: `${this.temperatureGcode} S0`,
        },
      });
    });
  }

  setTarget() {
    this.withLoader(async () => {
      await this.$apollo.mutate<SendGcodeMutation>({
        mutation: sendGcode,
        variables: <SendGcodeMutationVariables>{
          cmd: `${this.temperatureGcode} S${this.targetEdit}`,
        },
      });
    });
  }

  @Watch("target")
  targetWatch(newTarget: number, oldTarget: number) {
    if (this.targetEdit === oldTarget) {
      this.targetEdit = newTarget;
    }
  }

  get temperaturePresets() {
    if (!this.temperaturePresetsRaw) return [];
    return this.temperaturePresetsRaw.map((tp) => ({
      name: tp.name,
      value: <number>tp[this.temperaturePresetKey],
    }));
  }

  selectPreset(index: number) {
    this.targetEdit = this.temperaturePresets[index].value;
    this.showPresetsDropdown = false;
  }

  hidePresetsDropdown() {
    setTimeout(() => (this.showPresetsDropdown = false), 200); // delay before hiding the dropdown so that the  browser has time to register the click event
  }
}
</script>

<style scoped>
.slider {
  margin-top: 20px !important;
  width: 210px;
}
.is-danger .temperature-input {
  color: white;
}

.dropdown-trigger {
  display: flex;
  flex-direction: row;
 align-items: center;
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

.temp-value {
  width: 100px;
}

.tag-label {
  color: #838383;
}

.label {
  margin-top: 10px;
}

.preset-field {
  margin-right: 8px;
}

.orb {
  width: 10px;
  height: 10px;
  border-radius: 1000px;
  margin-right: 5px;
}
</style>
