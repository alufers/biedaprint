<template>
  <div>
    <div class="columns">
      <div class="column is-one-third">
        <PrinterConnectionWidget class="box"/>
        <CurrentPrintWidget class="box"/>
      </div>
      <div class="column">
        <div class="box">
        <TemperatureDisplay class=""/>
        <TemperatureControls class="" :heaters="heaters"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TemperatureDisplay from "@/components/TemperatureDisplay";
import TemperatureControls from "@/components/TemperatureControls";
import PrinterConnectionWidget from "@/components/PrinterConnectionWidget";

import connectionMixin from "@/connectionMixin";
// import TrackedValueTextDisplay from "@/components/TrackedValueTextDisplay";
import CurrentPrintWidget from "@/components/CurrentPrintWidget";

export default {
  name: "home",
  mixins: [connectionMixin],
  components: {
    TemperatureDisplay,
    TemperatureControls,
    PrinterConnectionWidget,
    // TrackedValueTextDisplay,
    CurrentPrintWidget
  },
  data() {
    return {
      heaters: ["Hotend", "Heatbed"],
    };
  },
};
</script>

<style scoped>
.box {
  padding: 16px;
}

</style>
