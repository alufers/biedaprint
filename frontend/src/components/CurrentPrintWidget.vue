<template>
  <div>
    <TrackedValueModel @change="isPrinting = $event" valueName="isPrinting"/>
    <TrackedValueModel @change="printProgress = $event" valueName="printProgress"/>
    <TrackedValueModel @change="printOriginalName = $event" valueName="printOriginalName"/>
    <article class="message is-primary" v-if="isPrinting">
      <div class="message-header">
        <p>Print status</p>
      </div>
      <div class="message-body">
        <progress class="progress is-primary" :value="printProgress" max="100">{{printProgress}}</progress>
        <div class="print-stat">
          <div>Print name</div>
          <div class="value">{{printOriginalName}}</div>
        </div>
        <div class="print-stat">
          <div>Print progress (lines)</div>
          <div class="value">{{printProgress.toFixed(2)}}%</div>
        </div>
        <br>
        <br>
        <button class="button is-danger" @click="abortJob">
          <span class="icon">
            <i class="fas fa-stop"></i>
          </span>
          <span>Abort job</span>
        </button>
      </div>
    </article>
    <article class="message is-dark" v-else>
      <div class="message-body msg-noprint">
        <p>Nothing is being printed at the moment. You can select or uload a file to be printed.</p>
        <br>
        <br>
        <router-link to="/print/gcode-files" class="button is-primary">View gcode files</router-link>
      </div>
    </article>
  </div>
</template>

<script>
import TrackedValueModel from "@/components/TrackedValueModel";
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      isPrinting: false,
      printProgress: 0,
      printOriginalName: ""
    };
  },
  components: { TrackedValueModel },
  methods: {
    abortJob() {
      this.connection.sendMessage("abortPrintJob");
    }
  }
};
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
