<template>
  <div>
    <canvas ref="chartCanvas" width="400" height="400"></canvas>
  </div>
</template>

<script>
import Vue from "vue";
import Chart from "chart.js";

export default {
  data() {
    return {
      valuesToShow: ["hotendTemperature", "targetHotendTemperature"],
      values: {}
    };
  },
  inject: ["connection"],
  mounted() {
    // for (let v of this.valuesToShow) {
    //   this.connection.sendMessage("subscribeToTrackedValue", {
    //     name: v
    //   });
    // }
    // this.connection.on("message.trackedValueUpdated", ({ value, name }) => {
    //   Vue.set(this.values, name, value);
    // });
    let chart = new Chart(this.$refs.chartCanvas.getContext("2d"), {
      type: "line",
      data: {
        labels: Array(300)
          .fill(0)
          .map((_, i) => i),
        datasets: []
      },
      options: {
        responsive: false
      }
    });
    for (let v of this.valuesToShow) {
      this.connection.sendMessage("getTrackedValue", {
        name: v
      });
      this.connection.sendMessage("subscribeToTrackedValue", {
        name: v
      });
    }
    this.connection.on("message.getTrackedValue", ({ trackedValue }) => {
      chart.data.datasets.push({
        borderColor: trackedValue.plotColor,
        label: trackedValue.name,
        data: trackedValue.history
      });
      chart.update();
    });
    this.connection.on("message.trackedValueUpdated", ({ value, name }) => {
      let dataset = chart.data.datasets.find(d => d.label === name);
      dataset.data.push(value);
      if (dataset.data.length > chart.data.labels.length) {
        dataset.data = dataset.data.slice(1);
      }
      chart.update();
    });
  }
};
</script>

<style>
</style>
