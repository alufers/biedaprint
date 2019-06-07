<template>
  <div>
    <div class="field">
      <label class="label">{{fieldDescriptor.label}}</label>
      <input type="number" class="input" :value="value" @input="onFieldInput">
      <p class="help" v-if="fieldDescriptor.description">{{fieldDescriptor.description}}</p>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Emit } from "vue-property-decorator";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";

@Component({})
export default class TextField extends Vue {
  @Prop({
    required: true,
    type: Object
  })
  fieldDescriptor: SettingsFieldDescriptor;

  @Prop({
    type: Number
  })
  value: number;

  @Emit("input")
  onFieldInput(e: Event) {
    return parseInt((<HTMLInputElement>e.target).value);
  }
}
</script>

<style>
</style>
