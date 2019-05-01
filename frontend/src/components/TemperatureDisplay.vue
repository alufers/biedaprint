<template>
  <div>
    <canvas ref="chartCanvas" width="400" height="400"></canvas>
  </div>
</template>

<script>
import Chart from "chart.js";
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      chart: null,
      valuesToShow: ["hotendTemperature", "targetHotendTemperature"],
      values: {}
    };
  },
  mounted() {
    this.chart = new Chart(this.$refs.chartCanvas.getContext("2d"), {
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
  },
  connectionSubscriptions: {
    "message.trackedValueUpdated"({ value, name }) {
      let dataset = this.chart.data.datasets.find(d => d.label === name);
      if (!dataset) {
        return; //wait for history
      }
      dataset.data.push(value);
      if (dataset.data.length > this.chart.data.labels.length) {
        dataset.data = dataset.data.slice(1);
      }
      this.chart.update();
    },
    "message.getTrackedValue"({ trackedValue }) {
      if (!this.valuesToShow.includes(trackedValue.name)) return;
      this.chart.data.datasets.push({
        borderColor: trackedValue.plotColor,
        label: trackedValue.name,
        data: trackedValue.history
      });
      this.chart.update();
    }
  }
};
</script>

<style>
</style>
