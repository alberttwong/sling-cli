import 'package:flutter/material.dart';
import 'dart:io';
import 'package:flutter/services.dart' show rootBundle;
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'dart:io' show Platform;

import 'pages/task/new/confirm_page.dart';
import 'pages/task/new/source_files_page.dart';
import 'pages/task/new/source_tables_page.dart';
import 'pages/task/new/target_config_page.dart';
import 'pages/task/status_page.dart';
import 'pages/task/history_page.dart';
import 'pages/task/new/type_page.dart';

import 'core/global.dart';
import 'pages/home.dart';

void main() {
  // init
  _setBinPath();

  // run
  runApp(const SlingApp());
}

void _setBinPath() async {
  // construct binary
  var binName = 'sling-mac';
  if (Platform.isLinux) {
    binName = 'sling-linux';
  } else if (Platform.isWindows) {
    binName = 'sling-win.exe';
  }
  // https://stackoverflow.com/questions/52353764/how-do-i-get-the-assets-file-path-in-flutter/53817467
  var directory = await getApplicationSupportDirectory();
  Global.binPath = join(directory.path, binName);
  var bin = await rootBundle.load('bin/' + binName);
  List<int> bytes =
      bin.buffer.asUint8List(bin.offsetInBytes, bin.lengthInBytes);
  await File(Global.binPath).writeAsBytes(bytes);
  print(Global.binPath);
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
        TaskNewType.routeName: (context) => const TaskNewType(),
        TaskNewSourceFilesPage.routeName: (context) =>
            const TaskNewSourceFilesPage(),
        TaskNewSourceTablesPage.routeName: (context) =>
            const TaskNewSourceTablesPage(),
        TaskNewTargetConfigPage.routeName: (context) =>
            const TaskNewTargetConfigPage(),
        TaskConfirmPage.routeName: (context) => const TaskConfirmPage(),
        TaskStatusPage.routeName: (context) => const TaskStatusPage(),
        TaskHistoryPage.routeName: (context) => const TaskHistoryPage(),
      },
    );
  }
}
