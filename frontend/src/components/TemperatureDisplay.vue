<!-- 
TemperatureDisplay shows a Chart.js chart of all the temperatures and targets.
Uses trackedValues to update it using subscriptions.
-->

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
import Chart, { ChartDataSets } from "chart.js";
import LoadableMixin from "../LoadableMixin";
import LoaderGuard from "./LoaderGuard.vue";
import {
  TrackedValue,
  GetTrackedValueByNameWithMetaQuery,
  GetTrackedValueByNameWithMetaQueryVariables,
  SubscribeToTrackedValueUpdatedByNameSubscription,
  SubscribeToTrackedValueUpdatedByNameSubscriptionVariables
} from "../graphql-models-gen";
import getTrackedValueByNameWithMeta from "../../../graphql/queries/getTrackedValueByNameWithMeta.graphql";
import { QueryResult } from "vue-apollo/types/vue-apollo";
import gql from "graphql-tag";

@Component({
  components: { LoaderGuard }
})
export default class TemperatureDisplay extends mixins(LoadableMixin) {
  chart: Chart | null = null;
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
          if (dataset.data.length > this!.chart.data.labels.length) {
            dataset.data = dataset.data.slice(1);
          }
          this.chart.update();
        });
      }
    });
    this.$nextTick(() => {
      if (!this.$refs.chartCanvas) {
        throw new Error("Chart canvas not ready!");
      }
      this.chart = new Chart(
        (this.$refs.chartCanvas as HTMLCanvasElement).getContext("2d")!,
        {
          type: "line",
          data: {
            labels: Array(300)
              .fill(0)
              .map((_, i) => i.toString()),
            datasets: Object.keys(tvMetas).map(
              k =>
                <ChartDataSets>{
                  _ddd: k,
                  borderColor: tvMetas[k].plotColor,
                  borderDash: tvMetas[k].plotDash,
                  label: tvMetas[k].name,
                  data: tvMetas[k].history,
                  pointRadius: 0
                }
            )
          },
          options: {
            responsive: false,
            tooltips: {
              enabled: false
            },
            scales: {
              yAxes: [
                {
                  ticks: {
                    callback(val: number) {
                      return (
                        val.toString() +
                        " " +
                        tvMetas[Object.keys(tvMetas)[0]].unit
                      );
                    },
                    beginAtZero: true,
                    suggestedMax: 300
                  }
                }
              ]
            }
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
