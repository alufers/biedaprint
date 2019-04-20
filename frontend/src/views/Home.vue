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
  </div>
</template>

<script>
import TemperatureDisplay from "@/components/TemperatureDisplay";
export default {
  name: "home",
  inject: ["connection"],
  components: { TemperatureDisplay },
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
  created() {
    this.connection.on(
      "message.getSerialStatus",
      ({ status }) => (this.serialStatus = status)
    );
  }
};
</script>
