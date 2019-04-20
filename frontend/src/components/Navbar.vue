<template>
  <nav class="navbar is-dark" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <img src="@/assets/logo.png" alt="Biedaprint logo" srcset>
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
        <router-link class="navbar-item" to="/" active-class="huba">Biedaprint</router-link>

        <div class="navbar-item has-dropdown is-hoverable">
          <router-link class="navbar-link" to="/control">Control</router-link>
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
          <span
            class="tag"
            :class="{'is-danger': socketStatus === 'closed'}"
          >{{socketStatus}}</span>
        </span>
        <span class="navbar-item">
          Serial status: &nbsp;
          <span
            class="tag"
            :class="{'is-danger': serialStatus === 'disconnected' || serialStatus === 'error', 'is-success': serialStatus === 'connected' }"
          >{{serialStatus}}</span>
        </span>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  props: ["socketStatus", "serialStatus"],
  data() {
    return {
      navbarActive: false
    };
  }
};
</script>

<style>
</style>
