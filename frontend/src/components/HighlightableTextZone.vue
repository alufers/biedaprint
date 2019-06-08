<!-- 
HighlightableTextZone provides informatio nabout the fragments of text to highlight in HighlightableText via the provide/inject api.
It is used in places where we want to visually emphasize search queries in the searched content.
-->
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
      $highlightsComputed: () => this.highlightsComputed // weird hack to inject reactive values
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
