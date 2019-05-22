<template>
  <div>
    <div class="columns">
      <div class="column is-one-third">
        <PrinterConnectionWidget class="box"/>
        <!-- <CurrentPrintWidget class="box"/> -->
      </div>
      <div class="column">
        <div class="box">
          <TemperatureDisplay/>
          <!--<HeaterControl
            name="Hotend"
            temperatureTrackedValueName="hotendTemperature"
            targetTrackedValueName="targetHotendTemperature"
            temperatureGcode="M104"
          />
          <HeaterControl
            name="Hotbed"
            temperatureTrackedValueName="hotbedTemperature"
            targetTrackedValueName="targetHotbedTemperature"
            temperatureGcode="M140"
          />-->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TemperatureDisplay from "@/components/TemperatureDisplay";
import HeaterControl from "@/components/HeaterControl";
import PrinterConnectionWidget from "@/components/PrinterConnectionWidget";

import connectionMixin from "@/connectionMixin";
// import TrackedValueTextDisplay from "@/components/TrackedValueTextDisplay";
import CurrentPrintWidget from "@/components/CurrentPrintWidget";

export default {
  name: "home",
  mixins: [connectionMixin],
  components: {
    TemperatureDisplay,
    HeaterControl,
    PrinterConnectionWidget,
    // TrackedValueTextDisplay,
    CurrentPrintWidget
  },
  data() {
    return {};
  }
};
</script>

<style scoped>
.box {
  padding: 16px;
}
</style>
