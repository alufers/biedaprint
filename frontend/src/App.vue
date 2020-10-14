<template>
  <div>
    <Navbar/>
    <section class="section">
      <div class="container">
        <Alerts/>
        <router-view></router-view>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Navbar from "./components/Navbar.vue";
import Alerts from "./components/Alerts.vue";
import "bulma/css/bulma.css";
import "bulma-extensions/dist/css/bulma-extensions.min.css";
import "@fortawesome/fontawesome-free/css/all.css";

@Component({
  components: {
    Navbar,
    Alerts
  }
})
export default class App extends Vue {}
</script>
