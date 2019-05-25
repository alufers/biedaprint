import Vue from "vue";
import Component from "vue-class-component";
import { Provide } from "vue-property-decorator";
import { Alert, AlertType } from "./modules/AlertsModule";

@Component({})
export default class LoadableMixin extends Vue {
  loading = false;

  async withLoader<T = any>(cbFn: () => Promise<T>): Promise<void> {
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
      await cbFn();
      shouldWarn = true;
    } catch (e) {
      this.$store.dispatch("AlertsModule/addAlert", <Alert>{
        content: e.message,
        title: "Error",
        type: AlertType.error
      });
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
