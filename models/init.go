package models

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

//MongoConfig - struct definition
type MongoConfig struct {
	MongoHost string
	MongoPort int64
	MongoUser string
	MongoPwd  string
}

//IMongoConnection - interface definition
type IMongoConnection interface {
	GetConnection() (*mongo.Database, context.Context, error)
	GetClient() *mongo.Client
	GetContext() context.Context
	GetDb() *mongo.Database
}

//MongoConnection - struct definition
type MongoConnection struct {
	db          *mongo.Database
	client      *mongo.Client
	mongoConfig *DBConfigs1
}

type IMongoClient interface {
}

//MongoClient - struct definition
type MongoClient struct {
	conn IMongoConnection
}

func NewMongoClient(conn IMongoConnection) IMongoClient {
	return &MongoClient{conn: conn}
}

// InitializeConnection -- connect to mongodb
func initializeConnection(connectionString string, dbname string) (*mongo.Client, *mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) //nolint
	var cliErr error
	clientOptions := options.Client()
	clientOptions.SetReadPreference(readpref.Primary())
	clientOptions.SetServerSelectionTimeout(2 * time.Second)
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	client, cliErr := mongo.Connect(ctx, clientOptions.ApplyURI(connectionString))
	if cliErr != nil {
		return nil, nil, fmt.Errorf("couldn't connect to mongo: %v", cliErr)
	}
	pingErr := client.Ping(ctx, readpref.Primary())
	if pingErr != nil {
		return nil, nil, fmt.Errorf("can't create connection: %v", pingErr)
	}
	db := client.Database(dbname)
	return client, db, nil
}

// GetConnection -- create context and return with db
func (conn *MongoConnection) GetConnection() (*mongo.Database, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
	return conn.db, ctx, nil
}

//GetContext - function definition
func (conn *MongoConnection) GetContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
	return ctx
}

// GetClient -- get client instance
func (conn *MongoConnection) GetClient() *mongo.Client {
	return conn.client
}

// GetDb -- get db instance
func (conn *MongoConnection) GetDb() *mongo.Database {
	return conn.db
}

//ConnectAndMigrate - initializes connection and applies migration
func ConnectAndMigrate(serviceMongoConfig *DBConfigs1, dbName string) (IMongoConnection, error) {
	if dbName == "" {
		return nil, fmt.Errorf("DbName must be supplied")
	}
	mongoConfig := serviceMongoConfig
	// if mongoConfig == nil {
	// 	mongoConfig = InitMongoConfig()
	// }

	var (
		mongoURI  string
		mongoDisp string
	)
	if mongoConfig.User > "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%v", mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port)
		mongoDisp = fmt.Sprintf("mongodb://****:****@%s:%v", mongoConfig.Host, mongoConfig.Port)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%v", mongoConfig.Host, mongoConfig.Port)
		mongoDisp = fmt.Sprintf("mongodb://%s:%v", mongoConfig.Host, mongoConfig.Port)
	}

	fmt.Println("Connecting to mongo database at %s/%s...", mongoDisp, dbName)
	// log.Infof("Connecting to mongo database at %s/%s...", mongoDisp, dbName)
	mongoClient, mongoDb, mongoErr := initializeConnection(mongoURI, dbName)
	if mongoErr != nil {
		return nil, mongoErr
	}

	migrateSchema(mongoClient, dbName)

	return &MongoConnection{
		db:          mongoDb,
		client:      mongoClient,
		mongoConfig: mongoConfig,
	}, nil
}

func migrateSchema(client *mongo.Client, dbName string) {
	fmt.Println("Applying migrations..")
	// log.Info("Applying migrations..")
	mongoDriver, mgrErr := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName, MigrationsCollection: "", TransactionMode: false})
	if mgrErr != nil {
		fmt.Println(mgrErr)
		// log.Fatal(mgrErr)
	}

	mongoMigrate, mInstErr := migrate.NewWithDatabaseInstance("file://resources/migrations", dbName, mongoDriver)
	if mInstErr != nil {
		fmt.Println(mInstErr)
		// log.Fatal(mInstErr)
	}

	migCurVer, _, _ := mongoMigrate.Version()
	migrateErr := mongoMigrate.Up()
	if migrateErr != nil && migrateErr.Error() != "no change" {
		fmt.Println(migrateErr)
		// log.Fatal(migrateErr)
	}

	migNewVer, _, _ := mongoMigrate.Version()
	if migNewVer != migCurVer {
		fmt.Println("Migrated from %v to %v schema version", migCurVer, migNewVer)
		// log.Infof("Migrated from %v to %v schema version", migCurVer, migNewVer)
	}
}
