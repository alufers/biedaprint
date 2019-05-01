<template>
  <div>
    <Navbar :socketStatus="socketStatus" :serialStatus="serialStatus"/>
    <section class="section">
      <div class="container">
        <Alerts/>
        <!-- <keep-alive> -->
        <router-view></router-view>
        <!-- </keep-alive> -->
      </div>
    </section>
  </div>
</template>

<script>
import Connection from "./Connection";
import Alerts from "./components/Alerts";
import Navbar from "./components/Navbar";
import "bulma/css/bulma.css";
import "@fortawesome/fontawesome-free/css/all.css";

export default {
  components: {
    Alerts,
    Navbar
  },
  data() {
    return {
      connection: new Connection(),
      socketStatus: "disconnected",
      serialStatus: "?",
      navbarActive: false
    };
  },
  provide() {
    return {
      connection: this.connection
    };
  },
  methods: {
    sendJSON(type, data) {
      this.socket.send(
        JSON.stringify({
          type,
          data
        })
      );
    }
  },
  created() {
    this.connection.on("statusChanged", sta => (this.socketStatus = sta));
    this.connection.on("message.getSerialStatus", data => {
      this.serialStatus = data.status;
      setTimeout(() => this.connection.sendMessage("getSerialStatus"), 1000);
    });

    this.connection.connect();
    this.connection.sendMessage("getSerialStatus");
  }
};
</script>