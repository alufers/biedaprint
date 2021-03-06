import Component, { createDecorator } from "vue-class-component";
import Vue from "vue";
import { VueApolloQueryDefinition } from "vue-apollo/types/options";

export default function ApolloQuery<R = any>(
  opts: VueApolloQueryDefinition<R, Vue>
) {
  return createDecorator((options, key) => {
    @Component
    class ApolloQueryMixin extends Vue {
      async created() {
        this.$apollo.addSmartQuery(key, opts);
      }
    }
    if (!options.mixins) {
      options.mixins = [];
    }
    options.mixins = [...options.mixins, ApolloQueryMixin];
  });
}
