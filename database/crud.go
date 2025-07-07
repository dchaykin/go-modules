package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/dchaykin/go-modules/datamodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	insertOne(ctx context.Context, record any) error
	replaceOne(ctx context.Context, filter bson.M, replacement any, allowInsert bool) error

	updateEntity(ctx context.Context, doc datamodel.DomainEntity) error
	updateOne(ctx context.Context, filter bson.M, doc any) error

	aggregate(ctx context.Context, match, group bson.M, result any) error

	findOne(ctx context.Context, filter bson.M, doc any) (bool, error)
	findEntity(ctx context.Context, filter bson.M, doc datamodel.DomainEntity) (bool, error)
	findMany(ctx context.Context, filter bson.M, docList any) error
	findWithOptions(ctx context.Context, filter bson.M, result any, sort bson.D, offset, limit int64) error

	get() *mongo.Collection

	removeOne(ctx context.Context, filter bson.M) error
	removeMany(ctx context.Context, filter bson.M) error

	createIndex(ctx context.Context, mod mongo.IndexModel, opts ...*options.CreateIndexesOptions) error
}

func (c mongoCollection) get() *mongo.Collection {
	return c.collection
}

func (c mongoCollection) createIndex(ctx context.Context, mod mongo.IndexModel, opts ...*options.CreateIndexesOptions) error {
	_, err := c.collection.Indexes().CreateOne(ctx, mod, opts...)
	return err
}

func (c mongoCollection) aggregate(ctx context.Context, match, group bson.M, result any) error {
	pipeline := []bson.M{
		{
			"$match": match,
		},
		{
			"$group": group,
		}}
	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

func (c mongoCollection) updateOne(ctx context.Context, filter bson.M, record any) error {
	_, err := c.collection.UpdateOne(ctx, filter, bson.M{"$set": record})
	if err != nil {
		return err
	}
	return nil
}

func (c mongoCollection) replaceOne(ctx context.Context, filter bson.M, replacement any, allowInsert bool) error {
	var opts *options.ReplaceOptions = nil
	if allowInsert {
		opts = &options.ReplaceOptions{}
		opts.SetUpsert(true)
	}

	if _, err := c.collection.ReplaceOne(ctx, filter, replacement, opts); err != nil {
		return err
	}
	return nil
}

func (c mongoCollection) insertOne(ctx context.Context, record any) error {
	if _, err := c.collection.InsertOne(ctx, record); err != nil {
		return err
	}
	return nil
}

func (c mongoCollection) updateEntity(ctx context.Context, doc datamodel.DomainEntity) error {
	_, err := c.collection.UpdateOne(ctx, bson.M{"entity.uuid": doc.UUID()}, bson.M{"$set": doc})
	if err != nil {
		return err
	}

	return nil
}

func (c mongoCollection) removeOne(ctx context.Context, filter bson.M) error {
	if _, err := c.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (c mongoCollection) removeMany(ctx context.Context, filter bson.M) error {
	if _, err := c.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}

func (c mongoCollection) findEntity(ctx context.Context, filter bson.M, doc datamodel.DomainEntity) (found bool, err error) {
	result := c.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	err = result.Decode(doc)
	if err != nil {
		return false, err
	}
	return true, err
}

func (c mongoCollection) findOne(ctx context.Context, filter bson.M, doc any) (found bool, err error) {
	result := c.collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	err = result.Decode(doc)
	if err != nil {
		return false, err
	}
	return true, err
}

func (c mongoCollection) findWithOptions(ctx context.Context, filter bson.M, result any, sort bson.D, offset, limit int64) error {
	findOpt := options.Find()
	if offset != 0 {
		findOpt.SetSkip(offset)
	}
	if limit != 0 {
		findOpt.SetLimit(limit)
	}
	if sort != nil {
		if findOpt == nil {
			findOpt = &options.FindOptions{}
		}
		findOpt.SetSort(sort)
	}

	cursor, err := c.collection.Find(ctx, filter, findOpt)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

func (c mongoCollection) findMany(ctx context.Context, filter bson.M, result any) error {
	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, result)
	if err == nil && cursor.Err() != nil {
		err = cursor.Err()
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	return nil
}

func (ms mongoSession) Extract(coll Collection, filter bson.M, result *[]any, sort bson.D, offset, limit int64) (totalCount int64, err error) {
	err = mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		totalCount, err = coll.get().CountDocuments(context.Background(), filter)
		if err != nil {
			return err
		}
		return coll.findWithOptions(sc, filter, result, sort, offset, limit)
	})
	return totalCount, err
}

func (ms mongoSession) ReplaceOne(coll Collection, filter bson.M, replacement any, allowInsert bool) (err error) {
	return mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		return coll.replaceOne(sc, filter, replacement, allowInsert)
	})

}

func (ms mongoSession) UpdateOne(coll Collection, filter bson.M, doc any) (err error) {
	return mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		return coll.updateOne(sc, filter, doc)
	})

}

func (ms mongoSession) GetDatabaseNames() ([]string, error) {
	cli, err := getMongoClient()
	if err != nil {
		return nil, err
	}
	return cli.client.ListDatabaseNames(context.Background(), bson.D{})
}

func (ms mongoSession) GetCollection(databaseName, collectionName string) Collection {
	return mongoCollection{
		collection: client.client.Database(databaseName).Collection(collectionName),
	}
}

func (ms mongoSession) InsertOne(coll Collection, record any) error {
	return mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		return coll.insertOne(sc, record)
	})
}

func (ms mongoSession) GetDatabase(name string) *mongo.Database {
	cli, err := getMongoClient()
	if err != nil {
		return nil
	}
	return cli.DB(name)
}

func (ms mongoSession) ReplaceEntityByUUID(doc datamodel.DomainEntity, allowInsert bool) error {
	if doc.UUID() == "" {
		return fmt.Errorf("could not upsert an entity: no uuid has been set")
	}
	coll := ms.GetCollection(doc.DatabaseName(), doc.CollectionName())
	return mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		return coll.replaceOne(sc, bson.M{"entity.uuid": doc.UUID()}, doc, allowInsert)
	})
}

func (ms mongoSession) FindEntity(coll Collection, filter bson.M, doc datamodel.DomainEntity) (found bool, err error) {
	err = mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		found, err = coll.findEntity(sc, filter, doc)
		return err
	})
	return found, err
}

func (ms mongoSession) FindOne(coll Collection, filter bson.M, doc any) (found bool, err error) {
	err = mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		found, err = coll.findOne(sc, filter, doc)
		if err != nil {
			return err
		}
		return nil
	})

	return found, err
}

func (ms mongoSession) FindMany(coll Collection, filter bson.M, docList any) error {
	err := mongo.WithSession(context.Background(), ms.session, func(sc mongo.SessionContext) error {
		return coll.findMany(sc, filter, docList)
	})

	return err
}
