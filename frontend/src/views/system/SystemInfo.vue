<template>
  <div>
    <h2 class="title">System information</h2>
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
export default {
  inject: ["connection"],
  data() {
    return {
      systemInfo: null
    };
  },
  created() {
    this.connection.sendMessage("getSysteminfo");
    this.connection.on("message.getSysteminfo", info => {
      this.systemInfo = info;
    });
  }
};
</script>

<style>
</style>
