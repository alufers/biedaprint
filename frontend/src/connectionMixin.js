import Vue from "vue";

export default {
  inject: ["connection"],
  data() {
    return {
      connSubscribers: []
    };
  },
  mounted() {
    let subs = this.$options.connectionSubscriptions;
    if (!subs) {
      return;
    }

    for (let k of Object.keys(subs)) {
      let boundFunc = subs[k].bind(this);
      this.connection.on(k, boundFunc);

      this.connSubscribers.push({
        unbind: () => {
          // lets use closures to hide the blackmagick fuckery from vue
          this.connection.removeListener(k, boundFunc);
        }
      });
    }
  },
  destroyed() {
    for (let l of this.connSubscribers) {
      l.unbind();
    }
  }
};
