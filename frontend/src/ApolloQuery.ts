import Component, { createDecorator } from "vue-class-component";
import Vue from "vue";
import { VueApolloQueryOptions } from 'vue-apollo/types/options';

export default function ApolloQuery<R=any>(opts: VueApolloQueryOptions<R, Vue>) {
  return createDecorator((options, key) => {
    @Component
    class ApolloQueryMixin extends Vue {
      async created() {
          this.$apollo.addSmartQuery(key, opts)
      }
    }
    options.mixins = [...options.mixins, ApolloQueryMixin];
  });
}
