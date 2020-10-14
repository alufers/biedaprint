<!-- 
TemperatureDisplay shows an Apex chart of all the temperatures and targets.
Uses trackedValues to update it using subscriptions.
-->

<template>
  <div>
    <h2 class="subtitle">Temperature graph</h2>
    <LoaderGuard>
      <ApexChart
        type="line"
        height="400"
        ref="chart"
        :options="chartOptions"
        :series="chartData"
      ></ApexChart>
    </LoaderGuard>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import ApexChartClass from "apexcharts";
import { ApexOptions } from "apexcharts";
import ApexChart from "vue-apexcharts";
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
// import { QueryResult } from "vue-apollo/types/vue-apollo";
import gql from "graphql-tag";
import { Watch } from "vue-property-decorator";

@Component({
  components: { LoaderGuard, ApexChart }
})
export default class TemperatureDisplay extends mixins(LoadableMixin) {
  readonly valuesToShow = [
    "hotendTemperature",
    "targetHotendTemperature",
    "hotbedTemperature",
    "targetHotbedTemperature",
    "fanSpeed"
  ];
  chartOptions: ApexOptions = {
    chart: {
      id: "realtime",
      height: 400,
      type: "line",
      animations: {
        enabled: true,
        easing: "linear",
        dynamicAnimation: {
          speed: 1000
        }
      },
      toolbar: {
        show: false
      },
      zoom: {
        enabled: false
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      curve: "smooth",
      colors: [],
      dashArray: []
    },
    markers: {
      size: 0
    },
    xaxis: {
      type: "numeric",
      min: 0,
      range: 300
    },
    yaxis: {
      min: 0,
      max: 300
    },
    legend: {
      show: false
    }
  };
  chartData: ApexAxisChartSeries = [];
  async created() {
    let tvMetas: { [k: string]: TrackedValue } = {};
    await this.withLoader(async () => {
      for (let valueToShow of this.valuesToShow) {
        // first we grab information about the tracked value
        let { data } = await this.$apollo.query<
          GetTrackedValueByNameWithMetaQuery
        >({
          query: getTrackedValueByNameWithMeta,
          variables: <GetTrackedValueByNameWithMetaQueryVariables>{
            name: valueToShow
          },
          fetchPolicy: "network-only"
        });
        tvMetas[valueToShow] = data.trackedValue;
        this.chartOptions.stroke.colors.push(data.trackedValue.plotColor);
        (<number[]>this.chartOptions.stroke.dashArray).push(
          data.trackedValue.plotDash.length > 0 ? 5 : 0
        );
        this.chartData.push({
          name: data.trackedValue.name,
          data: data.trackedValue.history
        });
        // after creating the series we subscribe for updates

        let observable = this.$apollo.subscribe<
          // QueryResult<SubscribeToTrackedValueUpdatedByNameSubscription>
          any
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
          let value = ev.data.trackedValueUpdated;
          const series = this.chartData.find(t => t.name === valueToShow);

          if (!series) return; // bad update
          series.data.push(value);
          if (series.data.length > tvMetas[series.name].maxHistoryLength) {
            series.data = series.data.slice(1);
          }
        });
      }
    });
  }
  @Watch("chartData", { deep: true })
  onChartDataChanged() {
    if (!this.$refs.chart) return;
    ((this.$refs.chart as any) as ApexChartClass).updateSeries([
      ...this.chartData
    ]);
  }
}
</script>

<style scoped>
.chart {
  width: 100%;
}
</style>
