<template>
  <HighlightableTextCSSWrapper>
    <slot></slot>
  </HighlightableTextCSSWrapper>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Provide } from "vue-property-decorator";
import HighlightableTextCSSWrapper from "./HighlightableTextCSSWrapper.vue";

@Component({
  components: {
    HighlightableTextCSSWrapper
  },
  provide() {
    return {
      $highlightsComputed: () => this.highlightsComputed
    };
  }
})
export default class HighlightableTextZone extends Vue {
  @Prop()
  highlights: string[];

  get highlightsComputed() {
    return this.highlights;
  }
}
</script>

<style>
</style>
