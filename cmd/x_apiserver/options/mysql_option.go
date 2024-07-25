package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

type MysqlClientOptions struct {
	Password        string
	UserName        string
	Host            string
	Port            int
	ConnMaxLifetime int //in seconds
	MaxIdleConns    int
	MaxOpenConns    int
	Database        string
}

func NewMysqlClientOptions() *MysqlClientOptions {
	return &MysqlClientOptions{
		Port:            3306,
		ConnMaxLifetime: 8 * 3600,
		MaxOpenConns:    100,
		MaxIdleConns:    50,
	}
}

func (m *MysqlClientOptions) Valid() []error {

	errors := []error{}

	if result := net.ParseIP(m.Host); result == nil {
		errors = append(errors, fmt.Errorf("--mysql-host invalid ip address: %s", m.Host))
	}

	if len(m.UserName) == 0 {
		errors = append(errors, fmt.Errorf("--mysql-user can't be empty"))
	}

	if len(m.Password) == 0 {
		errors = append(errors, fmt.Errorf("--mysql-password can't be empty"))
	}

	if m.MaxOpenConns < 0 {
		errors = append(errors, fmt.Errorf("--mysql-max-open-conns can't be less than 0"))
	}

	if m.MaxIdleConns < 0 {
		errors = append(errors, fmt.Errorf("--mysql-max-idle can't be less than 0"))
	}

	if m.ConnMaxLifetime < 0 {
		errors = append(errors, fmt.Errorf("--mysql-max-lifetime can't be less than 0"))
	}

	if len(m.Database) == 0 {
		errors = append(errors, fmt.Errorf("--mysql-database can't be empty"))
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (m *MysqlClientOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&m.Host, "mysql-host", m.Host, "set mysql host address.")
	fs.IntVar(&m.Port, "mysql-port", m.Port, "set mysql connect port. default is 3306")
	fs.StringVar(&m.Password, "mysql-password", m.Password, "set mysql connect password.")
	fs.StringVar(&m.UserName, "mysql-user", m.UserName, "set mysql connect username.")
	fs.IntVar(&m.ConnMaxLifetime, "mysql-max-lifetime", m.ConnMaxLifetime, "set mysql connect max lifetime. default is 8*3600s")
	fs.IntVar(&m.MaxOpenConns, "mysql-max-open-conn", m.MaxOpenConns, "set mysql max open connections. default is 100")
	fs.IntVar(&m.MaxIdleConns, "mysql-max-idle", m.MaxIdleConns, "set mysql max idle connection. default is 10")
	fs.StringVar(&m.Database, "mysql-database", m.Database, "set mysql database name.")
}
