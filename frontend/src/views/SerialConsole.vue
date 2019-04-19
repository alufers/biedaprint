<template>
  <div>
    <h2 class="title">Serial console</h2>
    <div class="box is-family-code console" ref="console">
      <div v-for="(l, i) in lines" :key="i">{{l}}</div>
    </div>
    <input
      class="input"
      type="text"
      placeholder="Send commands"
      v-model="currentCommand"
      @keyup.enter="sendCommand"
    >
  </div>
</template>

<script>
export default {
  inject: ["connection"],
  data() {
    return {
      backbuffer: "...\n",
      currentCommand: ""
    };
  },
  methods: {
    sendCommand() {
      this.connection.sendMessage("serialWrite", {
        data: this.currentCommand + "\r\n"
      });
      this.currentCommand = "";
    }
  },
  computed: {
    lines() {
      return this.backbuffer.split("\n");
    }
  },
  created() {
    this.connection.on("message.serialConsole", ({ data }) => {
      this.backbuffer += data;
      this.$refs.console.scrollTop = this.$refs.console.scrollHeight;
    });
  }
};
</script>

<style scoped>
.console {
  min-height: 500px;
  max-height: 500px;
  overflow-y: scroll;
}
</style>
