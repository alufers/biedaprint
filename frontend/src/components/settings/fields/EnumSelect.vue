<template>
  <div>
    <div class="field">
      <label class="label">
        <HighlightableText>{{fieldDescriptor.title}}</HighlightableText>
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
import HighlightableText from "../../HighlightableText";
import { JsonSchema } from "../../../util/settingsSchema";

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
  fieldDescriptor: JsonSchema;

  @Prop({
    type: String
  })
  value: string;

  @Emit("input")
  onFieldInput(e: Event) {
    return (<HTMLInputElement>e.target).value;
  }

  get options() {
    return this.fieldDescriptor.enum.map(e => ({ value: e, label: e }));
  }
}
</script>

<style>
</style>
