package server

import (
	"net/http"
	"strings"

	"github.com/flarco/dbio/connection"
	"github.com/flarco/dbio/iop"
	"github.com/flarco/g"
	"github.com/flarco/sling/core/env"
	"github.com/labstack/echo/v4"
)

// RouteName is the name of a route
type RouteName string

const (
	RouteConnections RouteName = "/connections"
	RouteTasks       RouteName = "/tasks"
	RouteExecute     RouteName = "/execute"
	RouteHistory     RouteName = "/history"
	RouteTerminate   RouteName = "/terminate"
	RouteStatus      RouteName = "/status"
	RouteSchemata    RouteName = "/schemata"
	RouteTable       RouteName = "/table"
	RouteWs          RouteName = "/ws"
)

func (r RouteName) String() string {
	return string(r)
}

func addRoutes(e *echo.Echo) {
	e.GET(RouteConnections.String(), GetConnections)
	e.POST(RouteConnections.String(), PostConnections)
	e.DELETE(RouteConnections.String(), DeleteConnections)

	e.GET(RouteTasks.String(), GetTasks)
	e.POST(RouteTasks.String(), PostTasks)
	e.POST(RouteExecute.String(), PostExecute)
	e.GET(RouteHistory.String(), GetHistory)
	e.POST(RouteTerminate.String(), PostTerminate)
	e.GET(RouteStatus.String(), GetStatus)
	e.GET(RouteSchemata.String(), GetSchemata)
}

// Request is the typical request struct
type Request struct {
	Name      string      `json:"name" query:"name"`
	Conn      string      `json:"conn" query:"conn"`
	Database  string      `json:"database" query:"database"`
	Schema    string      `json:"schema" query:"schema"`
	Table     string      `json:"table" query:"table"`
	Procedure string      `json:"procedure" query:"procedure"`
	Data      interface{} `json:"data" query:"data"`
}

// Response is the typical response struct
type Response struct {
	ID    string                 `json:"id"`
	Data  map[string]interface{} `json:"data"`
	Error string                 `json:"error"`
}

func LoadConnection(connID string) (conn env.Conn, err error) {
	// os.Setenv("DBIO_USE_POOL", "TRUE")
	localConns := env.GetLocalConns()
	for _, localConn := range localConns {
		if strings.EqualFold(connID, localConn.Name) {
			conn = localConn
		}
	}
	if conn.Name == "" {
		err = g.Error("could not find connection: %s", connID)
	}
	return
}

// GetConnections (conn_id)
func GetConnections(c echo.Context) (err error) {
	req := map[string]string{}
	if err = c.Bind(&req); err != nil {
		return g.ErrJSON(http.StatusBadRequest, err, "invalid get connections request")
	}
	var conn env.Conn
	var conns []env.Conn

	connID := req["conn_id"]
	test := req["test"]

	if connID != "" {
		conn, err = LoadConnection(connID)
		if err != nil {
			return g.ErrJSON(http.StatusBadRequest, err)
		}
		conns = []env.Conn{conn}
	}

	switch {
	case connID != "" && test != "":
		if err = TestConnection(conn.Connection, "", 15); err != nil {
			return g.ErrJSON(http.StatusBadRequest, err, "error testing connectivity for %s", connID)
		}
		return c.JSON(http.StatusOK, Response{Data: g.M()})
	case connID == "":
		conns = env.GetLocalConns() // return all connections
	}
	return c.JSON(http.StatusOK, Response{Data: g.M("conns", conns)})

}

// PostConnections adds a connection to local profile
func PostConnections(c echo.Context) (err error) {
	req := map[string]interface{}{}
	if err = c.Bind(&req); err != nil {
		return g.ErrJSON(http.StatusBadRequest, err, "invalid set connection request")
	}

	// create connection to test
	conn, err := connection.NewConnectionFromMap(req)
	if err != nil {
		return g.ErrJSON(http.StatusBadRequest, err, "unable to create connection")
	}

	if err = TestConnection(conn, "", 15); err != nil {
		return g.ErrJSON(http.StatusBadRequest, err, "error testing connectivity for %s", conn.Name)
	}

	// save to profile

	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

func DeleteConnections(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// GetTasks (source_conn, target_conn)
func GetTasks(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// PostTasks
func PostTasks(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// PostExecute (task_id)
func PostExecute(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// GetHistory (task_id)
func GetHistory(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// PostTerminate (execution_id)
func PostTerminate(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// GetStatus (execution_id, log=[true,false])
func GetStatus(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Response{Data: g.M()})
}

// GetSchemata (conn_id, level=[database,schema,table,column], database, schema, query, url)
func GetSchemata(c echo.Context) (err error) {
	type Level string
	const LevelDatabase Level = "database"
	const LevelSchema Level = "schema"
	const LevelTable Level = "table"
	const LevelColumn Level = "column"

	req := map[string]string{}
	if err = c.Bind(&req); err != nil {
		return g.ErrJSON(http.StatusBadRequest, err, "invalid get schemata request")
	}

	connID := req["conn_id"]
	level := Level(req["level"]) // database,schema,table,column
	database := req["database"]
	schema := req["schema"]
	query := req["query"]
	url := req["url"]

	if connID == "" {
		return g.ErrJSON(http.StatusBadRequest, g.Error("invalid get schemata request: must provide connection id"))
	} else if level == "" {
		return g.ErrJSON(http.StatusBadRequest, g.Error("invalid get schemata request: must provide level"))
	}

	conn, err := LoadConnection(connID)
	if err != nil {
		return g.ErrJSON(http.StatusBadRequest, err)
	} else if database == "" {
		database = conn.Connection.Info().Database
	}

	if url != "" {
		data, err := DataProcedureFile(conn.Connection, "get_file_columns", url, 15)
		if err != nil {
			return g.ErrJSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, Response{Data: data.ToJSONMap()})
	}

	var data iop.Dataset
	switch level {
	case LevelDatabase:
		data, err = DataProcedureDB(conn.Connection, "get_databases", "", "", "", 15)
		if err != nil {
			return g.ErrJSON(http.StatusInternalServerError, err)
		}
	case LevelSchema:
		data, err = DataProcedureDB(conn.Connection, "get_schemas", database, "", "", 15)
		if err != nil {
			return g.ErrJSON(http.StatusInternalServerError, err)
		}
	case LevelTable:
		data, err = DataProcedureDB(conn.Connection, "get_tables", database, schema, "", 15)
		if err != nil {
			return g.ErrJSON(http.StatusInternalServerError, err)
		}
	case LevelColumn:
		data, err = DataProcedureDB(conn.Connection, "get_sql_columns", database, schema, query, 15)
		if err != nil {
			return g.ErrJSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, Response{Data: data.ToJSONMap()})
}
