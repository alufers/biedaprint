<!-- 
The Alerts component is used to show global information about errors on top of the page.
It is only used once in the main App component, because it pulls global data about alerts from the vuex store.
Various functions should use the store's actions to show their alerts in a consistent way.
-->
<template>
  <div class="alerts-container">
    <div class="message" v-for="alert in alerts" :key="alert.id" :class="alertClass(alert)">
      <div class="message-header">
        <p>{{alert.title || "Alert"}}</p>
        <button class="delete" aria-label="delete" @click="removeAlertById(alert.id)"></button>
      </div>
      <div class="message-body">
        <p>{{alert.content}}</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Alert, AlertType } from "../modules/AlertsModule";

@Component({})
export default class Alerts extends Vue {
  get alerts(): Alert[] {
    return this.$store.state.AlertsModule.alerts;
  }

  alertClass(alert: Alert) {
    return {
      "is-danger": alert.type === AlertType.error,
      "is-info": alert.type === AlertType.info,
      "is-success": alert.type === AlertType.success
    };
  }

  removeAlertById(id: number) {
    this.$store.dispatch("AlertsModule/removeAlertById", id);
  }
}
</script>

<style scoped>
.alerts-container {
  margin-bottom: 30px;
}
</style>
