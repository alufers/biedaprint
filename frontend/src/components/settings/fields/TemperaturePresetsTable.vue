<template>
  <div>
    <div class="field">
      <label class="label">
        <HighlightableText>{{ fieldDescriptor.title }}</HighlightableText>
      </label>

      <table class="table">
        <thead>
          <th>Name</th>
          <th>Hotend temperature (°C)</th>
          <th>Hotbed temperature (°C)</th>
          <th>Fan speed</th>
          <th></th>
        </thead>
        <tbody>
          <tr v-for="(tp, i) in value" :key="i">
            <td>
              <input class="input" type="text" v-model="tp.name" />
            </td>
            <td>
              <input
                class="input"
                type="number"
                v-model.number="tp.hotendTemperature"
              />
            </td>
            <td>
              <input
                class="input"
                type="number"
                v-model.number="tp.hotbedTemperature"
              />
            </td>
            <td>
              <input class="input" type="number" v-model.number="tp.fanSpeed" />
            </td>
            <td>
              <button
                class="button is-danger"
                @click="deleteTemperaturePreset(i)"
              >
                <span class="icon is-small">
                  <i class="fas fa-trash"></i>
                </span>
              </button>
            </td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="4">
              <button class="button is-primary" @click="addTemperaturePreset()">
                <span class="icon is-small">
                  <i class="fas fa-plus"></i>
                </span>
                <span>Add temperature preset</span>
              </button>
            </td>
          </tr>
        </tfoot>
      </table>
      <p class="help">
        <HighlightableText>{{ fieldDescriptor.description }}</HighlightableText>
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Emit, Watch } from "vue-property-decorator";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import { Presets as TemperaturePreset } from "../../../types/settings";
import HighlightableText from "../../HighlightableText";
import { JsonSchema } from "../../../util/settingsSchema";

@Component({
  components: {
    HighlightableText,
  },
})
export default class TemperaturePresetsTable extends Vue {
  @Prop({
    required: true,
    type: Object,
  })
  fieldDescriptor: JsonSchema;

  @Prop({
    type: Array,
  })
  value: TemperaturePreset[];

  addTemperaturePreset() {
    this.$emit("input", [
      ...this.value,
      <TemperaturePreset>{
        __typename: "TemperaturePreset",
        name: "New",
        hotendTemperature: 0,
        hotbedTemperature: 0,
        fanSpeed: 0,
      },
    ]);
  }

  deleteTemperaturePreset(i: number) {
    this.$emit(
      "input",
      this.value.filter((_, ix) => ix !== i)
    );
  }

  @Watch("value", { deep: true })
  watchValue() {
    this.$emit("input", this.value);
  }
}
</script>

<style></style>
