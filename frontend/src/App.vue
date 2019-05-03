<template>
  <div>
    <Navbar :socketStatus="socketStatus"/>
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

    this.connection.connect();
  }
};
</script>