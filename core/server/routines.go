package server

import (
	"github.com/flarco/dbio/connection"
	"github.com/flarco/dbio/saas"

	"github.com/flarco/dbio/filesys"

	"github.com/flarco/dbio/database"
	"github.com/flarco/dbio/iop"
	"github.com/flarco/g"
)

// TestConnection tests a connection for access.
// if `permsSchema` is provided, will check read/write permisssions on that schema
func TestConnection(c connection.Connection, testObject string, timeOut ...int) (err error) {
	if c.Info().Type.IsDb() {
		dbConn, err := database.NewConn(c.URL(), g.MapToKVArr(c.DataS())...)
		if err != nil {
			return g.Error(err, "could not initialize database connection with provided credentials/url.")
		}
		defer dbConn.Close()
		err = dbConn.Connect(timeOut...)
		if err != nil {
			return g.Error(err, "could not connect with provided credentials/url")
		}

		if testObject != "" {
			err = database.TestPermissions(dbConn, testObject)
			if err != nil {
				return g.Error(err, "could not successfully test permissions")
			}
		}
	} else if c.Info().Type.IsFile() {
		fs, err := filesys.NewFileSysClientFromURL(c.URL(), g.MapToKVArr(c.DataS())...)
		if err != nil {
			return g.Error(err, "could not initialize file connection with provided credentials/url.")
		}
		_, err = fs.List(c.URL())
		if err != nil {
			return g.Error(err, "could not initialize file connection with provided credentials/url.")
		}

		if testObject != "" {
			err = filesys.TestFsPermissions(fs, testObject)
			if err != nil {
				return g.Error(err, "could not successfully test permissions")
			}
		}
	} else if c.Info().Type.IsAPI() {
		api, err := saas.NewAPIClientFromConn(c)
		if err != nil {
			return g.Error(err, "could not initialize api connection with provided credentials/url.")
		}
		api.Init()
		// TODO: add api.Ping()

	} else {
		return g.Error("could not test conn type %s", connection.GetTypeNameLong(c))
	}
	return
}

// DataProcedureDB is a procedure / routine process
func DataProcedureDB(dConn connection.Connection, procedure, database, schema, sql string, timeOut ...int) (data iop.Dataset, err error) {

	dbConn, err := dConn.ConnSetDatabase(database).AsDatabase()
	if err != nil {
		err = g.Error(err, "could not initialize database connection with provided credentials/url.")
		return data, err
	}
	defer dbConn.Close()
	err = dbConn.Connect(timeOut...)
	if err != nil {
		err = g.Error(err, "could not connect with provided credentials/url")
		return data, err
	}

	switch procedure {
	case "get_conn_spec":
		data = dbConn.Template().ToData()

	case "get_databases":
		data, err = dbConn.GetDatabases()
		if err != nil {
			err = g.Error(err, "could not get databases")
			return data, err
		}

	case "get_schemas":
		data, err = dbConn.GetSchemas()
		if err != nil {
			err = g.Error(err, "could not get schemas")
			return data, err
		}

	case "get_schema_objects":
		schemata, err := dbConn.GetSchemata(schema, "")
		if err != nil {
			err = g.Error(err, "could not get tables for schema: "+schema)
			return data, err
		}
		for _, s := range schemata.Database().Schemas {
			data = s.ToData()
		}
		if len(data.Rows) == 0 {
			err = g.Error("no schemata found for %s", schema)
		}

	case "get_tables":
		data, err = dbConn.GetTables(schema)
		if err != nil {
			err = g.Error(err, "could not get tables for schema: "+schema)
			return data, err
		}

	case "get_views":
		data, err = dbConn.GetViews(schema)
		if err != nil {
			err = g.Error(err, "could not get views for schema: "+schema)
			return data, err
		}
	case "get_sql_columns":
		sql = g.R(sql, "upsert_where_cond", "1=0")
		columns, err := dbConn.GetSQLColumns(sql)
		if err != nil {
			err = g.Error(err, "Could not get columns")
			return data, err
		}
		data = iop.NewDataset(iop.NewColumnsFromFields("column_id", "column_name", "column_type"))
		for _, col := range columns {
			row := []interface{}{col.Position, col.Name, col.Type}
			data.Rows = append(data.Rows, row)
		}
	}
	return
}

// DataProcedureFile is a procedure / routine process for file conns
func DataProcedureFile(dConn connection.Connection, procedure, path string, timeOut ...int) (data iop.Dataset, err error) {
	fileConn, err := dConn.AsFile()
	if err != nil {
		err = g.Error(err, "could not initialize file connection with provided credentials/url.")
		return data, err
	}

	switch procedure {
	case "get_list":
		paths, err := fileConn.List(path)
		if err != nil {
			err = g.Error(err, "could not list path: "+path)
			return data, err
		}
		data = iop.NewDataset(iop.NewColumnsFromFields("name", "is_dir"))
		for _, path := range paths {
			row := []interface{}{path, nil}
			data.Rows = append(data.Rows, row)
		}
	case "get_list_recursive":
		paths, err := fileConn.ListRecursive(path)
		if err != nil {
			err = g.Error(err, "could not list path: "+path)
			return data, err
		}
		data = iop.NewDataset(iop.NewColumnsFromFields("name", "is_dir"))
		for _, path := range paths {
			row := []interface{}{path, nil}
			data.Rows = append(data.Rows, row)
		}
	case "get_file_columns":
		df, err := fileConn.ReadDataflow(path, 100)
		if err != nil {
			err = g.Error(err, "Could not get columns for: "+path)
			return data, err
		}

		for range df.StreamCh {
			// get columns from all files to see if columns match
		}
		if err := df.Context.Err(); err != nil {
			return data, g.Error(err, "error getting columns")
		}

		df.SetEmpty()
		df.Close()
		data = iop.NewDataset(iop.NewColumnsFromFields("column_id", "column_name", "column_type"))
		for _, col := range df.Columns {
			row := []interface{}{col.Position, col.Name, col.Type}
			data.Rows = append(data.Rows, row)
		}
	}
	return
}

// DataQueryDb is a custom SQL query execution
func DataQueryDb(dConn connection.Connection, query string, limit int, timeOut ...int) (data iop.Dataset, err error) {

	if !dConn.Info().Type.IsDb() {
		err = g.Error("cannot query connection type: %s", dConn.Info().Type)
		return data, err
	}

	dbConn, err := database.NewConn(dConn.URL(), g.MapToKVArr(dConn.DataS())...)
	if err != nil {
		err = g.Error(err, "could not initialize database connection with provided credentials/url.")
		return
	}
	defer dbConn.Close()
	err = dbConn.Connect(timeOut...)
	if err != nil {
		err = g.Error(err, "could not connect with provided credentials/url")
		return
	}

	data, err = dbConn.Query(query, limit)
	if err != nil {
		err = g.Error(err, "could not execute query")
		return
	}

	return
}
