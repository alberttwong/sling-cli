import 'package:flutter/material.dart';
import 'package:sling/core/task.dart';
import 'package:sling/pages/task/new/connections_page.dart';
import 'package:sling/helpers/helpers.dart';

class TaskNewType extends StatelessWidget {
  static const routeName = '/task/new/type';
  const TaskNewType({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final task = ModalRoute.of(context)!.settings.arguments as Task;
    return Center(
      // Center is a layout widget. It takes a single child and positions it
      // in the middle of the parent.
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          makeBlueButton(context, 'File To Database', fontSize: 30,
              onPressed: () {
            task.type = 'file-db';
            Navigator.pushNamed(
              context,
              TaskNewConnectionsPage.routeName,
              arguments: task,
            );
          }),
          const SizedBox(height: 15),
          makeBlueButton(context, 'Database To Database', fontSize: 30,
              onPressed: () {
            task.type = 'db-db';
            Navigator.pushNamed(
              context,
              TaskNewConnectionsPage.routeName,
              arguments: task,
            );
          }),
          const SizedBox(height: 15),
          makeBlueButton(context, 'Database To File', fontSize: 30,
              onPressed: () {
            task.type = 'db-file';
            Navigator.pushNamed(
              context,
              TaskNewConnectionsPage.routeName,
              arguments: task,
            );
          }),
        ],
      ),
    );
  }
}
