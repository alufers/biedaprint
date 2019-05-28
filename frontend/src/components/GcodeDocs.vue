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
      <HighlightableTextZone :highlights="searchHighlighs">
        <div class="panel-block" v-for="doc in filteredDocs" :key="doc.code">
          <div>
            <div>
              <code
                class="tag code-margin"
                :class="{'has-background-grey-lighter':code.startsWith('M')}"
                v-for="code in doc.codes"
                :key="code"
              >{{code}}</code>
              <HighlightableText>{{doc.title}}</HighlightableText>
            </div>
            <p class="has-text-grey">
              <HighlightableText>{{doc.brief}}</HighlightableText>
            </p>
            <div v-if="filteredDocs.length <= 2" class="params-list">
              <div class="menu">
                <p class="menu-label">Parameters</p>
                <ul class="menu-list">
                  <li v-for="(param, i) in doc.parameters" :key="i">
                    <span class="tag code-margin optional-mark" v-if="param.optional">Optional</span>
                    <div class="tags has-addons">
                      <code class="tag is-primary">{{param.tag === true ? 'Y' : param.tag}}</code>
                      <template v-if="param.values">
                        <code
                          class="tag"
                          v-for="(v, k) in param.values.filter(pv => pv.type !=='bool')"
                          :key="k"
                        >&lt;{{v.tag || v.type}}&gt;</code>
                      </template>
                    </div>
                    <div class="clearfix"></div>
                    <p class="has-text-grey is-size-7 param-desc">
                      <HighlightableText>{{param.description}}</HighlightableText>
                    </p>
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
      </HighlightableTextZone>
    </nav>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import _gcodeDocsData from "../assets/gcode-docs.json";
import HighlightableTextZone from "./HighlightableTextZone.vue";
import HighlightableText from "./HighlightableText";
import Fuse from "fuse.js";

const gcodeDocsData: any = _gcodeDocsData;

interface FuseData {
  key: string;
  code: string;
  title: string;
  brief: string;
}

@Component({
  components: {
    HighlightableTextZone,
    HighlightableText
  }
})
export default class GcodeDocs extends Vue {
  @Prop({ type: String })
  currentCommand!: string;
  searchQuery = "";
  dataKeys = Object.keys(gcodeDocsData);
  forceLocalSearch = false;
  fuse: Fuse<FuseData> = null;
  created() {
    this.fuse = new Fuse<FuseData, Fuse.FuseOptions<FuseData>>(
      this.dataForFuse,
      {
        keys: ["code", "title", "brief"],
        id: "key",
        caseSensitive: false,
        tokenize: true
      }
    );
  }

  get dataForFuse(): FuseData[] {
    return this.dataKeys.map(k => ({
      key: k,
      code: k,
      title: gcodeDocsData[k].title,
      brief: gcodeDocsData[k].brief
    }));
  }
  get filteredKeys() {
    let query = this.searchQuery.trim();

    if (this.currentCommand.trim() !== "" && !this.forceLocalSearch) {
      let commandQuery = this.currentCommand
        .trim()
        .toUpperCase()
        .split(" ")[0];
      return this.dataKeys.filter(k => k.indexOf(commandQuery) !== -1);
    }
    if (query === "") {
      return this.dataKeys.slice(0, 10);
    }
    return this.fuse.search(query, { limit: 10 }) as unknown as string[];
  }

  get filteredDocs() {
    return this.filteredKeys.map(k => gcodeDocsData[k]);
  }

  get searchHighlighs() {
    return this.searchQuery
      .split(" ")
      .map(q => q.trim().toLowerCase())
      .filter(q => q.length > 0);
  }

  @Watch("searchQuery")
  searchQueryWatch() {
    this.forceLocalSearch = true;
  }

  @Watch("currentCommand")
  currentCommandWatch() {
    this.forceLocalSearch = false;
  }
}
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
