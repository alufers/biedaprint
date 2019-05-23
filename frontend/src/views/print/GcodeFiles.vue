<template>
  <div>
    <h2 class="title">Gcode files</h2>
    <GcodeUploadZone/>
    <table class="table is-fullwidth is-hoverable" v-if="!!gcodeFileMetas">
      <thead>
        <tr>
          <th>Name</th>
          <th>Lines</th>
          <th>Time</th>
          <th>Filament used</th>
          <th>Layers</th>
          <th>Upload date</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="f in gcodeFileMetas" :key="f.gcodeFileName">
          <td>{{f.originalName}}</td>
          <td>{{f.totalLines}}</td>
          <td>{{f.printTime | formatDuration}}</td>
          <td>{{f.filamentUsedMm}} mm</td>
          <td>{{f.layerIndexes && f.layerIndexes.length}}</td>
          <td>{{f.uploadDate | formatDate}}</td>
          <td>
            <div class="field has-addons">
              <p class="control">
                <button class="button is-primary" @click="startPrintJob(f.gcodeFileName)">
                  <span class="icon is-small">
                    <i class="fas fa-print"></i>
                  </span>
                </button>
              </p>

              <p class="control">
                <button class="button is-danger" @click="gcodeFileToDelete = f">
                  <span class="icon is-small">
                    <i class="fas fa-trash"></i>
                  </span>
                </button>
              </p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="modal is-active" v-if="gcodeFileToDelete">
      <div class="modal-background" @click="gcodeFileToDelete = null"></div>
      <div class="modal-content"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">Are you sure?</p>
          <button class="delete" aria-label="close" @click="gcodeFileToDelete = null"></button>
        </header>
        <section class="modal-card-body">
          <p>
            Do you really want to delete
            <code>{{gcodeFileToDelete.originalName}}</code>?
          </p>
        </section>
        <footer class="modal-card-foot">
          <button class="button" @click="gcodeFileToDelete = null">Cancel</button>
          <button
            class="button is-danger"
            @click="deleteGcodeFile(gcodeFileToDelete.gcodeFileName)"
          >
            <span class="icon">
              <i class="fas fa-trash"></i>
            </span>
            <span>Delete</span>
          </button>
        </footer>
      </div>
      <button class="modal-close is-large" aria-label="close" @click="gcodeFileToDelete = null"></button>
    </div>
  </div>
</template>


<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import { DateTime, Duration } from "luxon";
import GcodeUploadZone from "../../components/GcodeUploadZone.vue";
import { startPrintJob } from "../../../../queries/startPrintJob.graphql";
import { deleteGcodeFile } from "../../../../queries/deleteGcodeFile.graphql";
import LoadableMixin from "../../LoadableMixin";
import {
  StartPrintJobMutation,
  StartPrintJobMutationVariables,
  DeleteGcodeFileMutation,
  DeleteGcodeFileMutationVariables,
  GetGcodeFileMetasQuery,
  GcodeFileMeta
} from "../../graphql-models-gen";
import { getGcodeFileMetas } from "../../../../queries/getGcodeFileMetas.graphql";
import ApolloQuery from "../../ApolloQuery";
import { Watch } from "vue-property-decorator";

@Component({
  components: {
    GcodeUploadZone
  },
  filters: {
    formatDate(value: string) {
      let dt = DateTime.fromISO(value);
      return dt.toISODate() + " " + dt.toLocaleString(DateTime.TIME_24_SIMPLE);
    },
    formatDuration(value: number) {
      let dur = Duration.fromObject({
        days: 0,
        hours: 0,
        minutes: 0,
        seconds: value
      });
      let durObj = dur.normalize().toObject();

      return Object.keys(durObj)
        .filter(k => durObj[k] !== 0 && k !== "seconds")
        .map(k => durObj[k].toFixed(0) + " " + k)
        .join(", ");
    }
  }
})
export default class GcodeFiles extends mixins(LoadableMixin) {
  @ApolloQuery({
    query: getGcodeFileMetas
  })
  gcodeFileMetas: GcodeFileMeta[] = null;
  gcodeFileToDelete: GcodeFileMeta = null; // used to show the confirm modal


  deleteGcodeFile(gcodeFilename: string) {
    this.gcodeFileToDelete = null;
    this.withLoader(async () => {
      await this.$apollo.mutate<DeleteGcodeFileMutation>({
        mutation: deleteGcodeFile,
        variables: <DeleteGcodeFileMutationVariables>{
          gcodeFilename
        }
      });
    });
  }
  startPrintJob(gcodeFilename: string) {
    this.withLoader(async () => {
      await this.$apollo.mutate<StartPrintJobMutation>({
        mutation: startPrintJob,
        variables: <StartPrintJobMutationVariables>{
          gcodeFilename
        }
      });
      this.$router.push("/"); // redirect to main page
    });
  }
}
</script>

<!--<script>
import GcodeUploadZone from "@/components/GcodeUploadZone";
import connectionMixin from "@/connectionMixin";
import { DateTime, Duration } from "luxon";

export default {
  mixins: [connectionMixin],
  data() {
    return {
      gcodeFiles: null,
      gcodeFileToDelete: null
    };
  },
  components: {
    GcodeUploadZone
  },
  connectionSubscriptions: {
    "message.getGcodeFileMetas"(metas) {
      this.gcodeFiles = metas;
    }
  },
  methods: {
    deleteGcodeFile(gcodeFileName) {
      this.gcodeFileToDelete = null;
      this.connection.sendMessage("deleteGcodeFile", {
        gcodeFileName
      });
      this.connection.sendMessage("getGcodeFileMetas");
    },
    startPrintJob(gcodeFileName) {
      this.connection.sendMessage("startPrintJob", {
        gcodeFileName
      });
      this.$router.push("/"); // redirect to main page
    }
  },
 ,
  created() {
    this.connection.sendMessage("getGcodeFileMetas");
  }
};
</script>-->

<style>
</style>
