package server_test

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/flarco/g"
	"github.com/flarco/g/net"
	"github.com/flarco/sling/core/server"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

var (
	srv = server.NewServer("", "")
)

func TestAll(t *testing.T) {
	go srv.Start()
	time.Sleep(1 * time.Second)
	assert.NoError(t, nil)
	// testConnection(t)
	testSchemata(t)
}

func extractData(data1 map[string]interface{}) map[string]interface{} {
	return cast.ToStringMap(data1["data"])
}

func getRequest(route server.RouteName, data1 map[string]interface{}) (data2 map[string]interface{}, err error) {
	headers := map[string]string{"Content-Type": "application/json"}
	vals := url.Values{}
	for k, v := range data1 {
		switch v.(type) {
		case map[string]interface{}:
			v = string(g.MarshalMap(v.(map[string]interface{})))
		}
		val := cast.ToString(v)
		if val == "" {
			continue
		}
		vals.Set(k, val)
	}
	url := g.F("http://localhost:%s%s?%s", srv.Port, route.String(), vals.Encode())
	g.P(url)
	_, respBytes, err := net.ClientDo("GET", url, nil, headers, 5)
	if err != nil {
		err = g.Error(err)
		return
	}
	err = g.Unmarshal(string(respBytes), &data2)
	if err != nil {
		err = g.Error(err)
		return
	}
	return
}

func postRequest(route server.RouteName, data1 map[string]interface{}) (data2 map[string]interface{}, err error) {
	headers := map[string]string{"Content-Type": "application/json"}
	url := g.F("http://localhost:%s%s", srv.Port, route.String())
	g.P(url)
	_, respBytes, err := net.ClientDo("POST", url, strings.NewReader(g.Marshal(data1)), headers)
	if err != nil {
		err = g.Error(err)
		g.Unmarshal(string(respBytes), &data2)
		return
	}

	err = g.Unmarshal(string(respBytes), &data2)
	if err != nil {
		err = g.Error(err)
		return
	}
	return
}

func testConnection(t *testing.T) {
	m := g.M()

	// get connections
	data, err := getRequest(server.RouteConnections, m)
	if !g.AssertNoError(t, err) {
		return
	}
	conns := cast.ToSlice(extractData(data)["conns"])
	assert.Greater(t, len(conns), 1)

	// create connection
	// update connection
	// delete connection
}

func testSchemata(t *testing.T) {
	m := g.M(
		"conn_id", "PG_BIONIC_URL",
		"level", "database",
	)

	// get connections
	respData, err := getRequest(server.RouteSchemata, m)
	if !g.AssertNoError(t, err) {
		return
	}
	data := extractData(respData)
	g.PP(data)

}
