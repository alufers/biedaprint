const path = require("path");

module.exports = {
  chainWebpack: config => {
    config.module
      .rule("gql")
      .include.add(path.resolve(__dirname, "../graphql/queries"));
  }
};
