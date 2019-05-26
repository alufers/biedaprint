import 'package:flutter/material.dart';
import 'package:graphql_flutter/graphql_flutter.dart';


import 'widgets/widgets.dart';
import 'dart:io' show Platform;

void main() async {
  final httpLink = HttpLink(
    uri: 'http://192.168.254.112:4444/query',
  ) as Link;
  final websocketLink = WebSocketLink(
    url: 'ws://192.168.254.112:4444/query',
    config: SocketClientConfig(
        autoReconnect: true, inactivityTimeout: Duration(seconds: 600)),
  );

  final link = httpLink.concat(websocketLink);

  final client = ValueNotifier(
    GraphQLClient(
      cache: InMemoryCache(),
      link: link,
    ),
  );
  runApp(
    GraphQLProvider(
      client: client,
      child: MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        // See https://github.com/flutter/flutter/wiki/Desktop-shells#fonts
        fontFamily: 'Roboto',
      ),
      home: MainPage(),
    );
  }
}

class MainPage extends StatefulWidget {
  @override
  _MainPageState createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  int _currentIndex = 0;
  final List<Widget> _children = [
    ConnectionWidget(),
    SystemInformationWidget(),
    GcodeFilesWidget(),
    ConsoleWidget(),
    SystemInformationWidget(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _children[_currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        onTap: onTabTapped,
        currentIndex: _currentIndex,
        items: [
          BottomNavigationBarItem(
            icon: new Icon(Icons.settings_input_hdmi),
            title: new Text('Connection'),
          ),
          BottomNavigationBarItem(
            icon: new Icon(Icons.keyboard),
            title: new Text('Controls'),
          ),
          BottomNavigationBarItem(
            icon: new Icon(Icons.file_download),
            title: new Text('Print'),
          ),
          BottomNavigationBarItem(
            icon: new Icon(Icons.swap_vertical_circle),
            title: new Text('Console'),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.info),
            title: Text('System'),
          )
        ],
      ),
    );
  }

  void onTabTapped(int index) {
    setState(
      () {
        _currentIndex = index;
      },
    );
  }
}
