import Vue from "vue";
import Component from "vue-class-component";

@Component({})
export default class LoadableMixin extends Vue {
  loading = false;

  async withLoader<T = any>(cbFn: () => Promise<T>): Promise<T> {
    this.loading = true;
    try {
      return await cbFn();
    } finally {
      this.loading = false;
    }
  }
}
