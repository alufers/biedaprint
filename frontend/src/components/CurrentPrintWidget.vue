<template>
  <div>
    <div v-if="isPrinting">
      <div class="print-stat">
        <h2 class="subtitle">Print status</h2>
        <button class="button is-danger" @click="abortJob" :class="isLoadingClass">
          <span class="icon">
            <i class="fas fa-stop"></i>
          </span>
          <span>Abort job</span>
        </button>
      </div>
      <div>
        <progress class="progress is-primary" :value="printProgress" max="100">{{printProgress}}</progress>
        <div class="print-stat">
          <div>Print name</div>
          <div class="value">{{printOriginalName}}</div>
        </div>
        <div class="print-stat">
          <div>Print progress (lines)</div>
          <div class="value">{{printProgress.toFixed(2)}}%</div>
        </div>
        <div class="print-stat">
          <div>Print start time</div>
          <div class="value">{{printStartTime | formatDate}}</div>
        </div>
        <div class="print-stat">
          <div>Layer</div>
          <div class="value">{{printCurrentLayer}}/{{printTotalLayers}}</div>
        </div>
      </div>
    </div>
    <div class v-else>
      <div class="msg-noprint">
        <p>Nothing is being printed at the moment. You can select or upload a file to be printed.</p>
        <br>
        <br>
        <router-link to="/print/gcode-files" class="button is-primary">View gcode files</router-link>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import { DateTime } from "luxon";
import TrackedValueSubscription from "../TrackedValueSubscription";
import LoadableMixin from "../LoadableMixin";
import { abortPrintJob } from "../../../queries/abortPrintJob.graphql";

@Component({
  filters: {
    formatDate(value) {
      if (!value) return "";
      let dt = DateTime.fromISO(value);
      return dt.toISODate() + " " + dt.toLocaleString(DateTime.TIME_24_SIMPLE);
    }
  }
})
export default class CurrentPrintWidget extends mixins(LoadableMixin) {
  @TrackedValueSubscription("isPrinting")
  isPrinting = false;
  @TrackedValueSubscription("printProgress")
  printProgress = 0;
  @TrackedValueSubscription("printOriginalName")
  printOriginalName = "";
  @TrackedValueSubscription("printStartTime")
  printStartTime: string = null;
  @TrackedValueSubscription("printCurrentLayer")
  printCurrentLayer = 0;
  @TrackedValueSubscription("printTotalLayers")
  printTotalLayers = 0;

  abortJob() {
    this.withLoader(async () => {
      await this.$apollo.mutate({
        mutation: abortPrintJob
      });
    });
  }
}
</script>


<style scoped>
.msg-noprint {
  text-align: center;
  padding-top: 60px;
  padding-bottom: 60px;
}
.print-stat {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
.print-stat .value {
  font-weight: bold;
}
</style>
