<template>
  <div>
    <div class="field">
      <label class="label">{{fieldDescriptor.label}}</label>
      <input :type="inputType" class="input" :value="value" @input="onFieldInput">
      <p class="help">{{fieldDescriptor.description}}</p>
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

  @Prop({})
  value: any;

  @Emit("input")
  onFieldInput(e: Event) {
    let rawValue = (<HTMLInputElement>e.target).value;
    if (this.fieldDescriptor.editComponent === "IntField") {
      return parseInt(rawValue);
    }
    if (this.fieldDescriptor.editComponent === "FloatField") {
      return parseFloat(rawValue);
    }
    return rawValue;
  }

  get inputType() {
    switch (this.fieldDescriptor.editComponent) {
      case "IntField":
      case "FloatField":
        return "number";
      default:
        return "text";
    }
  }
}
</script>

<style>
</style>
