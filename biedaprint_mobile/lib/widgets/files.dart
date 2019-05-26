import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';

var getGcodeFileMetas = '''
query getGcodeFileMetas {
  gcodeFileMetas {
    originalName
    gcodeFileName
    uploadDate
    totalLines
    printTime
    filamentUsedMm
    layerIndexes {
      lineNumber
      layerNumber
    }
  }
}
''';

class GcodeFilesWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Query(
        options: QueryOptions(
          document:
              getGcodeFileMetas, // this is the query string you just created
        ),
        builder: (result, {refetch}) {
          if (result.errors != null) {
            return Text(result.errors.toString());
          }

          if (result.loading) {
            return Text('Loading');
          }

          final files = result.data['gcodeFileMetas'];

          return ListView.builder(
            itemCount: files.length,
            itemBuilder: (context, index) {
              final file = files[index];
              final printTime = Duration(seconds: file['printTime'].toInt());
              return ListTile(
                trailing: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                  IconButton(
                    icon: Icon(Icons.info),
                    onPressed: () {},
                  ),
                  IconButton(
                    icon: Icon(Icons.print),
                    onPressed: () {},
                  )
                ]),
                title: Text(file['originalName']),
                subtitle: Text(
                  '${printTime.inMinutes.toString()} minutes; Lines: ${file['totalLines'].toString()}',
                ),
              );
            },
          );
        },
      ),
    );
  }
}
