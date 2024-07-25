package resources

import (
	"database/sql"
	"github.com/HugoWw/x_apiserver/pkg/client/httpclient"
)

type RestOption struct {
	HttpClient *httpclient.RTClient
	Db         *sql.DB
}
