import Component, { createDecorator } from "vue-class-component";
import Vue from "vue";
import gql from "graphql-tag";
import {
  GetTrackedValueByNameOnlyValueQuery,
  SubscribeToTrackedValueUpdatedByNameSubscription,
  SubscribeToTrackedValueUpdatedByNameSubscriptionVariables
} from "../graphql-models-gen";
import { QueryResult } from "vue-apollo/types/vue-apollo";

/**
 * TrackedValueSubscription is a decorator which binds to a vue instance property that will be updated every time the tracked value changes and an update event is sent via the subscription.
 * @param tvName the name of the tracked value
 */
export default function TrackedValueSubscription(
  tvName: string | (() => string)
) {
  return createDecorator((options, key) => {
    @Component
    class TrackedValueSubscriptionDecoratorMixin extends Vue {
      async created() {
        let withLoader = (cbFn: () => Promise<any>) => cbFn();
        if ((this as any).withLoader) {
          withLoader = (this as any).withLoader;
        }
        await withLoader(async () => {
          let resultingTvName: string;
          if (typeof tvName === "function") {
            resultingTvName = tvName.bind(this)();
          } else {
            resultingTvName = tvName;
          }
          let tv = await this.$apollo.query<
            GetTrackedValueByNameOnlyValueQuery
          >({
            variables: {
              name: resultingTvName
            },
            fetchPolicy: "network-only",
            query: gql`
              query getTrackedValueByNameOnlyValue($name: String!) {
                trackedValue(name: $name) {
                  value
                }
              }
            `
          });
          (this as any)[key] = tv.data.trackedValue.value;

          // create the real subscription
          let observable = this.$apollo.subscribe<
            QueryResult<SubscribeToTrackedValueUpdatedByNameSubscription>
          >({
            variables: <
              SubscribeToTrackedValueUpdatedByNameSubscriptionVariables
            >{
              name: resultingTvName
            },

            query: gql`
              subscription subscribeToTrackedValueUpdatedByName(
                $name: String!
              ) {
                trackedValueUpdated(name: $name)
              }
            `
          });

          observable.subscribe(val => {
            (this as any)[key] = val.data.trackedValueUpdated;
          });
        });
      }
    }
    if (!options.mixins) {
      options.mixins = [];
    }
    options.mixins = [
      ...options.mixins,
      TrackedValueSubscriptionDecoratorMixin
    ];
  });
}
