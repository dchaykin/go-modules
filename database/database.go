package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const MAX_CONNECT_ATTEMPTS = 3

type OnReadDomainEntity func(object datamodel.DomainEntity) error

func getMongoHost() string {
	return os.Getenv("MONGOHOST")
}

type mongoClient struct {
	client *mongo.Client
}

var client mongoClient

type DatabaseSession interface {
	GetDatabase(name string) *mongo.Database
	InsertOne(coll Collection, record interface{}) error
	ReplaceOne(coll Collection, filter bson.M, replacement interface{}, allowInsert bool) error
	UpdateOne(coll Collection, filter bson.M, replacement interface{}) error
	FindEntity(coll Collection, filter bson.M, doc datamodel.DomainEntity) (bool, error)
	FindOne(coll Collection, filter bson.M, doc interface{}) (bool, error)
	FindMany(coll Collection, filter bson.M, list interface{}) error
	Extract(coll Collection, filter bson.M, result *[]interface{}, sort bson.D, offset, limit int64) (int64, error)
	Aggregate(databaseName, collectionName string, match, group bson.M, result interface{}) error
	GetCollection(databaseName, collectionName string) Collection
	GetDatabaseNames() ([]string, error)
	GetCollectionNames(dbName string) ([]string, error)
	GetEntityByUUID(uuid string, requestedObject datamodel.DomainEntity) (bool, error)
	InsertEntity(entity datamodel.DomainEntity) error
	UpdateEntityByUUID(updatedObject datamodel.DomainEntity) error
	SaveEntityToHistory(entity datamodel.DomainEntity) error
	ReplaceEntityByUUID(entity datamodel.DomainEntity, allowInsert bool) error
	RemoveOne(collection Collection, selector bson.M) error
	RemoveEntity(entity datamodel.DomainEntity) error
	CreateIndex(c Collection, mod mongo.IndexModel, opts ...*options.CreateIndexesOptions) error
	Close() error
}

type mongoSession struct {
	session mongo.Session
}

func (ms mongoSession) Error() string {
	return ""
}

type mongoCollection struct {
	collection *mongo.Collection
}

func getMongoClient() (*mongoClient, error) {
	var err error
	if !client.isConnected() {
		if getMongoHost() == "" {
			return nil, fmt.Errorf("environment variable MONGOHOST is not set. Could not establish a mongo connection")
		}
		err = client.connect()
		if err != nil {
			return nil, fmt.Errorf("could not establish a connection to the mongo server: %v", err)
		}
	}
	return &client, nil
}

func (mc mongoClient) DB(name string) *mongo.Database {
	return mc.client.Database(name)
}

func (mc *mongoClient) connect() error {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      os.Getenv("MONGO_USERNAME"),
		Password:      os.Getenv("MONGO_PASSWORD"),
	}
	params := ""
	if os.Getenv("MONGO_WITH_TLS") == "true" {
		params = "tls=true"
	}
	if params != "" {
		params = "?" + params
	}

	connString := fmt.Sprintf("mongodb://%s/%s", os.Getenv("MONGOHOST"), params)
	clientOpts := options.Client().ApplyURI(connString).SetAuth(credential)

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	mc.client, err = mongo.Connect(ctx, clientOpts)
	if cancel != nil {
		defer cancel()
	}

	if err != nil {
		return err
	}
	return mc.ping(ctx)
}

func (mc *mongoClient) ping(ctx context.Context) error {
	return mc.client.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) isConnected() bool {
	if mc.client != nil {
		if err := mc.ping(context.Background()); err != nil {
			log.WrapError(err)
			return false
		}
		return true
	}
	return false
}

func (mc *mongoClient) Disconnect() error {
	return mc.client.Disconnect(context.Background())
}

