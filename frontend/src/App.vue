<template>
  <div>
    <Navbar/>
    <section class="section">
      <div class="container">
        <!-- <keep-alive> -->
        <router-view></router-view>
        <!-- </keep-alive> -->
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Navbar from "./components/Navbar.vue";
import "bulma/css/bulma.css";
import "@fortawesome/fontawesome-free/css/all.css";

@Component({
  components: {
    Navbar
  }
})
export default class BlankWrapper extends Vue {}
</script>
