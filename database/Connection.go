package database

import (
	"fmt"
	"log"
	"sync"

	"slabcode/config"
	"slabcode/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	connection *DatabaseConnection
	tables     []interface{}
	once       sync.Once
)

type DatabaseConnection struct {
	User       string
	Password   string
	Host       string
	Database   string
	Connection *gorm.DB
}

func init() {

	tables = append(tables, []interface{}{
		&models.User{},
		&models.Rol{},
		&models.Baned{},
		&models.Project{},
		&models.Task{},
	}...)

	getAllTables(GetDefaultConnection().Connection)
}

//GetConnection return connection to database
func GetConnection(user, password, database, host string) *DatabaseConnection {
	once.Do(func() {
		connection = &DatabaseConnection{
			User:     user,
			Password: password,
			Database: database,
			Host:     host,
		}
		var err error
		err, connection = connection.InitDatabase()
		if err != nil {
			log.Fatal(err.Error())
		}

	})
	return connection
}

//InitDatabase initialize database whit migrate and seeds
func (db *DatabaseConnection) InitDatabase() (error, *DatabaseConnection) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s",
		db.User,
		db.Password,
		db.Host,
		db.Database,
	)
	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Connection = dbConnection
	if config.Debug() {
		err := dropDatabase(db.Connection)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = migrateDatabase(db, tables...)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil, db
}

//GetDefaultConnection get connection with feault values
func GetDefaultConnection() *DatabaseConnection {
	var (
		credentials = config.GetDefaultDatabaseCredentials()
	)
	return GetConnection(
		credentials.UserName,
		credentials.Password,
		credentials.Database,
		credentials.Host,
	)

}

func migrateDatabase(db *DatabaseConnection, tables ...interface{}) error {

	for _, table := range tables {
		err := db.Connection.AutoMigrate(table)
		if err != nil {
			return err
		}
	}
	password, _ := bcrypt.
		GenerateFromPassword(
			[]byte("admin"),
			bcrypt.MinCost,
		)
	db.Connection.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Administrator)
	db.Connection.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Operator)
	db.Connection.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.User{
		UserName: "admin",
		Email:    "admin@slabcode.com",
		Password: string(password),
		RolID:    models.Administrator.ID,
	})
	return nil
}

func dropDatabase(db *gorm.DB) error {
	names := getAllTables(db)
	migrator := db.Migrator()
	for _, name := range names {
		err := migrator.DropTable(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllTables(db *gorm.DB) []string {
	names := make([]string, 0)
	db.Raw(`SELECT table_name
	FROM information_schema.tables
	WHERE table_schema = 'public'
	ORDER BY table_name;`).Scan(&names)
	return names
}
