import 'dart:async';

import 'package:flutter/material.dart';

import 'package:multicast_dns/multicast_dns.dart';

class DiscoveryWidget extends StatefulWidget {
  @override
  _DiscoveryWidgetState createState() => _DiscoveryWidgetState();
}

class _DiscoveryWidgetState extends State<DiscoveryWidget> {
  List<String> messageLog = <String>[];

  @override
  initState() {
    super.initState();
    this.enumerateServices();
  }

  enumerateServices() async {
    String name = "_biedaprint._tcp";
    bool verbose = true;
    final MDnsClient client = MDnsClient();
    await client.start();

    await for (PtrResourceRecord ptr in client
        .lookup<PtrResourceRecord>(ResourceRecordQuery.serverPointer(name))) {
      if (verbose) {
        print(ptr);
      }
      await for (SrvResourceRecord srv in client.lookup<SrvResourceRecord>(
          ResourceRecordQuery.service(ptr.domainName))) {
        if (verbose) {
          print(srv);
        }
        if (verbose) {
          await client
              .lookup<TxtResourceRecord>(
                  ResourceRecordQuery.text(ptr.domainName))
              .forEach(print);
        }
        await for (IPAddressResourceRecord ip
            in client.lookup<IPAddressResourceRecord>(
                ResourceRecordQuery.addressIPv4(srv.target))) {
          if (verbose) {
            print(ip);
          }
          addToLog('Service instance found at '
              '${srv.target}:${srv.port} with ${ip.address}.');
        }
        await for (IPAddressResourceRecord ip
            in client.lookup<IPAddressResourceRecord>(
                ResourceRecordQuery.addressIPv6(srv.target))) {
          if (verbose) {
            print(ip);
          }
          addToLog('Service instance found at '
              '${srv.target}:${srv.port} with ${ip.address}.');
        }
      }
    }
    client.stop();
  }

  void addToLog(String val) {
    setState(() {
      messageLog.add(val);
    });
  }

  @override
  Widget build(BuildContext context) {
    return new MaterialApp(
      home: new Scaffold(
          body: new ListView.builder(
        reverse: true,
        itemCount: messageLog.length,
        itemBuilder: (BuildContext context, int index) {
          return new Text(messageLog[index]);
        },
      )),
    );
  }
}
