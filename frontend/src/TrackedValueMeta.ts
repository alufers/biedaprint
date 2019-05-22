import Component, { createDecorator } from "vue-class-component";
import Vue from "vue";
import gql from "graphql-tag";
import { TrackedValue, GetTrackedValueByNameQuery, GetTrackedValueByNameWithMetaQuery, GetTrackedValueByNameWithMetaQueryVariables } from "./graphql-models-gen";
import getTrackedValueByNameWithMeta from "../../queries/getTrackedValueByNameWithMeta.graphql";

export default function TrackedValueMeta(tvName: string) {
  return createDecorator((options, key) => {
    @Component
    class TrackedValueMetaDecoratorMixin extends Vue {
      async created() {
        let withLoader = (cbFn: () => Promise<any>) => cbFn();
        if ((this as any).withLoader) {
          withLoader = (this as any).withLoader;
        }
        await withLoader(async () => {
          let tv = await this.$apollo.query<GetTrackedValueByNameWithMetaQuery>({
            variables: <GetTrackedValueByNameWithMetaQueryVariables> {
              name: tvName
            },
            query: getTrackedValueByNameWithMeta
          });
          this[key] = tv.data.trackedValue;
        });
      }
    }
    options.mixins = [...options.mixins, TrackedValueMetaDecoratorMixin];
  });
}
