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
        <template v-for="route in topLevelRoutes">
          <div
            v-if="route.children && route.children.length > 0"
            class="navbar-item has-dropdown is-hoverable"
            exact
            :key="route.path"
          >
            <router-link class="navbar-link" :to="route.path">{{route.menuName || route.name}}</router-link>
            <div class="navbar-dropdown" v-if="route.children && route.children.length > 0">
              <router-link
                class="navbar-item"
                :to="urlJoin(route.path, child.path)"
                v-for="child in route.children"
                :key="child.path"
              >{{child.menuName || child.name}}</router-link>
            </div>
          </div>
          <router-link
            :key="route.path"
            class="navbar-item"
            :to="route.path"
            v-else
          >{{route.menuName || route.name}}</router-link>
        </template>
      </div>
      <div class="navbar-end">
        <span class="navbar-item">
          <div class="tags has-addons">
            <span class="tag">Serial status</span>
            <span
              class="tag"
              :class="{'is-danger': serialStatus === 'disconnected' || serialStatus === 'error', 'is-success': serialStatus === 'connected' }"
            >{{serialStatus}}</span>
          </div>
        </span>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import TrackedValueSubscription from "../TrackedValueSubscription";
import { routerConfig } from "../router";
import urlJoin from "url-join";

@Component({})
export default class Navbar extends Vue {
  @TrackedValueSubscription("serialStatus")
  serialStatus = "?";
  navbarActive = false;
  get topLevelRoutes() {
    return routerConfig.routes;
  }
  urlJoin(...u: string[]) {
    return urlJoin(u);
  }
}
</script>

<style>
</style>
