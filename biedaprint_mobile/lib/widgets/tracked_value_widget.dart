import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';

const subscribeToTrackedValueUpdatedByName = '''
subscription subscribeToTrackedValueUpdatedByName(\$trackedValue: String!) {
  trackedValueUpdated(name: \$trackedValue)
}
''';

const getTrackedValueByNameOnlyValue = '''
query getTrackedValueByNameOnlyValue(\$trackedValue: String!) {
  trackedValue(name: \$trackedValue) {
    value
  }
}
''';

typedef Widget BuildChildWidget<T>(BuildContext context, T value);
typedef Widget BuildErrorWidget(
    BuildContext context, List<GraphQLError> errors);

class TrackedValueWidget<T> extends StatelessWidget {
  final String name;
  final BuildChildWidget<T> builder;
  final Widget loadingWidget;
  final BuildErrorWidget errorBuilder;

  TrackedValueWidget(
      {this.name, this.builder, this.loadingWidget, this.errorBuilder});
  @override
  Widget build(BuildContext context) {
    return Query(
        options: QueryOptions(
          document: getTrackedValueByNameOnlyValue, // this is the query string you just created
          variables: {
            'trackedValue': name,
          },
        ),
        builder: (result, {refetch}) {
          if (result.errors != null) {
            return errorBuilder(context, result.errors);
          }

          if (result.loading) {
            return loadingWidget;
          }

          return Subscription(
              'subscribeToTrackedValueUpdatedByName', subscribeToTrackedValueUpdatedByName,
              variables: {'trackedValue': name}, builder: ({
            loading,
            payload,
            dynamic error,
          }) {
            return builder(
                context,
                loading
                    ? result.data['trackedValue']['value']
                    : payload['trackedValueUpdated']);
          });
        });
  }
}
