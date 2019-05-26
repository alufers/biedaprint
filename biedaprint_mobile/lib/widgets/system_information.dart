import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';

var readSystemInformationQuery = '''
query getSystemInformation {
  systemInformation
}''';

class SystemInformationWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Query(
        options: QueryOptions(
          document:
              readSystemInformationQuery, // this is the query string you just created
        ),
        builder: (result, {refetch}) {
          if (result.errors != null) {
            return Text(result.errors.toString());
          }

          if (result.loading) {
            return Text('Loading');
          }

          final sysData = result.data['systemInformation'];
          List<Widget> systemInfoWidgets = [];
          sysData.forEach((k, v) => {
                systemInfoWidgets.add(Container(
                  alignment: Alignment.centerLeft,
                  margin: const EdgeInsets.all(10.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: <Widget>[Text(k, style: TextStyle(color: Theme.of(context).primaryColor),), Text(v.toString())],
                  ),
                ))
              });
          return Column(
            children: <Widget>[
              RaisedButton(child: Text('Refresh'), onPressed: refetch),
              ...systemInfoWidgets
            ],
          );
        },
      ),
    );
  }
}
