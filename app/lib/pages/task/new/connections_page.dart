import 'package:flutter/material.dart';
import 'package:sling/core/task.dart';
import 'package:sling/pages/task/new/target_config_page.dart';

class TaskNewConnectionsPage extends StatelessWidget {
  static const routeName = '/task/new/connections';
  const TaskNewConnectionsPage({Key? key}) : super(key: key);

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
              TaskNewTargetConfigPage.routeName,
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
              TaskNewTargetConfigPage.routeName,
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
              TaskNewTargetConfigPage.routeName,
              arguments: task,
            );
          },
        ),
      ],
    ));
  }
}
