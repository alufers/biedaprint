<template>
  <LoaderGuard>
    <h2 class="subtitle">Settings</h2>
    <progress class="progress is-large is-primary" max="100" v-if="loading">15%</progress>
    <div class="columns">
      <div class="column is-one-fifth">
        <SettingsMenu/>
      </div>
      <div class="column box">
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
            <label class="label">Serial parity</label>
            <div class="select">
              <select v-model.number="settings.parity">
                <option v-for="parity in parities" :key="parity">{{parity}}</option>
              </select>
            </div>
          </div>
          <div class="field">
            <label class="label">Serial data bits</label>
            <div class="select">
              <select v-model.number="settings.dataBits">
                <option v-for="dataBit in dataBits" :key="dataBit">{{dataBit}}</option>
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
          <div class="field">
            <label class="label">Temperature presets</label>
            <table class="table">
              <thead>
                <th>Name</th>
                <th>Hotend temperature (°C)</th>
                <th>Hotbed temperature (°C)</th>
                <th></th>
              </thead>
              <tbody>
                <tr v-for="(tp, i) in settings.temperaturePresets" :key="i">
                  <td>
                    <input class="input" type="text" v-model="tp.name">
                  </td>
                  <td>
                    <input class="input" type="number" v-model.number="tp.hotendTemperature">
                  </td>
                  <td>
                    <input class="input" type="number" v-model.number="tp.hotbedTemperature">
                  </td>
                  <td>
                    <button class="button is-danger" @click="deleteTemperaturePreset(i)">
                      <span class="icon is-small">
                        <i class="fas fa-trash"></i>
                      </span>
                    </button>
                  </td>
                </tr>
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="4">
                    <button class="button is-primary" @click="addTemperaturePreset()">
                      <span class="icon is-small">
                        <i class="fas fa-plus"></i>
                      </span>
                      <span>Add temperature preset</span>
                    </button>
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
          <button class="button is-primary" @click="save">Save</button>
        </div>
      </div>
    </div>
  </LoaderGuard>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../../LoadableMixin";
import gql from "graphql-tag";
import getSettingsAndSerialPorts from "../../../../graphql/queries/getSettingsAndSerialPorts.graphql";
import updateSettings from "../../../../graphql/queries/updateSettings.graphql";
import {
  GetSettingsAndSerialPortsQuery,
  UpdateSettingsMutation,
  UpdateSettingsMutationVariables,
  Settings as SettingsModel
} from "../../graphql-models-gen";
import LoaderGuard from "../../components/LoaderGuard.vue";
import SettingsMenu from "../../components/settings/SettingsMenu.vue";

@Component({
  components: {
    LoaderGuard,
    SettingsMenu
  }
})
export default class SettingsPage extends mixins(LoadableMixin) {
  readonly parities = ["NONE", "EVEN"];
  readonly dataBits = [5, 7, 8];
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
  serialPorts: string[] = [];
  settings: SettingsModel | null = null;
  created() {
    this.withLoader(async () => {
      let { data } = await this.$apollo.query<GetSettingsAndSerialPortsQuery>({
        query: getSettingsAndSerialPorts
      });
      delete data.settings.__typename;
      data.settings.temperaturePresets.forEach(tp => delete tp.__typename);
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
  deleteTemperaturePreset(i: number) {
    this.settings.temperaturePresets = this.settings.temperaturePresets.filter(
      (_, ix) => ix !== i
    );
  }
  addTemperaturePreset() {
    this.settings.temperaturePresets.push({
      name: "New",
      hotendTemperature: 0,
      hotbedTemperature: 0
    });
  }
}
</script>