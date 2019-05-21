import Component, { createDecorator } from "vue-class-component";
import Vue from "vue";
import gql from "graphql-tag";
import { TrackedValue, GetTrackedValueByNameQuery } from "./graphql-models-gen";

export default function TrackedValueMeta(tvName: string) {
  return createDecorator((options, key) => {
    @Component
    class TrackedValueDecoratorMixin extends Vue {
      async created() {
        let tv = await this.$apollo.query<GetTrackedValueByNameQuery>({
          variables: {
            name: tvName
          },
          query: gql`
            query getTrackedValueByName($name: String!) {
              trackedValue(name: $name) {
                name
                unit
                displayType
                plotColor
                value
                lastUpdate
                lastSent
                minUpdateInterval
                history
                maxHistoryLength
              }
            }
          `
        });
        this[key] = tv.data.trackedValue;
      }
    }
    options.mixins = [...options.mixins, TrackedValueDecoratorMixin];
  });
}
