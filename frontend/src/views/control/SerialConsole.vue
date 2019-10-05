<template>
  <div>
    <h2 class="title">Serial console</h2>
    <div class="columns">
      <div class="column is-9">
        <div class="console-wrapper">
          <button
            class="button is-primary is-rounded follow-log"
            :class="{'hidden': isScrolledToBottom}"
            @click="scrollToBottom"
          >
            <span class="icon is-small">
              <i class="fas fa-arrow-down"></i>
            </span>
            <span>Scroll to bottom</span>
          </button>
          <div class="box is-family-code console" ref="console" @scroll="updateScrolledToBottom">
            <div v-for="(l, i) in lines" :key="i">{{l}}</div>
          </div>
        </div>
        <div class="is-flex">
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
          />
          &nbsp;
          <button
            class="button is-primary"
            :class="isLoadingClass"
            @click="sendCommand"
          >Send</button>
        </div>
      </div>
      <div class="column is-4">
        <GcodeDocs @useGcode="useGcodeFromDocs" :currentCommand="currentCommand" />
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import GcodeDocs from "../../components/GcodeDocs.vue";
import sendConsoleCommand from "../../../../graphql/queries/sendConsoleCommand.graphql";
import getScrollbackBufferAndRecentCommands from "../../../../graphql/queries/getScrollbackBufferAndRecentCommands.graphql";
import serialConsoleDataSubscription from "../../../../graphql/queries/serialConsoleDataSubscription.graphql";
import LoadableMixin from "../../LoadableMixin";
import {
  SendConsoleCommandMutation,
  SendConsoleCommandMutationVariables,
  GetScrollbackBufferQuery,
  SerialConsoleDataSubscriptionSubscription,
  GetScrollbackBufferAndRecentCommandsQuery
} from "../../graphql-models-gen";
import gql from "graphql-tag";
// import { QueryResult } from "vue-apollo/types/vue-apollo";

@Component({
  components: {
    GcodeDocs
  }
})
export default class SerialConsole extends mixins(LoadableMixin) {
  scrollback = "...\n";
  currentCommand = "";
  recentCommands: string[] = [];
  currentRecentCommand = 0;
  isScrolledToBottom = true;
  created() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<
        GetScrollbackBufferAndRecentCommandsQuery
      >({
        query: getScrollbackBufferAndRecentCommands
      });
      this.scrollback = data.scrollbackBuffer;
      this.recentCommands = data.recentCommands;

      this.scrollToBottom();

      let obs = this.$apollo.subscribe<
        any // QueryResult<SerialConsoleDataSubscriptionSubscription>
      >({
        query: serialConsoleDataSubscription
      });
      obs.subscribe(val => {
        let shouldScroll = false;
        const consoleDiv = this.$refs.console as HTMLDivElement;
        if (
          consoleDiv.scrollTop ===
          consoleDiv.scrollHeight - consoleDiv.offsetHeight
        ) {
          shouldScroll = true;
        }
        this.scrollback += val.data.serialConsoleData;
        if (shouldScroll) {
          this.scrollToBottom();
        }
      });
    });
    this.updateScrolledToBottom();
  }
  async sendCommand() {
    if (this.loading) return;
    let commandToSend = this.currentCommand;
    this.recentCommands.push(this.currentCommand);
    this.currentRecentCommand = 0;
    this.currentCommand = "";
    await this.withLoader(async () => {
      await this.$apollo.mutate<SendConsoleCommandMutation>({
        mutation: sendConsoleCommand,
        variables: <SendConsoleCommandMutationVariables>{
          cmd: commandToSend
        }
      });
    });
  }
  useGcodeFromDocs(gcode: string) {
    this.currentCommand = gcode + " ";
    (this.$refs.commandInput as HTMLInputElement).focus();
  }

  previousRecentCommand() {
    if (this.recentCommands.length - this.currentRecentCommand > 0) {
      this.currentRecentCommand++;
      this.currentCommand = this.recentCommands[
        this.recentCommands.length - this.currentRecentCommand
      ];
    }
  }
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
  }
  resetCurrentRecentCommand(ev: KeyboardEvent) {
    if (ev.keyCode === 38 || ev.keyCode === 40) return;
    this.currentRecentCommand = 0;
  }
  get lines() {
    return this.scrollback.split("\n");
  }
  updateScrolledToBottom() {
    const consoleDiv = this.$refs.console as HTMLDivElement;
    if (!consoleDiv) return;
    this.isScrolledToBottom =
      consoleDiv.scrollTop ===
      consoleDiv.scrollHeight - consoleDiv.offsetHeight;
  }
  scrollToBottom() {
    this.$nextTick(() => {
      (this.$refs.console as HTMLDivElement).scrollTop = (this.$refs
        .console as HTMLDivElement).scrollHeight;
    });
    this.updateScrolledToBottom();
  }
}
</script>

<style scoped>
.console {
  min-height: 500px;
  max-height: 500px;
  overflow-y: scroll;
}
.console-wrapper {
  position: relative;
  margin-bottom: 10px;
}
.follow-log {
  position: absolute;
  top: 10px;
  right: 25px;
  opacity: 1;
  transform: scale(1);
  transition: 100ms ease-in-out;
}
.hidden {
  opacity: 0;
  transform: scale(0) translateY(-30px);
}
</style>
