import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;

const baseHost = 'localhost:9876';
const Map<String, String> defHeaders = {
  HttpHeaders.contentTypeHeader: 'application/json',
};

/*
connections: GET, POST, PUT, CONNECT
tasks: GET, POST, PUT (source_conn, target_conn)
execute: POST (task_id)
history: GET (task_id)
terminate: POST (execution_id)
status: GET (execution_id, log=[true,false])
schemata: GET (level=[database,schema,table,column], database, schema, table, sql, url)
ws: websocket stream
*/
enum Route { connections, tasks, execute, terminate, status, schemata, ws }

extension ParseToString on Route {
  String route() {
    // ignore: unnecessary_this
    return this.toString().split('.').last;
  }
}

class Response {
  int status;
  Map<String, Object> data;

  Response(this.status, this.data);

  bool isError() {
    return status >= 300 ? true : false;
  }

  String error() {
    if (data.containsKey('error')) {
      return data['error'].toString();
    }
    return '';
  }
}

Future<Response> apiGet(Route route, Map<String, Object> data) async {
  var resp = await http.get(Uri.http(baseHost, route.route(), data),
      headers: defHeaders);
  return Response(resp.statusCode, json.decode(resp.body));
}

Future<Response> apiPost(Route route, Map<String, Object> data) async {
  var resp = await http.post(Uri.http(baseHost, route.route()),
      body: data.cast<String, String>(), headers: defHeaders);
  return Response(resp.statusCode, json.decode(resp.body));
}

Future<Response> apiPut(Route route, Map<String, Object> data) async {
  var resp = await http.put(Uri.http(baseHost, route.route()),
      body: data.cast<String, String>(), headers: defHeaders);
  return Response(resp.statusCode, json.decode(resp.body));
}
