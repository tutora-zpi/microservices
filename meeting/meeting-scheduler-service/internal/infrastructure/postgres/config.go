package postgres

import "fmt"

type PostgresConfig struct {
	User     string
	Password string

	Host string
	Port string

	DatabaseName string
	SSLMode      bool
}

func (this *PostgresConfig) ConnectionString() string {
	sslMode := "disable"
	if this.SSLMode {
		sslMode = "enable"
	}

	return fmt.Sprintf("host = %s user = %s password = %s dbname = %s port = %s sslmode = %s", this.Host, this.User, this.Password, this.DatabaseName, this.Port, sslMode)
}
