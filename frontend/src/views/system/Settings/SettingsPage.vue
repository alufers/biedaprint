<template>
  <div>
    <h3 class="subtitle">{{pageData.name}}</h3>
    <p class="has-text-grey-light">{{pageData.description}}</p>
    <br>
    <br>
    <component
      v-for="field of fields"
      :key="field.paramName"
      :is="fieldComponent(field)"
      :fieldDescriptor="field"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import settingsSchema from "../../../assets/settings-schema.json";
import SettingsFieldDescriptor from "../../../types/SettingsFieldDescriptor";
import SettingsPageDescriptor from "../../../types/SettingsPageDescriptor";
import TextField from "../../../components/settings/fields/TextField.vue";

@Component({})
export default class SettingsPage extends Vue {
  get fields(): SettingsFieldDescriptor[] {
    return settingsSchema.fields.filter(
      (f: SettingsFieldDescriptor) => f.pageEnumName === this.pageData.enumName
    );
  }
  get pageData(): SettingsPageDescriptor {
    return settingsSchema.pages.find(
      p => p.paramName === this.$route.params.pageName
    );
  }
  fieldComponent(field: SettingsFieldDescriptor) {
    return TextField;
  }
}
</script>

<style>
</style>
