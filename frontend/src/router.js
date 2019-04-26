import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  linkActiveClass: "is-active",
  routes: [
    {
      path: "/",
      name: "home",
      component: Home
    },
    {
      path: "/control/manual",
      name: "manual",
      component: () =>
        import(/* webpackChunkName: "manual" */ "./views/control/Manual.vue")
    },
    {
      path: "/control/serial-console",
      name: "serial-console",
      component: () =>
        import(/* webpackChunkName: "serialconsole" */ "./views/control/SerialConsole.vue")
    },
    {
      path: "/system/settings",
      name: "settings",
      component: () =>
        import(/* webpackChunkName: "settings" */ "./views/system/Settings.vue")
    },
    {
      path: "/system/system-info",
      name: "system-info",
      component: () =>
        import(/* webpackChunkName: "systeminfo" */ "./views/system/SystemInfo.vue")
    },
    {
      path: "/print/gcode-files",
      name: "gcode-files",
      component: () =>
        import(/* webpackChunkName: "systeminfo" */ "./views/print/GcodeFiles.vue")
    }
  ]
});
