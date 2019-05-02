<template>
  <div>
    <div>Welcome to biedaprint!</div>
    <button
      class="button is-success"
      @click="connectToSerial"
      v-if="serialStatus === 'disconnected' || serialStatus == 'error'"
    >Connect to printer</button>
    <button
      class="button is-danger"
      @click="disconnectFromSerial"
      v-if="serialStatus === 'connected'"
    >Disconnect from printer</button>
    <div class="columns">
      <div class="column">
        <TemperatureDisplay/>
      </div>
      <div class="column">
        <CurrentPrintWidget/>
      </div>
    </div>
  </div>
</template>

<script>
import TemperatureDisplay from "@/components/TemperatureDisplay";
import connectionMixin from "@/connectionMixin";
// import TrackedValueTextDisplay from "@/components/TrackedValueTextDisplay";
import CurrentPrintWidget from "@/components/CurrentPrintWidget";

export default {
  name: "home",
  mixins: [connectionMixin],
  components: {
    TemperatureDisplay,
    // TrackedValueTextDisplay,
    CurrentPrintWidget
  },
  data() {
    return {
      serialStatus: null
    };
  },
  methods: {
    connectToSerial() {
      this.connection.sendMessage("connectToSerial");
    },
    disconnectFromSerial() {
      this.connection.sendMessage("disconnectFromSerial");
    }
  },
  connectionSubscriptions: {
    "message.getSerialStatus"({ status }) {
      this.serialStatus = status;
    }
  }
};
</script>
