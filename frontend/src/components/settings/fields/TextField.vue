<template>
  <div>
    <div class="field">
      <label for class="label">{{fieldDescriptor.label}}</label>
      <input type="text" class="input">
      <p class="help" v-if="fieldDescriptor.description">{{fieldDescriptor.description}}</p>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";

@Component({})
export default class TextField extends Vue {
  @Prop({
    required: true,
    type: Object
  })
  fieldDescriptor: SettingsFieldDescriptor;
}
</script>

<style>
</style>
