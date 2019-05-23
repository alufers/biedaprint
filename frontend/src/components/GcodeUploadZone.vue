<template>
  <div>
    <div
      class="box has-text-centered"
      @drop="handleDrop($event)"
      @dragover="handleDragover($event)"
    >
      <i class="fas fa-cloud-upload-alt"></i>
      <p>Drop your gcode files here</p>
      <div class="file-coenterer">
        <div class="file is-primary">
          <label class="file-label">
            <input
              class="file-input"
              type="file"
              name="file"
              @change="uploadFile($event.target.files[0])"
            >
            <span class="file-cta">
              <span class="file-icon">
                <i class="fas fa-upload"></i>
              </span>
              <span class="file-label">Choose a file…</span>
            </span>
          </label>
        </div>
      </div>
    </div>
    <div class="modal" :class="{'is-active': uploadModalOpen}">
      <div class="modal-background"></div>
      <div class="modal-content"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">Gcode upload</p>
          <button class="delete" aria-label="close" @click="uploadModalOpen = false"></button>
        </header>
        <section class="modal-card-body">
          <div class="notification is-success" v-if="isFinished">
            <button class="delete" @click="isFinished = false"></button>
            Upload finished
          </div>
          <progress class="progress" max="100">15%</progress>
        </section>
        <footer class="modal-card-foot">
          <!-- <button class="button">Cancel</button> -->
        </footer>
      </div>
      <!-- <button class="modal-close is-large" aria-label="close"></button> -->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component, { mixins } from "vue-class-component";
import LoadableMixin from "../LoadableMixin";
import { uploadGcode } from "../../../queries/uploadGcode.graphql";
import {
  UploadGcodeMutationVariables,
  UploadGcodeMutation,
  GetGcodeFileMetasQuery
} from "../graphql-models-gen";
import { getGcodeFileMetas } from "../../../queries/getGcodeFileMetas.graphql";

@Component({})
export default class GcodeUploadZone extends mixins(LoadableMixin) {
  uploadModalOpen = false;
  isFinished = false;

  handleDrop(ev) {
    // Prevent default behavior (Prevent file from being opened)
    ev.preventDefault();

    if (ev.dataTransfer.items) {
      // Use DataTransferItemList interface to access the file(s)
      for (var i = 0; i < ev.dataTransfer.items.length; i++) {
        // If dropped items aren't files, reject them
        if (ev.dataTransfer.items[i].kind === "file") {
          var file = ev.dataTransfer.items[i].getAsFile();
          this.uploadFile(file);
          return;
        }
      }
    } else {
      // Use DataTransfer interface to access the file(s)
      for (var i = 0; i < ev.dataTransfer.files.length; i++) {
        this.uploadFile(ev.dataTransfer.files[i]);
        return;
      }
    }
  }
  handleDragover(ev) {
    ev.preventDefault();
  }
  uploadFile(f) {
    this.isFinished = false;
    let formData = new FormData();
    formData.append("file", f);
    this.withLoader(async () => {
      this.uploadModalOpen = true;
      await this.$apollo.mutate<UploadGcodeMutation>({
        mutation: uploadGcode,
        variables: <UploadGcodeMutationVariables>{
          file: f
        },
        update(store, { data: { uploadGcode } }) {
          const data = store.readQuery<GetGcodeFileMetasQuery>({
            query: getGcodeFileMetas
          });
          data.gcodeFileMetas.unshift(uploadGcode);
          store.writeQuery<GetGcodeFileMetasQuery>({
            query: getGcodeFileMetas,
            data
          });
        }
      });
    });
    this.isFinished = true;
  }
}
</script>


<!--<script>
import connectionMixin from "@/connectionMixin";
/* eslint-disable */
export default {
  mixins: [connectionMixin],
  data() {
    return {
      uploadModalOpen: false,
      isFinished: false,
      uploadProgress: 0
    };
  },
  methods: {
   
  }
};
</script>-->

<style scoped>
.box {
  padding: 50px;
}
.file-coenterer {
  margin-top: 15px;
  display: flex;
  justify-content: center;
}
</style>
