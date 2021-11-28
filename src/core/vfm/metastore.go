package vfm

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MetaStoreURI = "mongodb+srv://caa8:hootfs@hootfsmetadata.qffjw.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

// NOTE, below lies the schema of the MongoDB metastore.
//
// Machines Collection:
// Machine_Object:
//  _id	 		Machine_ID
//
// Users Collection:
// User_Object:
//	_id 		User_ID (string)
//  // .. Potentially more info here if needed.
//
// Namespaces Collection:
// Namespace_Object:
//  _id         Namespace_UUID
//  RootObjects []VO_UUID
//  Users		[]User_ID
//
// Virtual_Objects collection:
// Virtual_Object:
//  _id         VO_UUID
//  ParentID    VO_UUID
//
//  //
//	ClosestRoot VO_UUID
//	Name        string
//
//  // A Namespace ID N will be held in this slice if and
//  // only if this object is a root object of N.
//	Namespaces  []Namespace_UUID
//
//	IsDir       bool
//  Machines    []Machine_ID  (exists only if file)
//	SubObjects  []VO_UUID	  (exists only if directory)

const (
	C_Machines        = "Machines"
	C_Users           = "Users"
	C_Namespaces      = "Namespaces"
	C_Virtual_Objects = "Virtual_Objects"
)

type MetaStore struct {
	client *mongo.Client
	DB     *mongo.Database
}

func NewMetaStore(db string) (*MetaStore, error) {
	clientOptions := options.Client().ApplyURI(MetaStoreURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	return &MetaStore{
		client: client,
		DB:     client.Database(db),
	}, nil
}

func (ms *MetaStore) Disconnect() error {
	return ms.client.Disconnect(context.TODO())
}

var NotImplemented = fmt.Errorf("Function Not Implemented!")

func (ms *MetaStore) CreateMachine(new_machine Machine_ID) error {
	// Must make sure this is not a repeat machine ID.

	// If not, we will insert into the collection.

	return NotImplemented
}

func (ms *MetaStore) DeleteMachine(old_machine Machine_ID) error {
	return NotImplemented
}
