<template>
  <div>
    <div class="columns">
      <div class="column is-one-third">
        <PrinterConnectionWidget class="box" />
        <CurrentPrintWidget class="box" />
      </div>
      <div class="column">
        <div class="box">
          <TemperatureDisplay />
          <HeaterControl
            name="Hotend"
            temperatureTrackedValueName="hotendTemperature"
            temperaturePresetKey="hotendTemperature"
            targetTrackedValueName="targetHotendTemperature"
            temperatureGcode="M104"
          />
          <HeaterControl
            name="Hotbed"
            temperatureTrackedValueName="hotbedTemperature"
            temperaturePresetKey="hotbedTemperature"
            targetTrackedValueName="targetHotbedTemperature"
            temperatureGcode="M140"
          />
          <FanSpeedControl
            name="Fan speed"
            fanSpeedTrackedValue="fanSpeed"
            temperaturePresetKey="fanSpeed"
            temperatureGcode="M106"
          />
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import TemperatureDisplay from "@/components/TemperatureDisplay.vue";
import HeaterControl from "@/components/HeaterControl.vue";
import FanSpeedControl from "@/components/FanSpeedControl.vue";
import PrinterConnectionWidget from "@/components/PrinterConnectionWidget.vue";
import CurrentPrintWidget from "@/components/CurrentPrintWidget.vue";

@Component({
  components: {
    TemperatureDisplay,
    HeaterControl,
    PrinterConnectionWidget,
    CurrentPrintWidget,
    FanSpeedControl,
  },
})
export default class Home extends Vue {}
</script>

<style scoped>
.box {
  padding: 16px;
}
</style>
