import {
  Module,
  VuexModule,
  MutationAction,
  Mutation
} from "vuex-module-decorators";

export enum AlertType {
  info,
  success,
  error
}

export interface Alert {
  id?: number;
  type: AlertType;
  title: string;
  content: string;
  icon?: string;
}

@Module({ namespaced: true })
export default class AlertsModule extends VuexModule<AlertsModule> {
  alerts: Alert[] = [];

  @MutationAction({ mutate: ["alerts"] })
  async addAlert(alert: Alert) {
    if (typeof this.state === "function") throw new Error("ddd");
    return {
      alerts: [
        ...this.state.alerts,
        { ...alert, id: Math.floor(Math.random() * 10000000) }
      ]
    };
  }

  @MutationAction({ mutate: ["alerts"] })
  async removeAlertById(id: number) {
    if (typeof this.state === "function") throw new Error("ddd");
    return {
      alerts: this.state.alerts.filter(a => a.id !== id)
    };
  }
}
