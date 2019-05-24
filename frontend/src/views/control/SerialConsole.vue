<template>
  <div>
    <h2 class="title">Serial console</h2>
    <div class="columns">
      <div class="column is-9">
        <div class="box is-family-code console" ref="console">
          <div v-for="(l, i) in lines" :key="i">{{l}}</div>
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
          >
          &nbsp;
          <button
            class="button is-primary"
            :class="isLoadingClass"
            @click="sendCommand"
          >Send</button>
        </div>
      </div>
      <div class="column is-4">
        <GcodeDocs @useGcode="useGcodeFromDocs" :currentCommand="currentCommand"/>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import GcodeDocs from "../../components/GcodeDocs.vue";
import sendConsoleCommand from "../../../../queries/sendConsoleCommand.graphql";
import getScrollbackBufferAndRecentCommands from "../../../../queries/getScrollbackBufferAndRecentCommands.graphql";
import serialConsoleDataSubscription from "../../../../queries/serialConsoleDataSubscription.graphql";
import LoadableMixin from "../../LoadableMixin";
import {
  SendConsoleCommandMutation,
  SendConsoleCommandMutationVariables,
  GetScrollbackBufferQuery,
  SerialConsoleDataSubscriptionSubscription,
  GetScrollbackBufferAndRecentCommandsQuery
} from "../../graphql-models-gen";
import gql from "graphql-tag";
import { QueryResult } from "vue-apollo/types/vue-apollo";

@Component({
  components: {
    GcodeDocs
  }
})
export default class SerialConsole extends mixins(LoadableMixin) {
  scrollback = "...\n";
  currentCommand = "";
  recentCommands = [];
  currentRecentCommand = 0;
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
        QueryResult<SerialConsoleDataSubscriptionSubscription>
      >({
        query: serialConsoleDataSubscription
      });
      obs.subscribe(val => {
        this.scrollback += val.data.serialConsoleData;
        this.scrollToBottom();
      });
    });
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
  useGcodeFromDocs(gcode) {
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
  resetCurrentRecentCommand(ev) {
    if (ev.keyCode === 38 || ev.keyCode === 40) return;
    this.currentRecentCommand = 0;
  }
  get lines() {
    return this.scrollback.split("\n");
  }
  scrollToBottom() {
    this.$nextTick(() => {
      (this.$refs.console as HTMLDivElement).scrollTop = (this.$refs
        .console as HTMLDivElement).scrollHeight;
    });
  }
}
</script>

<style scoped>
.console {
  min-height: 500px;
  max-height: 500px;
  overflow-y: scroll;
}
</style>
