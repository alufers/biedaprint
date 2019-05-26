import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:biedaprint_mobile/widgets/widgets.dart';

const baudRates = [19200, 28800, 38400, 57600, 115200, 2500000];
const connectToSerial = '''
mutation {
  connectToSerial(void: null)
}
''';

const updateSettings = '''
mutation updateSettings(\$newSettings: NewSettings!) {
  updateSettings(settings: \$newSettings) {
    serialPort @include(if: false)  # shitty hack not to return anything [shrug]
  }
}
''';

const disconnectFromSerial = '''
mutation {
  disconnectFromSerial(void: null)
}
''';

const getSettingsAndSerialPorts = '''
query getSettingsAndSerialPorts {
  settings {
    serialPort
    baudRate
    dataBits
    parity
    dataPath
    startupCommand
    scrollbackBufferSize
    temperaturePresets {
      name
      hotendTemperature
      hotbedTemperature
    }
  }
  serialPorts
}
''';

class ConnectionWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.all(20.0),
      child: Container(
        width: double.infinity,
        child: TrackedValueWidget<String>(
          name: 'serialStatus',
          builder: (context, value) {
            final isConnected = value == 'connected';
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: <Widget>[
                isConnected ? Container() : _drawQuickSettings(context),
                _drawConnectionButton(context, isConnected)
              ],
            );
          },
          errorBuilder: (context, errors) => Text(errors.toString()),
          loadingWidget: Text('Loading'),
        ),
      ),
    );
  }

  Widget _drawConnectionButton(BuildContext context, bool isConnected) {
    return Container(
      width: double.infinity,
      child: Mutation(
        options: MutationOptions(
          document: isConnected
              ? disconnectFromSerial
              : connectToSerial, // this is the mutation string you just created
        ),
        builder: (runMutation, _) {
          return RaisedButton(
            onPressed: () => runMutation({}),
            color: isConnected ? Colors.red : Colors.green,
            textColor: Colors.white,
            child: Text(
                isConnected ? 'Disconnect from printer' : 'Connect to printer'),
          );
        },
      ),
    );
  }


  Widget _drawQuickSettings(BuildContext context) {
    // @TODO rewrite it, messy af
    return Query(
                  options: QueryOptions(document: getSettingsAndSerialPorts),
                  builder: (result, {refetch}) {
                    if (result.errors != null) {
                      return Text(result.errors.toString());
                    }

                    if (result.loading) {
                      return Text('Loading');
                    }
                    final serialPorts =
                        List<String>.from(result.data['serialPorts']);

                    var settings = result.data['settings'];
                    return Mutation(
                      options: MutationOptions(
                        document: updateSettings,
                      ),
                      builder: (runmutation, context) => Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('Baudrate',
                                  textAlign: TextAlign.start,
                                  style: TextStyle(
                                      color: Colors.grey, fontSize: 16)),
                              new DropdownButton<int>(
                                value: result.data['settings']['baudRate'],
                                items: baudRates.map((value) {
                                  return new DropdownMenuItem<int>(
                                    value: value,
                                    child: new Text(value.toString()),
                                  );
                                }).toList(),
                                onChanged: (value) {
                                  settings['baudRate'] = value;
                                  runmutation({
                                    'newSettings': settings,
                                  });
                                },
                              ),
                              Text('Serial port',
                                  textAlign: TextAlign.start,
                                  style: TextStyle(
                                      color: Colors.grey, fontSize: 16)),
                              new DropdownButton<String>(
                                value: result.data['settings']['serialPort'],
                                items: serialPorts.map((value) {
                                  return new DropdownMenuItem<String>(
                                    value: value,
                                    child: new Text(value),
                                  );
                                }).toList(),
                                onChanged: (value) async {
                                  settings['serialPort'] = value;
                                  runmutation({
                                    'newSettings': settings,
                                  });
                                },
                              )
                            ],
                          ),
                          onCompleted: (_) => refetch(),
                    );
                  },
                );
  }
}
