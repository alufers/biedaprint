<template>
  <section class="section">
    <div class="container">
      <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
          <a class="navbar-item" href="/">
            <img src="./logo.png" alt="Biedaprint logo" srcset>
          </a>
          <a
            role="button"
            class="navbar-burger burger"
            aria-label="menu"
            aria-expanded="false"
            data-target="navbarBasicExample"
            @click="navbarActive = !navbarActive"
            :class="{'is-active': navbarActive}"
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>
        <div id="navbarBasicExample" class="navbar-menu" :class="{'is-active': navbarActive}">
          <div class="navbar-start">
            <router-link class="navbar-item" to="/">Biedaprint</router-link>

            <div class="navbar-item has-dropdown is-hoverable">
              <router-link class="navbar-link" to="/system">Control</router-link>
              <div class="navbar-dropdown">
                <router-link class="navbar-item" to="/control/manual">Manual</router-link>
                <router-link class="navbar-item" to="/control/serial-console">Serial console</router-link>
              </div>
            </div>
            <div class="navbar-item has-dropdown is-hoverable">
              <router-link class="navbar-link" to="/system">System</router-link>
              <div class="navbar-dropdown">
                <router-link class="navbar-item" to="/system/settings">Settings</router-link>
                <router-link class="navbar-item" to="/system/system-info">System information</router-link>
              </div>
            </div>
          </div>
          <div class="navbar-end">
            <span class="navbar-item">
              Socket status: &nbsp;
              <span class="tag">{{socketStatus}}</span>
            </span>
            <span class="navbar-item">
              Serial status: &nbsp;
              <span class="tag">{{serialStatus}}</span>
            </span>
          </div>
        </div>
      </nav>
      <Alerts/>
      <keep-alive>
        <router-view></router-view>
      </keep-alive>
    </div>
  </section>
</template>

<script>
import Connection from "./Connection";
import Alerts from "./components/Alerts";
import "bulma/css/bulma.css";
import "@fortawesome/fontawesome-free/css/all.css";

export default {
  components: {
    Alerts
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
    });
    this.connection.on("open", () => {
      setInterval(() => this.connection.sendMessage("getSerialStatus"), 1000);
    });
    this.connection.connect();
  }
};
</script>