<template>
  <div>
    <h2 class="title">System information</h2>
    <button class="button is-primary" @click="loadData">
      <i class="fas fa-sync"></i>
    </button>
    <div v-if="systemInfo">
      <table class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(value, key) in systemInfo" :key="key">
            <td>{{key}}</td>
            <td>{{value}}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      systemInfo: null
    };
  },
  created() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.connection.sendMessage("getSysteminfo");
    }
  },
  connectionSubscriptions: {
    "message.getSysteminfo"(info) {
      this.systemInfo = info;
    }
  }
};
</script>

<style>
</style>
