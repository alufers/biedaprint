<template>
  <div>
    <article class="message" :class="'is-'+alert.type" v-for="alert in alerts" :key="alert.id">
      <div class="message-header">
        <p>Alert</p>
        <button class="delete" aria-label="delete" @click="deleteAlert(alert.id)"></button>
      </div>
      <div class="message-body">{{alert.content}}</div>
    </article>
  </div>
</template>

<script>
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      alerts: []
    };
  },
  methods: {
    deleteAlert(id) {
      this.alerts = this.alerts.filter(a => a.id !== id);
    }
  },
  created() {
  },
  connectionSubscriptions: {
    "message.alert"(a) {
      this.alerts.push({
        ...a,
        id: Math.random()
      });
    }
  }
};
</script>

<style>
</style>
