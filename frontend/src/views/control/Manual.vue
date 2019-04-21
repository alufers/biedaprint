<template>
  <div>
    <h3 class="title">Manual control</h3>
    <label class="label">Homing</label>
    <div class="field has-addons">
      <p class="control">
        <!-- home all axes -->
        <button
          class="button is-dark is-centered is-outlined"
          @click="sendGCODE('G28')"
          title="Home all axes"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <!-- home X -->
        <button
          class="button is-danger is-centered is-outlined"
          @click="sendGCODE('G28 X')"
          title="Home X"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <button
          class="button is-success is-centered is-outlined"
          @click="sendGCODE('G28 Y')"
          title="Home Y"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
      <p class="control">
        <!-- home Z -->
        <button
          class="button is-primary is-centered is-outlined"
          @click="sendGCODE('G28 Z')"
          title="Home Z"
        >
          <i class="fas fa-home"></i>
        </button>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  inject: ["connection"],
  methods: {
    sendGCODE(gcode) {
      this.connection.sendMessage("sendGCODE", { data: gcode });
    }
  }
};
</script>

<style scoped>
td {
  text-align: center;
}
button {
  width: 50px;
  height: 50px;
}
</style>
