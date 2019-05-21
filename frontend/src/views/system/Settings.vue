<template>
  <div>
    <h2 class="subtitle">Settings</h2>
    <progress class="progress is-large is-primary" max="100" v-if="loading">15%</progress>
    <div v-if="settings">
      <div class="field">
        <label class="label">Serial</label>
        <div class="select">
          <select v-model="settings.serialPort">
            <option v-for="serial in serialPorts" :key="serial">{{serial}}</option>
          </select>
        </div>
      </div>
      <div class="field">
        <label class="label">Baud rate</label>
        <div class="select">
          <select v-model.number="settings.baudRate">
            <option v-for="rate in rates" :key="rate">{{rate}}</option>
          </select>
        </div>
      </div>
      <div class="field">
        <label class="label">Scrollback buffer size</label>
        <input class="input" type="number" v-model.number="settings.scrollbackBufferSize">
      </div>
      <div class="field">
        <label class="label">Data path</label>
        <input class="input" type="text" v-model="settings.dataPath">
      </div>
      <div class="field">
        <label class="label">Startup command</label>
        <input class="input" type="text" v-model="settings.startupCommand">
      </div>
      <button class="button is-primary" @click="save">Save</button>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import gql from "graphql-tag";
import getSettingsAndSerialPorts from "../../../../queries/getSettingsAndSerialPorts.graphql";
import updateSettings from "../../../../queries/updateSettings.graphql";
import {
  GetSettingsAndSerialPortsQuery,
  UpdateSettingsMutation,
  UpdateSettingsMutationVariables,
  Settings as SettingsModel
} from "../../graphql-models-gen";

@Component({})
export default class SettingsPage extends mixins(LoadableMixin) {
  readonly rates = [
    300,
    600,
    1200,
    2400,
    4800,
    9600,
    14400,
    19200,
    28800,
    38400,
    57600,
    115200,
    2500000
  ];
  serialPorts = [];
  settings: SettingsModel = null;
  created() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<GetSettingsAndSerialPortsQuery>({
        query: getSettingsAndSerialPorts
      });
      delete data.settings.__typename;
      this.settings = data.settings;
      this.serialPorts = data.serialPorts;
    });
  }
  save() {
    this.withLoader(async () => {
      await this.$apollo.mutate<UpdateSettingsMutation>({
        mutation: updateSettings,
        variables: <UpdateSettingsMutationVariables>{
          newSettings: this.settings
        }
      });
    });
  }
}
</script>