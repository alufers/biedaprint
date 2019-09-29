<template>
  <div class="field">
    <label>
      <HighlightableText>{{fieldDescriptor.title}}</HighlightableText>
    </label>
    <div class="field-body">
      <div class="field" :class="{'has-addons': !!fieldDescriptor.unit}">
        <div class="control is-expanded">
          <input :type="inputType" class="input" :value="value" @input="onFieldInput">
        </div>
        <p class="control" v-if="!!fieldDescriptor.unit">
          <a class="button is-static">{{fieldDescriptor.unit}}</a>
        </p>
      </div>
    </div>
    <p class="help">
      <HighlightableText>{{fieldDescriptor.description}}</HighlightableText>
    </p>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Emit } from "vue-property-decorator";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import HighlightableText from "../../HighlightableText";
import { JsonSchema } from "../../../util/settingsSchema";

@Component({
  components: {
    HighlightableText
  }
})
export default class TextField extends Vue {
  @Prop({
    required: true,
    type: Object
  })
  fieldDescriptor: JsonSchema;

  @Prop({})
  value: any;

  @Emit("input")
  onFieldInput(e: Event) {
    let rawValue = (<HTMLInputElement>e.target).value;
    if (this.fieldDescriptor.type === "integer") {
      return parseInt(rawValue);
    }
    if (this.fieldDescriptor.type === "float") {
      return parseFloat(rawValue);
    }
    return rawValue;
  }

  get inputType() {
    switch (this.fieldDescriptor.type) {
      case "float":
      case "integer":
        return "number";
      default:
        return "text";
    }
  }
}
</script>

<style>
</style>
