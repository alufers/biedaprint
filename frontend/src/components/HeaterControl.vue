<template>
  <div>
    <TrackedValueModel @change="temperature = $event" :valueName="temperatureTrackedValueName"/>
    <TrackedValueModel @change="target = $event" :valueName="targetTrackedValueName"/>
    <div class="field">
      <label class="label">{{name}}</label>

      <div class="print-stat">
        <div class="controls tags has-addons are-medium">
          <span class="tag">Actual</span>
          <span class="tag">{{temperature}}</span>
        </div>

        <div class="controls tags has-addons are-medium">
          <span class="tag">Target</span>
          <span class="tag is-primary">
            <input
              class="temperature-input"
              type="number"
              min="0"
              max="300"
              v-model.number="targetEdit"
              @keyup.enter="setTarget"
            >
          </span>
        </div>

        <div class="field has-addons">
          <p class="control">
            <a class="button is-primary" @click="setTarget">SET</a>
          </p>
          <p class="control">
            <a class="button" @click="heaterOff">OFF</a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TrackedValueModel from "@/components/TrackedValueModel";
import connectionMixin from "@/connectionMixin";

export default {
  mixins: [connectionMixin],
  props: {
    name: String,
    temperatureTrackedValueName: String,
    targetTrackedValueName: String,
    temperatureGcode: String
  },
  methods: {
    heaterOff() {
      this.connection.sendMessage("sendGCODE", { data: `${this.temperatureGcode} S0` });
    },
    setTarget() {
      this.connection.sendMessage("sendGCODE", {
        data: `${this.temperatureGcode} S${this.targetEdit}`
      });
    }
  },
  watch: {
    target(newTarget, oldTarget) {
      if (this.targetEdit === oldTarget) {
        this.targetEdit = newTarget;
      }
    }
  },
  data() {
    return {
      target: 0,
      targetEdit: 0,
      temperature: 0
    };
  },
  components: { TrackedValueModel }
};
</script>

<style scoped>
.temperature-input {
  font-size: 1rem;
  color: white;
  border: none transparent;
  background: transparent;
  outline: none;
  width: auto;
}

.print-stat {
  height: 35px;
  display: flex;
  align-items: flex-start;
  align-content: center;
  justify-content: flex-start;
}
.print-stat .value {
  font-weight: bold;
}

.controls {
  padding-top: 2px;
  margin-right: 8px;
  height: 30px;
}
</style>
