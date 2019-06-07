<template>
  <LoaderGuard>
    <h2 class="subtitle">Settings</h2>
    <progress class="progress is-large is-primary" max="100" v-if="loading">15%</progress>
    <div class="columns">
      <div class="column is-one-fifth">
        <SettingsMenu/>
      </div>
      <div class="column box">
        <router-view/>
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