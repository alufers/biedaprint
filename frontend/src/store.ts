import Vue from "vue";
import Vuex from "vuex";
import AlertsModule from './modules/AlertsModule';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    AlertsModule
  }
});
