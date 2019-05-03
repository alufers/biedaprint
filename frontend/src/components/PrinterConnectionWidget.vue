<template>
  <div>
    <h2 class="subtitle">Connect to printer</h2>
    <div v-if="settings && (serialStatus === 'disconnected' || serialStatus == 'error')">
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
      <button
        class="button is-success"
        @click="connectToSerial"
        v-if="serialStatus === 'disconnected' || serialStatus == 'error'"
      >Connect to printer</button>
    </div>
    <button
      class="button is-danger"
      @click="disconnectFromSerial"
      v-if="serialStatus === 'connected'"
    >Disconnect from printer</button>
  </div>
</template>
<script>
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      rates: [
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
      ],
      serialPorts: [],
      settings: null,
      serialStatus: null
    };
  },
  methods: {
    connectToSerial() {
      this.connection.sendMessage("connectToSerial");
    },
    disconnectFromSerial() {
      this.connection.sendMessage("disconnectFromSerial");
    },
    save() {
      this.connection.sendMessage("saveSettings", this.settings);
    }
  },
  created() {
    this.connection.sendMessage("serialList");
    this.connection.sendMessage("getSettings");
  },
  connectionSubscriptions: {
    "message.serialList"({ ports }) {
      this.serialPorts = ports;
    },
    "message.getSettings"(set) {
      this.settings = set;
    },
    "message.getSerialStatus"({ status }) {
      this.serialStatus = status;
    }
  }
};
</script>

<style scoped>
.button {
    width: 100%;
}
</style>