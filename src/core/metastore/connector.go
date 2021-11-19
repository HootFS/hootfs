package metastore

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The metastore will have all sorts of operations on it.
// Should the metastore just become the new virtual file manager...
// Maybe...

const MetaStoreURI = "mongodb+srv://caa8:hootfs@hootfsmetadata.qffjw.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

type MetaStore struct {
	client *mongo.Client
	testDB *mongo.Database
}

// MetaStore will have the following backend structure.
//
// The Database "prod" will hold all data accessed by this MetaStore impl.
// Unless you ask to make a test collection. This goes to the "test" database.
//
// "prod" will be an unordered dictionary with the following keys.
// "namespaces" : ...
// "directories" : ...
// "files" : ...

func NewMetaStore() (*MetaStore, error) {
	clientOptions := options.Client().ApplyURI(MetaStoreURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	return &MetaStore{
		client: client,
		testDB: client.Database("test"),
	}, nil
}

func (ms *MetaStore) TestCollection(collectionName string) *mongo.Collection {
	return ms.testDB.Collection(collectionName)
}

func (ms *MetaStore) Disconnect() error {
	return ms.client.Disconnect(context.TODO())
}
