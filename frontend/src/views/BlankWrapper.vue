<template>
  <div>
    <router-view/>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

@Component({})
export default class BlankWrapper extends Vue {}
</script>

<style>
</style>
