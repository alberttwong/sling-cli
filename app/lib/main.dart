import 'package:flutter/material.dart';
import 'dart:io';

void main() {
  runApp(const SlingApp());
}

class SlingApp extends StatelessWidget {
  const SlingApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Sling',
      theme: ThemeData(
        // This is the theme of your application.
        primarySwatch: Colors.blue,
      ),
      home: const HomePage(title: 'Flutter Demo Home Page'),
      initialRoute: HomePage.routeName,
      routes: {
        HomePage.routeName: (context) =>
            const HomePage(title: 'Flutter Demo Home Page'),
        NewTaskPage.routeName: (context) => const NewTaskPage(),
        SourceConfigPage.routeName: (context) => const SourceConfigPage(),
        TargetConfigPage.routeName: (context) => const TargetConfigPage(),
        ConfirmPage.routeName: (context) => const ConfirmPage(),
        ExecutionStatusPage.routeName: (context) => const ExecutionStatusPage(),
        HistoryPage.routeName: (context) => const HistoryPage(),
      },
    );
  }
}

class HomePage extends StatefulWidget {
  static const routeName = '/home';
  const HomePage({Key? key, required this.title}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<HomePage> createState() => _HomePageState();
}

class Task {
  final String id;
  late String type;
  late String source;
  // late String message;

  Task(this.id);
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
    void testSling() async {
      var result = await Process.start('ls', ['-l']);
      print(result.stdout);
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
            ClipRRect(
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
                    onPressed: () {
                      Navigator.pushNamed(
                        context,
                        NewTaskPage.routeName,
                        arguments: Task('12341258536.cha'),
                      );
                    },
                    child: const Text('New Task'),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 15),
            ClipRRect(
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
                    onPressed: () => showDialog<String>(
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
                    ),
                    child: const Text('Re-Run Task'),
                  ),
                ],
              ),
            ),
            TextButton(
              child: const Text('test sling'),
              onPressed:  () => {
                testSling()
              }, 
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

class NewTaskPage extends StatelessWidget {
  static const routeName = '/new-task';
  const NewTaskPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final task = ModalRoute.of(context)!.settings.arguments as Task;
    return Center(
      // Center is a layout widget. It takes a single child and positions it
      // in the middle of the parent.
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          ClipRRect(
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
                  onPressed: () {
                    task.type = 'file-db';
                    Navigator.pushNamed(
                      context,
                      SourceConfigPage.routeName,
                      arguments: task,
                    );
                  },
                  child: const Text('File To Database'),
                ),
              ],
            ),
          ),
          const SizedBox(height: 15),
          ClipRRect(
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
                  onPressed: () {
                    task.type = 'db-db';
                    Navigator.pushNamed(
                      context,
                      SourceConfigPage.routeName,
                      arguments: task,
                    );
                  },
                  child: const Text('Database To Database'),
                ),
              ],
            ),
          ),
          const SizedBox(height: 15),
          ClipRRect(
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
                  onPressed: () {
                    task.type = 'db-file';
                    Navigator.pushNamed(
                      context,
                      SourceConfigPage.routeName,
                      arguments: task,
                    );
                  },
                  child: const Text('Database To File'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class SourceConfigPage extends StatelessWidget {
  static const routeName = '/source-config';
  const SourceConfigPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final task = ModalRoute.of(context)!.settings.arguments as Task;
    return Scaffold(
        body: ListView(
      children: <Widget>[
        ListTile(
          leading: const Icon(Icons.map),
          title: const Text('PG_URL'),
          onTap: () {
            task.source = 'PG_URL';
            Navigator.pushNamed(
              context,
              TargetConfigPage.routeName,
              arguments: task,
            );
          },
        ),
        ListTile(
          leading: const Icon(Icons.photo_album),
          title: const Text('SNOWFLAKE'),
          onTap: () {
            task.source = 'SNOWFLAKE';
            Navigator.pushNamed(
              context,
              TargetConfigPage.routeName,
              arguments: task,
            );
          },
        ),
        ListTile(
          leading: const Icon(Icons.phone),
          title: const Text('BIGQUERY'),
          onTap: () {
            task.source = 'BIGQUERY';
            Navigator.pushNamed(
              context,
              TargetConfigPage.routeName,
              arguments: task,
            );
          },
        ),
      ],
    ));
  }
}

class TargetConfigPage extends StatelessWidget {
  static const routeName = '/target-config';
  const TargetConfigPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final task = ModalRoute.of(context)!.settings.arguments as Task;
    return Container();
  }
}

class ConfirmPage extends StatelessWidget {
  static const routeName = '/confirm';
  const ConfirmPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container();
  }
}

class ExecutionStatusPage extends StatelessWidget {
  static const routeName = '/execution';
  const ExecutionStatusPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container();
  }
}

class HistoryPage extends StatelessWidget {
  static const routeName = '/history';
  const HistoryPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container();
  }
}
