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
    <TemperatureDisplay/>
    <TrackedValueTextDisplay valueName="printProgress"/>
  </div>
</template>

<script>
import TemperatureDisplay from "@/components/TemperatureDisplay";
import connectionMixin from "@/connectionMixin";
import TrackedValueTextDisplay from "@/components/TrackedValueTextDisplay";

export default {
  name: "home",
  mixins: [connectionMixin],
  components: { TemperatureDisplay, TrackedValueTextDisplay },
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
