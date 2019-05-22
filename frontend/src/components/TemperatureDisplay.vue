<template>
  <div>
    <h2 class="subtitle">Temperature graph</h2>
    <LoaderGuard>
      <canvas ref="chartCanvas" class="chart" height="400"></canvas>
    </LoaderGuard>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import Chart from "chart.js";
import LoadableMixin from "../LoadableMixin";
import LoaderGuard from "./LoaderGuard.vue";
import {
  TrackedValue,
  GetTrackedValueByNameWithMetaQuery,
  GetTrackedValueByNameWithMetaQueryVariables,
  SubscribeToTrackedValueUpdatedByNameSubscription,
  SubscribeToTrackedValueUpdatedByNameSubscriptionVariables
} from "../graphql-models-gen";
import getTrackedValueByNameWithMeta from "../../../queries/getTrackedValueByNameWithMeta.graphql";
import { QueryResult } from "vue-apollo/types/vue-apollo";
import gql from "graphql-tag";

@Component({
  components: { LoaderGuard }
})
export default class TemperatureDisplay extends mixins(LoadableMixin) {
  chart: Chart = null;
  readonly valuesToShow = [
    "hotendTemperature",
    "targetHotendTemperature",
    "hotbedTemperature",
    "targetHotbedTemperature"
  ];

  async created() {
    let tvMetas: { [k: string]: TrackedValue } = {};
    await this.withLoader(async () => {
      for (let valueToShow of this.valuesToShow) {
        let { data } = await this.$apollo.query<
          GetTrackedValueByNameWithMetaQuery
        >({
          query: getTrackedValueByNameWithMeta,
          variables: <GetTrackedValueByNameWithMetaQueryVariables>{
            name: valueToShow
          }
        });
        tvMetas[valueToShow] = data.trackedValue;
        let observable = this.$apollo.subscribe<
          QueryResult<SubscribeToTrackedValueUpdatedByNameSubscription>
        >({
          variables: <
            SubscribeToTrackedValueUpdatedByNameSubscriptionVariables
          >{
            name: valueToShow
          },

          query: gql`
            subscription subscribeToTrackedValueUpdatedByName($name: String!) {
              trackedValueUpdated(name: $name)
            }
          `
        });

        observable.subscribe(ev => {
          if (!this.chart) return;
          let value = ev.data.trackedValueUpdated;
          let dataset = this.chart.data.datasets.find(
            d => d.label === valueToShow
          );
          if (!dataset) {
            return; //wait for history
          }
          dataset.data.push(value);
          if (dataset.data.length > this.chart.data.labels.length) {
            dataset.data = dataset.data.slice(1);
          }
          this.chart.update();
        });
      }
    });
    this.$nextTick(() => {
      this.chart = new Chart(
        (this.$refs.chartCanvas as HTMLCanvasElement).getContext("2d"),
        {
          type: "line",
          data: {
            labels: Array(300)
              .fill(0)
              .map((_, i) => i),
            datasets: Object.keys(tvMetas).map(k => ({
              _ddd: k,
              borderColor: tvMetas[k].plotColor,
              label: tvMetas[k].name,
              data: tvMetas[k].history
            }))
          },
          options: {
            responsive: false
          }
        }
      );
    });
  }
}
</script>

<style scoped>
.chart {
  width: 100%;
}
</style>
