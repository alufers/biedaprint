import Vue from "vue";
import VueApollo from "vue-apollo";
import { createApolloClient } from "vue-cli-plugin-apollo/graphql-client";
import { createUploadLink } from "apollo-upload-client";
// Install the vue plugin
Vue.use(VueApollo);

// Http endpoint
const httpEndpoint =
  process.env.VUE_APP_GRAPHQL_HTTP || "http://localhost:4444/query";

let wsEndpoint = process.env.VUE_APP_GRAPHQL_WS;

if (!wsEndpoint) {
  wsEndpoint = "ws://localhost:4444/query";
}

function websocketRelativeUrl() {
  var l = window.location;
  return (
    (l.protocol === "https:" ? "wss://" : "ws://") +
    l.hostname +
    (l.port != "80" && l.port != "443" ? ":" + l.port : "") +
    "/query"
  );
}

if (process.env.VUE_APP_GRAPHQL_WS_AUTO_RELATIVE) {
  wsEndpoint = websocketRelativeUrl();
}

// Config
const defaultOptions = {
  // You can use `https` for secure connection (recommended in production)
  httpEndpoint,
  // You can use `wss` for secure connection (recommended in production)
  // Use `null` to disable subscriptions
  wsEndpoint,
  // Enable Automatic Query persisting with Apollo Engine
  persisting: false,
  // Use websockets for everything (no HTTP)
  // You need to pass a `wsEndpoint` for this to work
  websocketsOnly: false,
  // Is being rendered on the server?
  ssr: false,

  // Override default apollo link
  // note: don't override httpLink here, specify httpLink options in the
  // httpLinkOptions property of defaultOptions.
  link: createUploadLink({
    uri: httpEndpoint
  })

  // Override default cache
  // cache: myCache

  // Override the way the Authorization header is set
  // getAuth: (tokenName) => ...

  // Additional ApolloClient options
  // apollo: { ... }

  // Client local data (see apollo-link-state)
  // clientState: { resolvers: { ... }, defaults: { ... } }
};

// Call this in the Vue app file
export function createProvider(options = {}) {
  // Create apollo client
  const { apolloClient, wsClient } = createApolloClient({
    ...defaultOptions,
    ...options
  });
  apolloClient.wsClient = wsClient;

  // Create vue apollo provider
  const apolloProvider = new VueApollo({
    defaultClient: apolloClient,
    defaultOptions: {
      $query: {
        // fetchPolicy: 'cache-and-network',
      }
    },
    errorHandler(error) {
      // eslint-disable-next-line no-console
      console.log(
        "%cError",
        "background: red; color: white; padding: 2px 4px; border-radius: 3px; font-weight: bold;",
        error.message
      );
    }
  });

  return apolloProvider;
}
