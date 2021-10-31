package env

import (
	"os"
	"sort"
	"strings"

	"github.com/flarco/dbio"
	"github.com/flarco/dbio/connection"
	"github.com/flarco/g"
	"github.com/spf13/cast"
)

var (
	HomeDir         = os.Getenv("SLING_DIR")
	DbNetDir        = os.Getenv("DBNET_DIR")
	HomeDirProfile  = ""
	DbNetDirProfile = ""
)

func init() {
	if HomeDir == "" {
		HomeDir = g.UserHomeDir() + "/sling"
		os.Setenv("SLING_DIR", HomeDir)
	}
	if DbNetDir == "" {
		DbNetDir = g.UserHomeDir() + "/dbnet"
	}
	os.MkdirAll(HomeDir, 0755)

	HomeDirProfile = HomeDir + "/profiles.yaml"
	DbNetDirProfile = DbNetDir + "/.dbnet.yaml"

	os.Setenv("DBIO_PROFILE_PATHS", g.F("%s,%s", HomeDirProfile, DbNetDirProfile))
}

type Conn struct {
	Name        string
	Description string
	Source      string
	Connection  connection.Connection
}

func GetLocalConns() []Conn {
	conns := []Conn{}
	for key, val := range g.KVArrToMap(os.Environ()...) {
		if !strings.Contains(val, ":/") || strings.Contains(val, "{") {
			continue
		}
		conn, err := connection.NewConnectionFromURL(key, val)
		if err != nil {
			e := g.F("could not parse %s: %s", key, g.ErrMsgSimple(err))
			g.Warn(e)
			continue
		}

		if connection.GetTypeNameLong(conn) == "" || conn.Info().Type == dbio.TypeUnknown || conn.Info().Type == dbio.TypeFileHTTP {
			continue
		}
		conns = append(conns, Conn{conn.Info().Name, connection.GetTypeNameLong(conn), "env variable", conn})
	}

	dbtConns, err := connection.ReadDbtConnections()
	if !g.LogError(err) {
		for _, conn := range dbtConns {
			conns = append(conns, Conn{conn.Info().Name, connection.GetTypeNameLong(conn), "dbt profiles", conn})
		}
	}

	if g.PathExists(HomeDirProfile) {
		profileConns, err := connection.ReadConnections(HomeDirProfile)
		if !g.LogError(err) {
			for _, conn := range profileConns {
				conns = append(conns, Conn{conn.Info().Name, connection.GetTypeNameLong(conn), "sling profiles", conn})
			}
		}
	}

	if g.PathExists(DbNetDirProfile) {
		profileConns, err := connection.ReadConnections(DbNetDirProfile)
		if !g.LogError(err) {
			for _, conn := range profileConns {
				conns = append(conns, Conn{conn.Info().Name, connection.GetTypeNameLong(conn), "dbnet yaml", conn})
			}
		}
	}

	sort.Slice(conns, func(i, j int) bool {
		return cast.ToString(conns[i].Name) < cast.ToString(conns[j].Name)
	})
	return conns
}
