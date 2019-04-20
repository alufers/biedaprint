<template>
  <div>
    <h2 class="subtitle">Settings</h2>
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
      <button class="button is-primary" @click="save">Save</button>
    </div>
  </div>
</template>
<script>
export default {
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
      settings: null
    };
  },
  inject: ["connection"],
  methods: {
    save() {
      this.connection.sendMessage("saveSettings", this.settings);
    }
  },
  created() {
    this.connection.sendMessage("serialList");
    this.connection.sendMessage("getSettings");
    this.connection.on("message.serialList", ({ ports }) => {
      this.serialPorts = ports;
    });
    this.connection.on("message.getSettings", set => {
      this.settings = set;
    });
  }
};
</script>