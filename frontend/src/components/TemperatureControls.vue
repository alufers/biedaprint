<template>
  <div>
    <h2 class="subtitle">Temperature control</h2>
    <div v-for="(heater, index) in heaters" :key="`heater-${index}`" class="field" >
      <label class="label">{{heater}} temperature</label>

      <div class="print-stat">
        <div class="controls tags has-addons are-medium">
          <span class="tag">Actual</span>
          <span class="tag is-primary">135</span>
        </div>

        <div class="controls tags has-addons are-medium">
          <span class="tag">Target</span>
          <span class="tag is-info">
            <input
              class="temperature-input"
              placeholder="220"
              type="number"
              min="0"
              max="300"
              v-model="hotendTemperature"
            >
          </span>
        </div>

        <div class="field has-addons">
          <p class="control">
            <a class="button">SET</a>
          </p>
          <p class="control">
            <a class="button is-danger">OFF</a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TrackedValueModel from "@/components/TrackedValueModel";
import connectionMixin from "@/connectionMixin";
import { DateTime } from "luxon";

export default {
  props: {
    heaters: Array
  },
  mixins: [connectionMixin],
  data() {
    return {
      hotendTemperature: 0
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