func OpenSession() (DatabaseSession, error) {
	cli, err := getMongoClient()
	if err != nil {
		return nil, err
	}

	result := mongoSession{}
	if result.session, err = cli.client.StartSession(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ms *mongoSession) Close() error {
	if ms.session == nil {
		return fmt.Errorf("could not close an empty session")
	}
	ms.session.EndSession(context.Background())
	return nil
}

func HasMongoAccess() bool {

	s, err := OpenSession()
	if err != nil {
		log.WrapError(err)
		return false
	}

	defer func() {
		if err := s.Close(); err != nil {
			log.WrapError(err)
		}
	}()

	return true
}

func (mc mongoClient) GetCollection(dbName string, collectionName string) (collection Collection, err error) {
	db := mc.DB(dbName)
	if db == nil {
		return nil, fmt.Errorf("could not connect to the database %s", dbName)
	}
	return mc.getCollectionByName(db, collectionName)
}

func (mc mongoClient) getCollectionHistory(dataObject datamodel.DomainEntity) (collection Collection, err error) {
	dbName := dataObject.DatabaseName()
	db := mc.DB(dbName)
	if db == nil {
		return nil, fmt.Errorf("could not connect to the database %s", dbName)
	}
	return mc.getCollectionByName(db, dataObject.CollectionName()+"-history")
}

func (mc mongoClient) getCollectionByName(db *mongo.Database, collectionName string) (result Collection, err error) {
	collection := db.Collection(collectionName)
	if collection == nil {
		return nil, fmt.Errorf("could not connect to the collection %s.%s", db.Name(), collectionName)
	}
	return mongoCollection{collection: collection}, nil
}

func (ms mongoSession) Aggregate(dbName, collName string, match, group bson.M, result interface{}) error {
	collection, err := client.GetCollection(dbName, collName)
	if err != nil {
		return err
	}
	return collection.aggregate(context.Background(), match, group, result)
}

func (ms mongoSession) GetEntityByUUID(uuid string, requestedObject datamodel.DomainEntity) (bool, error) {
	if uuid == "" {
		return false, fmt.Errorf("GetObjectByUUID failed. Got an empty UID")
	}

	collection, err := client.GetCollection(requestedObject.DatabaseName(), requestedObject.CollectionName())
	if err != nil {
		return false, err
	}

	found, err := collection.findOne(context.Background(), bson.M{"entity.uuid": uuid}, requestedObject)
	if err != nil {
		return false, fmt.Errorf("GetObjectByRefNo failed. Could not create a query for %v: %v", requestedObject, err)
	}

	return found, err
}

func (ms mongoSession) UpdateEntityByUUID(updatedObject datamodel.DomainEntity) error {
	if updatedObject.UUID() == "" {
		return fmt.Errorf("UpdateEntityByUID failed. Got an empty UID in %v", updatedObject)
	}

	collection, err := client.GetCollection(updatedObject.DatabaseName(), updatedObject.CollectionName())
	if err != nil {
		return err
	}

	err = collection.updateEntity(context.Background(), updatedObject)
	return err
}

func (ms mongoSession) SaveEntityToHistory(entity datamodel.DomainEntity) error {
	collection, err := client.getCollectionHistory(entity)
	if err != nil {
		return err
	}

	err = collection.replaceOne(context.Background(), bson.M{"entity.uuid": entity.UUID()}, entity, true)

	return err
}

func (ms mongoSession) CreateIndex(c Collection, mod mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	return c.createIndex(context.Background(), mod, opts...)
}

func (ms mongoSession) RemoveOne(coll Collection, selector bson.M) error {
	return coll.removeOne(context.Background(), selector)
}

func (ms mongoSession) RemoveEntity(entity datamodel.DomainEntity) error {
	err := ms.SaveEntityToHistory(entity)
	if err != nil {
		return err
	}

	collection, err := client.GetCollection(entity.DatabaseName(), entity.CollectionName())
	if err != nil {
		return err
	}

	selector := bson.M{"entity.uuid": entity.UUID()}

	err = collection.removeOne(context.Background(), selector)

	return err
}

func (ms mongoSession) InsertEntity(entity datamodel.DomainEntity) error {
	if entity.UUID() == "" {
		return fmt.Errorf("cannot insert an entity with an empty UID. Entity: %v", entity)
	}
	collection, err := client.GetCollection(entity.DatabaseName(), entity.CollectionName())
	if err != nil {
		return err
	}

	err = collection.insertOne(context.Background(), entity)

	return err
}

func (ms mongoSession) GetCollectionNames(dbName string) ([]string, error) {
	db := client.DB(dbName)
	if db == nil {
		return nil, fmt.Errorf("could not connect to the database %s", dbName)
	}
	return db.ListCollectionNames(context.Background(), bson.D{})
}
