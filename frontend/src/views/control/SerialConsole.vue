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
          @keyup="resetCurrentRecentCommand"
          @keyup.enter="sendCommand"
          @keyup.up="previousRecentCommand"
          @keyup.down="nextRecentCommand"
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
  name: "SerialConsole",
  mixins: [connectionMixin],
  components: {
    GcodeDocs
  },
  data() {
    return {
      scrollback: "...\n",
      currentCommand: "",
      recentCommands: [],
      currentRecentCommand: 0
    };
  },
  methods: {
    sendCommand() {
      this.connection.sendMessage("sendConsoleCommand", {
        data: this.currentCommand
      });
      this.recentCommands.push(this.currentCommand);
      this.currentRecentCommand = 0;
      this.currentCommand = "";
    },
    useGcodeFromDocs(gcode) {
      this.currentCommand = gcode + " ";
      this.$refs.commandInput.focus();
    },
    previousRecentCommand() {
      if (this.recentCommands.length - this.currentRecentCommand > 0) {
        this.currentRecentCommand++;
        this.currentCommand = this.recentCommands[
          this.recentCommands.length - this.currentRecentCommand
        ];
      }
    },
    nextRecentCommand() {
      if (this.currentRecentCommand > 0) {
        this.currentRecentCommand--;
        if (this.currentRecentCommand === 0) {
          this.currentCommand = "";
          return;
        }
        this.currentCommand = this.recentCommands[
          this.recentCommands.length - this.currentRecentCommand
        ];
      }
    },
    resetCurrentRecentCommand(ev) {
      if (ev.keyCode === 38 || ev.keyCode === 40) return;
      this.currentRecentCommand = 0;
    }
  },
  computed: {
    lines() {
      return this.scrollback.split("\n");
    }
  },
  created() {
    this.connection.sendMessage("getScrollbackBuffer");
    this.connection.sendMessage("getRecentCommands");
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
    },
    "message.getRecentCommands"(recentCommands) {
      this.recentCommands = recentCommands;
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
