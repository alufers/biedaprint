import Vue from "vue";
import Component, { createDecorator } from "vue-class-component";
import getTrackedValueByNameWithMeta from "../../graphql/queries/getTrackedValueByNameWithMeta.graphql";
import {
  GetTrackedValueByNameWithMetaQuery,
  GetTrackedValueByNameWithMetaQueryVariables
} from "./graphql-models-gen";

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
          let tv = await this.$apollo.query<GetTrackedValueByNameWithMetaQuery>(
            {
              variables: <GetTrackedValueByNameWithMetaQueryVariables>{
                name: tvName
              },
              query: getTrackedValueByNameWithMeta
            }
          );
          (this as any)[key] = tv.data.trackedValue;
        });
      }
    }
    if (!options.mixins) { // create an array if it doesn't exixt so that typescript won't be angry when spreading the array
      options.mixins = []; 
    }
    options.mixins = [...options.mixins, TrackedValueMetaDecoratorMixin];
  });
}
