<template>
  <div>
    <nav class="panel">
      <p class="panel-heading">Gcode docs</p>
      <div class="panel-block">
        <p class="control has-icons-left">
          <input class="input is-small" type="text" placeholder="Search" v-model="searchQuery">
          <span class="icon is-small is-left">
            <i class="fas fa-search" aria-hidden="true"></i>
          </span>
        </p>
      </div>
      <div class="panel-block" v-for="doc in filteredDocs" :key="doc.codes.toString()">
        <div>
          <div>
            <code
              class="tag code-margin"
              :class="{'has-background-grey-lighter':code.startsWith('M')}"
              v-for="code in doc.codes"
              :key="code"
            >{{code}}</code>
            <span>{{doc.title}}</span>
          </div>
          <p class="has-text-grey">{{doc.brief}}</p>
          <div v-if="filteredDocs.length <= 2" class="params-list">
            <div class="menu">
              <p class="menu-label">Parameters</p>
              <ul class="menu-list">
                <li v-for="(param, i) in doc.parameters" :key="i">
                  <span class="tag code-margin optional-mark" v-if="param.optional">Optional</span>
                  <div class="tags has-addons">
                    <code class="tag is-primary">{{param.tag === true ? 'Y' : param.tag}}</code>
                    <code
                      class="tag"
                      v-for="(v, k) in param.values.filter(pv => pv.type !=='bool')"
                      :key="k"
                    >&lt;{{v.tag || v.type}}&gt;</code>
                  </div>
                  <div class="clearfix"></div>
                  <p class="has-text-grey is-size-7 param-desc">{{param.description}}</p>
                </li>
              </ul>
            </div>
          </div>

          <div class="buttons buttons-marg">
            <button
              class="button is-primary is-outlined is-small"
              @click="$emit('useGcode', doc.codes[0])"
            >Use</button>
            <a
              class="button is-text is-small"
              :href="'http://marlinfw.org/docs/gcode/' + doc.base + '.html'"
              target="_blank"
            >More...</a>
          </div>
        </div>
      </div>
    </nav>
  </div>
</template>

<script>
import gcodeDocsData from "@/assets/gcode-docs.json";
export default {
  props: ["currentCommand"],
  data() {
    return {
      searchQuery: "",
      dataKeys: Object.keys(gcodeDocsData),
      forceLocalSearch: false
    };
  },
  computed: {
    filteredKeys() {
      let query = this.searchQuery.trim().toUpperCase();
      let searchBriefs = true;
      if (this.currentCommand.trim() !== "" && !this.forceLocalSearch) {
        query = this.currentCommand
          .trim()
          .toUpperCase()
          .split(" ")[0];
        searchBriefs = false;
      }
      let keys = this.dataKeys.filter(k => k.indexOf(query) !== -1);
      if (keys.length < 10 && searchBriefs) {
        this.dataKeys
          .filter(
            k => gcodeDocsData[k].brief.toUpperCase().indexOf(query) !== -1
          )
          .forEach(k => keys.push(k));
      }
      return keys.slice(0, 10);
    },
    filteredDocs() {
      return this.filteredKeys.map(k => gcodeDocsData[k]);
    }
  },
  watch: {
    searchQuery() {
      this.forceLocalSearch = true;
    },
    currentCommand() {
      this.forceLocalSearch = false;
    }
  }
};
</script>

<style scoped>
.code-margin {
  margin-right: 4px;
}
.buttons-marg {
  margin-top: 4px;
}
.params-list {
  margin-top: 10px;
}
.params-list .menu-label {
  margin-bottom: 0;
}
.params-list li {
  margin-top: 10px;
}
.param-desc {
  margin-left: 10px;
}
.tags {
  margin-bottom: 0;
}
.optional-mark {
  float: right;
}
.clearfix {
  clear: both;
}
</style>
