import Vue from "vue";
import Component from "vue-class-component";

@Component({})
export default class LoadableMixin extends Vue {
  loading = false;

  async withLoader<T = any>(cbFn: () => Promise<T>): Promise<T> {
    this.loading = true;
    let shouldWarn = false;
    setTimeout(
      () =>
        shouldWarn &&
        console.warn(
          "LoadableMixin.withLoader(...) executed synchronously! Did you forget about an await?"
        )
    );
    try {
      let val = await cbFn();
      shouldWarn = true;
      return val;
    } finally {
      this.loading = false;
    }
  }
  get isLoadingClass() {
    return {
      "is-loading": this.loading
    };
  }
}
