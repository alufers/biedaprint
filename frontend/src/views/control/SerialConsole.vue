<template>
  <div>
    <h2 class="title">Serial console</h2>
    <div class="columns">
      <div class="column is-9">
        <div class="box is-family-code console" ref="console">
          <div v-for="(l, i) in lines" :key="i">{{l}}</div>
        </div>
        <input
          class="input"
          type="text"
          placeholder="Send commands"
          v-model="currentCommand"
          @keyup.enter="sendCommand"
          ref="commandInput"
        >
      </div>
      <div class="column is-4">
        <GcodeDocs @useGcode="useGcodeFromDocs" :currentCommand="currentCommand"/>
      </div>
    </div>
  </div>
</template>

<script>
import GcodeDocs from "@/components/GcodeDocs.vue";
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  components: {
    GcodeDocs
  },
  data() {
    return {
      scrollback: "...\n",
      currentCommand: ""
    };
  },
  methods: {
    sendCommand() {
      this.connection.sendMessage("serialWrite", {
        data: this.currentCommand + "\r\n"
      });
      this.currentCommand = "";
    },
    useGcodeFromDocs(gcode) {
      this.currentCommand = gcode + " ";
      this.$refs.commandInput.focus();
    }
  },
  computed: {
    lines() {
      return this.scrollback.split("\n");
    }
  },
  created() {
    this.connection.sendMessage("getScrollbackBuffer");
    this.connection.once("message.getScrollbackBuffer", ({ data }) => {
      this.scrollback += data;
      this.$nextTick(() => {
        this.$refs.console.scrollTop = this.$refs.console.scrollHeight;
      });
    });
  },
  connectionSubscriptions: {
    "message.serialConsole"({ data }) {
      this.scrollback += data;
      this.$nextTick(() => {
        this.$refs.console.scrollTop = this.$refs.console.scrollHeight;
      });
    }
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
