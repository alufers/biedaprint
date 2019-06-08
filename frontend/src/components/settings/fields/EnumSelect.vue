<template>
  <div>
    <div class="field">
      <label class="label">
        <HighlightableText>{{fieldDescriptor.label}}</HighlightableText>
      </label>
      <div class="control">
        <div class="select">
          <select :value="value" @input="onFieldInput">
            <option v-for="opt in options" :key="opt.value" :value="opt.value">{{opt.label}}</option>
          </select>
        </div>
      </div>
      <p class="help">
        <HighlightableText>{{fieldDescriptor.description}}</HighlightableText>
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Emit } from "vue-property-decorator";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import settingsSchema from "../../../assets/settings-schema.json";
import HighlightableText from "../../HighlightableText";

@Component({
  components: {
    HighlightableText
  }
})
export default class EnumSelect extends Vue {
  @Prop({
    required: true,
    type: Object
  })
  fieldDescriptor: SettingsFieldDescriptor;

  @Prop({
    type: String
  })
  value: string;

  @Emit("input")
  onFieldInput(e: Event) {
    return (<HTMLInputElement>e.target).value;
  }

  get options() {
    return settingsSchema.enums.find(
      (e: any) => e.name === this.fieldDescriptor.enumTypeName
    ).values as any;
  }
}
</script>

<style>
</style>
