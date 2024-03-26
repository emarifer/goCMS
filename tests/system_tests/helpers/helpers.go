package helpers

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/emarifer/gocms/database"
	admin_api "github.com/emarifer/gocms/internal/admin_app/api"
	"github.com/emarifer/gocms/internal/app/api"
	"github.com/emarifer/gocms/internal/repository"
	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

// SEE NOTE BELOW:
func init() {
	os.Setenv("GO_ENV", "testing")
}

func RunDatabaseServer(as settings.AppSettings) {
	pro := CreateTestDatabase(as.DatabaseName)
	engine := sqle.NewDefault(pro)
	engine.Analyzer.Catalog.MySQLDb.AddRootAccount()

	session := memory.NewSession(sql.NewBaseSession(), pro)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase(as.DatabaseName)

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", as.DatabaseHost, as.DatabasePort),
	}
	s, err := server.NewServer(config, engine, memory.NewSessionBuilder(pro), nil)
	if err != nil {
		panic(err)
	}
	if err = s.Start(); err != nil {
		panic(err)
	}
}

func CreateTestDatabase(name string) *memory.DbProvider {
	db := memory.NewDatabase(name)
	db.BaseDatabase.EnablePrimaryKeyIndexes()

	pro := memory.NewDBProvider(db)

	return pro
}

func WaitForDb(ctx context.Context, as settings.AppSettings) (*sqlx.DB, error) {
	// can be done <for range 50> in go 1.22
	for range 50 {
		db, err := database.NewMariaDBConnection(ctx, &as)

		if err == nil {
			return db, nil
		}

		time.Sleep(25 * time.Millisecond)
	}

	return nil, fmt.Errorf("database did not start")
}

// GetAppSettings gets the settings for the http servers
// taking into account a unique port. Very hacky way to
// get a unique port: manually have to pass a new number
// for every test...
// TODO : Find a way to assign a unique port at compile time
func GetAppSettings(appNum int) settings.AppSettings {

	appSettings := settings.AppSettings{
		WebserverPort:    8080,
		DatabaseHost:     "localhost",
		DatabasePort:     3336 + appNum, // initial port
		DatabaseUser:     "root",
		DatabasePassword: "",
		DatabaseName:     "cms_db",
	}

	return appSettings
}

// StartApp configures and starts the application to be
// tested by returning a `*gin.Engine` but without attaching
// an http.Server to it
func StartApp(
	ctx context.Context, as settings.AppSettings, dbConn *sqlx.DB,
) (*gin.Engine, error) {

	repo := repository.New(dbConn)
	serv := service.New(repo)
	h := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(h)
	a := api.New(serv, logger, &as)

	cache := api.MakeCache(4, time.Minute*10)
	e, err := a.Start(
		gin.Default(),
		fmt.Sprintf(":%d", as.WebserverPort),
		&cache,
	)
	return e, err
}

// StartAdminApp configures and starts the administration application
// which will be tested by returning a `*gin.Engine` but without
// attaching an http.Server
func StartAdminApp(
	ctx context.Context, as settings.AppSettings, dbConn *sqlx.DB,
) (*gin.Engine, error) {

	repo := repository.New(dbConn)
	serv := service.New(repo)
	v := validator.New()
	a := admin_api.New(serv, v, &as)

	e, err := a.Start(
		gin.Default(),
		fmt.Sprintf(":%d", as.WebserverPort),
	)
	return e, err
}

/* HOW DO I KNOW I'M RUNNING WITHIN "GO TEST". SEE:
https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test#59444829
*/

/* PROBLEM WITH `memory.NewSessionBuilder`: HOW DO I GET THE MAIN BRANCH WITH GOLANG. SEE:
https://stackoverflow.com/questions/42761820/how-to-get-another-branch-instead-of-default-branch-with-go-get
https://www.youtube.com/live/rmgLKG4kmMw?si=8JeeZohJ5myot4yw&t=2524
*/

/* PROBLEM WITH `ERROR 2002 (HY000): Can't connect to local MySQL server through socket '/var/run/mysqld/mysqld.sock'`. SEE:

IT IS NECESSARY TO INSTALL THE MYSQL CLIENT FOR LINUX:
sudo apt install mysql-client-core-8.0

https://www.dailyrazor.com/blog/cant-connect-to-local-mysql-server-through-socket/

SOLUTION:
sudo mysql -h 127.0.0.1 -u root -P3336
OR:
mariadb -h 127.0.0.1 -u root -P3336
*/

/* EXAMPLE. SEE:
https://github.com/dolthub/go-mysql-server
https://docs.dolthub.com/sql-reference/sql-support

This is an example of how to implement a MySQL server.
After running the example, you may connect to it using the following:

> mysql --host=localhost --port=3306 --user=root mydb --execute="SELECT * FROM mytable;"
+----------+-------------------+-------------------------------+----------------------------+
| name     | email             | phone_numbers                 | created_at                 |
+----------+-------------------+-------------------------------+----------------------------+
| Jane Deo | janedeo@gmail.com | ["556-565-566","777-777-777"] | 2022-11-01 12:00:00.000001 |
| Jane Doe | jane@doe.com      | []                            | 2022-11-01 12:00:00.000001 |
| John Doe | john@doe.com      | ["555-555-555"]               | 2022-11-01 12:00:00.000001 |
| John Doe | johnalt@doe.com   | []                            | 2022-11-01 12:00:00.000001 |
+----------+-------------------+-------------------------------+----------------------------+

The included MySQL client is used in this example, however any MySQL-compatible client will work.
*/
