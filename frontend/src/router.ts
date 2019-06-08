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
  routes: <(RouteConfig & { menuName: string; menuIcon: string })[]>[
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
      menuIcon: "fa-print",
      component: BlankWrapper,
      children: [
        {
          path: "gcode-files",
          name: "gcode-files",
          menuName: "Gcode files",
          menuIcon: "fa-file",
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
      menuIcon: "fa-gamepad",
      component: BlankWrapper,
      children: [
        {
          path: "manual",
          name: "manual",
          menuName: "Manual",
          menuIcon: "fa-gamepad",
          component: () =>
            import(
              /* webpackChunkName: "manual" */ "./views/control/Manual.vue"
            )
        },
        {
          path: "serial-console",
          name: "serial-console",
          menuName: "Serial console",
          menuIcon: "fa-terminal",
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
      menuIcon: "fa-cogs",
      children: [
        {
          path: "settings",
          name: "settings",
          menuName: "Settings",
          menuIcon: "fa-wrench",
          component: () =>
            import(
              /* webpackChunkName: "settings" */ "./views/system/Settings.vue"
            ),
          children: [
            {
              path: "search",
              name: "settings-search",
              component: () =>
                import(
                  /* webpackChunkName: "settingssearch" */ "./views/system/Settings/Search.vue"
                )
            },
            {
              path: ":pageName",
              component: () =>
                import(
                  /* webpackChunkName: "settings" */ "./views/system/Settings/SettingsPage.vue"
                )
            }
          ]
        },
        {
          path: "system-info",
          name: "system-info",
          menuName: "System information",
          menuIcon: "fa-info",
          component: () =>
            import(
              /* webpackChunkName: "systeminfo" */ "./views/system/SystemInfo.vue"
            )
        },
        {
          path: "updates",
          name: "updates",
          menuName: "Updates",
          menuIcon: "fa-download",
          component: () =>
            import(
              /* webpackChunkName: "updates" */ "./views/system/Updates.vue"
            )
        }
      ]
    }
  ]
};

export default new Router(routerConfig);
