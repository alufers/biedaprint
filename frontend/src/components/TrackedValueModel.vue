<template>
  <div></div>
</template>

<script>
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  props: ["valueName"],
  data() {
    return {
      value: null
    };
  },
  mounted() {
    this.connection.sendMessage("getTrackedValue", {
      name: this.valueName
    });
    this.connection.sendMessage("subscribeToTrackedValue", {
      name: this.valueName
    });
  },
  connectionSubscriptions: {
    "message.trackedValueUpdated"({ value, name }) {
      if (name === this.valueName) {
        this.value = value;
        this.$emit("change", value);
      }
    },
    "message.getTrackedValue"({ trackedValue }) {
      if (trackedValue.name === this.valueName) {
        this.value = trackedValue.value;
        this.$emit("change", trackedValue.value);
      }
    }
  }
};
</script>

<style>
</style>
