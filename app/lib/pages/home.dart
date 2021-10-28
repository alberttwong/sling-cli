import 'package:flutter/material.dart';
import 'dart:io';
import 'dart:convert';

import 'package:sling/core/global.dart';
import 'package:sling/core/task.dart';
import 'package:sling/pages/task/new/type_page.dart';

class HomePage extends StatefulWidget {
  static const routeName = '/home';
  const HomePage({Key? key, required this.title}) : super(key: key);

  final String title;

  @override
  State<HomePage> createState() => _HomePageState();
}

ClipRRect MakeBlueButton(
    BuildContext context, String text, Function()? onPressed) {
  var button = ClipRRect(
    borderRadius: BorderRadius.circular(4),
    child: Stack(
      children: <Widget>[
        Positioned.fill(
          child: Container(
            decoration: const BoxDecoration(
              gradient: LinearGradient(
                colors: <Color>[
                  Color(0xFF0D47A1),
                  Color(0xFF1976D2),
                  Color(0xFF42A5F5),
                ],
              ),
            ),
          ),
        ),
        TextButton(
          style: TextButton.styleFrom(
            padding: const EdgeInsets.all(16.0),
            primary: Colors.white,
            textStyle: const TextStyle(fontSize: 20),
          ),
          onPressed: onPressed,
          child: Text(text),
        ),
      ],
    ),
  );
  return button;
}

class _HomePageState extends State<HomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      // This call to setState tells the Flutter framework that something has
      // changed in this State, which causes it to rerun the build method below
      // so that the display can reflect the updated values. If we changed
      // _counter without calling setState(), then the build method would not be
      // called again, and so nothing would appear to happen.
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    void testSlingAsync() async {
      var result = await Process.start(Global.binPath, ['--version']);
      var stdout = await result.stdout
          .transform(utf8.decoder)
          .transform(const LineSplitter())
          .join();
      var stderr = await result.stderr
          .transform(utf8.decoder)
          .transform(const LineSplitter())
          .join();
      print(stdout);
      print(stderr);
    }

    void testSling() {
      var result = Process.runSync(Global.binPath, ['--version']);
      var version = result.stdout;
      showDialog<String>(
        context: context,
        builder: (BuildContext context) => AlertDialog(
          title: const Text('Success'),
          content: Text(version),
          actions: <Widget>[
            TextButton(
              onPressed: () => Navigator.pop(context, 'OK'),
              child: const Text('OK'),
            ),
          ],
        ),
      );
    }

    return Scaffold(
      // appBar: AppBar(
      //   // Here we take the value from the MyHomePage object that was created by
      //   // the App.build method, and use it to set our appbar title.
      //   title: Text(widget.title),
      // ),
      body: Center(
        // Center is a layout widget. It takes a single child and positions it
        // in the middle of the parent.
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            MakeBlueButton(context, 'New Task', () {
              Navigator.pushNamed(
                context,
                TaskNewType.routeName,
                arguments: Task('12341258536.cha'),
              );
            }),
            const SizedBox(height: 15),
            MakeBlueButton(
                context,
                'Re-Run Task',
                () => showDialog<String>(
                      context: context,
                      builder: (BuildContext context) => AlertDialog(
                        title: const Text('AlertDialog Title'),
                        content: const Text('AlertDialog description'),
                        actions: <Widget>[
                          TextButton(
                            onPressed: () => Navigator.pop(context, 'Cancel'),
                            child: const Text('Cancel'),
                          ),
                          TextButton(
                            onPressed: () => Navigator.pop(context, 'OK'),
                            child: const Text('OK'),
                          ),
                        ],
                      ),
                    )),
            TextButton(
              child: const Text('test sling'),
              onPressed: () => {testSling()},
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headline4,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
