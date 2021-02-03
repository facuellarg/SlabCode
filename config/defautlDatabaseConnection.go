package config

type DatabaseCredentials struct {
	UserName string
	Password string
	Database string
	Host     string
}

//GetDefaultDatabaseCredentials
// get default credentials for database
func GetDefaultDatabaseCredentials() DatabaseCredentials {

	return DatabaseCredentials{
		dbUserName,
		dbPassword,
		dbDatabase,
		dbHost,
	}
}
