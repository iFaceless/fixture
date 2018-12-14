package fixture

import (
	"github.com/xo/dburl"
)

type DatabaseURL struct {
	url *dburl.URL
}

func Parse(rawurl string) (*DatabaseURL, error) {
	u, err := dburl.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// only support mysql currently
	if u.Driver != "mysql" {
		return nil, ErrDriverNotSupported
	}
	dbURL := &DatabaseURL{u}
	return dbURL, nil
}

func (dbURL *DatabaseURL) DBName() string {
	if dbURL.url.Path != "" {
		return dbURL.url.Path[1:]
	}
	return ""
}

func (dbURL *DatabaseURL) Driver() string {
	return dbURL.url.Driver
}

func (dbURL *DatabaseURL) DSN() string {
	return dbURL.url.DSN
}

func (dbURL *DatabaseURL) String() string {
	return dbURL.url.String()
}
