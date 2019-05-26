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

          return Column(
            children: <Widget>[
              RaisedButton(child: Text('Refresh'), onPressed: refetch),
              Text(sysData['AppAlloc']),
              Text(sysData['AppName']),
              Text(sysData['AppNumGC'].toString()),
              Text(sysData['AppSysMemory']),
              Text(sysData['GCCPUFractionPercent']),
              Text(sysData['SystemFreeMemory']),
              Text(sysData['SystemTime']),
              Text(sysData['SystemTotalMemory']),
              Text(sysData['SystemUsedMemory']),
              Text(sysData['SystemUsedMemoryPercent']),
            ],
          );
        },
      ),
    );
  }
}
