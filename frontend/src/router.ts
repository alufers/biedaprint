import Vue from "vue";
import Router, { RouteConfig, RouterOptions } from "vue-router";
import Home from "./views/Home.vue";

Vue.use(Router);

const BlankWrapper = () =>
  import(/* webpackChunkName: "manual" */ "./views/BlankWrapper.vue");

export const routerConfig: RouterOptions = {
  mode: "history",
  base: process.env.BASE_URL,
  linkActiveClass: "is-active",
  routes: <(RouteConfig & { menuName: string })[]>[
    {
      path: "/",
      name: "home",
      component: Home,
      menuName: "Biedaprint"
    },
    {
      path: "/print",
      name: "print",
      menuName: "Print",
      component: BlankWrapper,
      children: [
        {
          path: "gcode-files",
          name: "gcode-files",
          menuName: "Gcode files",
          component: () =>
            import(
              /* webpackChunkName: "gcodeFiles" */ "./views/print/GcodeFiles.vue"
            )
        }
      ]
    },
    {
      path: "/control",
      name: "control",
      menuName: "Control",
      component: BlankWrapper,
      children: [
        {
          path: "manual",
          name: "manual",
          menuName: "Manual",
          component: () =>
            import(
              /* webpackChunkName: "manual" */ "./views/control/Manual.vue"
            )
        },
        {
          path: "serial-console",
          name: "serial-console",
          menuName: "Serial console",
          component: () =>
            import(
              /* webpackChunkName: "serialconsole" */ "./views/control/SerialConsole.vue"
            )
        }
      ]
    },
    {
      path: "/system",
      name: "system",
      menuName: "System",
      component: BlankWrapper,
      children: [
        {
          path: "settings",
          name: "settings",
          menuName: "Settings",
          component: () =>
            import(
              /* webpackChunkName: "settings" */ "./views/system/Settings.vue"
            )
        },
        {
          path: "system-info",
          name: "system-info",
          menuName: "System information",
          component: () =>
            import(
              /* webpackChunkName: "systeminfo" */ "./views/system/SystemInfo.vue"
            )
        }
      ]
    }
  ]
};

export default new Router(routerConfig);
