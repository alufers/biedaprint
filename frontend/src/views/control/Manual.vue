<template>
  <div>
    <h3 class="title">Manual control</h3>
    <table class="table">
      <tbody>
        <tr>
          <td>
            <!-- home all axes -->
            <button
              class="button is-dark is-centered is-outlined"
              @click="sendGCODE('G28')"
              title="Home all axes"
            >
              <i class="fas fa-home"></i>
            </button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>X</td>
          <td></td>
          <td></td>
          <td></td>
          <td>Z</td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <!-- home X -->
            <button
              class="button is-danger is-centered is-outlined"
              @click="sendGCODE('G28 X')"
              title="Home X"
            >
              <i class="fas fa-home"></i>
            </button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <!-- home Z -->
            <button
              class="button is-primary is-centered is-outlined"
              @click="sendGCODE('G28 Z')"
              title="Home Z"
            >
              <i class="fas fa-home"></i>
            </button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">-100</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">-100</button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">-10</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">-10</button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">-1</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">-1</button>
          </td>
        </tr>
        <tr>
          <td>Y</td>
          <td>
            <button
              class="button is-success is-centered is-outlined"
              @click="sendGCODE('G28 Y')"
              title="Home Y"
            >
              <i class="fas fa-home"></i>
            </button>
          </td>
          <td>
            <button class="button is-success is-centered">-100</button>
          </td>
          <td>
            <button class="button is-success is-centered">-10</button>
          </td>
          <td>
            <button class="button is-success is-centered">-1</button>
          </td>
          <td>
            <button class="button is-centered is-info is-outlined">&bull;</button>
          </td>
          <td>
            <button class="button is-success is-centered">+1</button>
          </td>
          <td>
            <button class="button is-success is-centered">+10</button>
          </td>
          <td>
            <button class="button is-success is-centered">+100</button>
          </td>
          <td></td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">+1</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">+10</button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">+1</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">+10</button>
          </td>
        </tr>
        <tr>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-danger is-centered">+100</button>
          </td>
          <td></td>
          <td></td>
          <td></td>
          <td>
            <button class="button is-primary is-centered">+100</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  inject: ["connection"],
  methods: {
    sendGCODE(gcode) {
      this.connection.sendMessage("sendGCODE", gcode);
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
