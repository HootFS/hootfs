package vfm

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MetaStoreURI = "mongodb+srv://caa8:comp413@cluster0.jmblg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

const (
	C_Machines        = "Machines"
	C_Users           = "Users"
	C_Namespaces      = "Namespaces"
	C_Virtual_Objects = "Virtual_Objects"
)

// NOTE, the below four structs mirror how data will be
// stored in the metastore.
// Each type of struct will have its own collection
// for being stored in. These are seen above with all
// collection names starting with the prefix C_.

type Machine struct {
	MID Machine_ID

	// Potentially more here...
}

type User struct {
	UID User_ID

	// Potentially more here...
}

type Namespace struct {
	NID         Namespace_ID
	RootObjects []VO_ID

	// Users which have access to this Namespace.
	Users []User_ID
}

type VObject struct {
	VOID     VO_ID
	ParentID VO_ID

	// This will be this object's closest parent which is a root of a
	// namespace. If this object has no parent, this will be the NULL
	// VO_ID.
	ClosestRoot VO_ID

	Name string

	// A namespace N will be in this slice if and only if this object
	// is a root object of N.
	namespaces []Namespace_ID

	IsDir bool

	// If this object is a file, this slice will exist, and will
	// contain all machines this file resides on.
	Machines []Machine_ID

	// If this object is a directory, this slice will exist, and
	// will hold all objects contained in this directory.
	SubObjects []VO_ID
}

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

var (
	ErrNotImplemented      = errors.New("Not Implemented!")
	ErrMachineExists       = errors.New("Machine Already Exists!")
	ErrMachineDoesNotExist = errors.New("Machine Does Not Exist!")
)

func (ms *MetaStore) CreateMachine(new_machine Machine_ID) error {
	filter := bson.D{{"mid", new_machine}}

	// Must make sure this is not a repeat machine ID.
	match := ms.DB.Collection(C_Machines).FindOne(context.TODO(), filter)

	if match.Err() == nil {
		return ErrMachineExists
	}

	// Error sending request.
	if match.Err() != mongo.ErrNoDocuments {
		return match.Err()
	}

	_, err := ms.DB.Collection(C_Machines).InsertOne(context.TODO(),
		Machine{MID: new_machine})

	return err
}

func (ms *MetaStore) DeleteMachine(old_machine Machine_ID) error {
	// First we must make sure the machine actually exists.
	filter := bson.D{{"mid", old_machine}}
	res, err := ms.DB.Collection(C_Machines).DeleteOne(context.TODO(), filter)

	// TODO ----------------------------
	// Must make sure to delete this machine number from all File Objects!
	// ---------------------------------

	if res.DeletedCount == 0 {
		return ErrMachineDoesNotExist
	}

	return err
}
