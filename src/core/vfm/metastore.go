package vfm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MetaStoreURI = "mongodb+srv://caa8:comp413@cluster0.jmblg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

// NOTE, the below structs mirror how data will be
// stored in the metastore.
// Each type of struct will have its own collection
// for being stored in.

type Machine struct {
	MID Machine_ID

	// Potentially more here...
}

type User struct {
	UID User_ID

	// Potentially more here...
}

// Namespaces will mirror the Namespace struct found in
// virtual_file_manager.go

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
	Namespaces []Namespace_ID

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
	db     *mongo.Database
}

func (ms *MetaStore) machines() *mongo.Collection {
	return ms.db.Collection("machines")
}

func (ms *MetaStore) users() *mongo.Collection {
	return ms.db.Collection("users")
}

func (ms *MetaStore) namespaces() *mongo.Collection {
	return ms.db.Collection("namespaces")
}

func (ms *MetaStore) vobjects() *mongo.Collection {
	return ms.db.Collection("vobjects")
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
		db:     client.Database(db),
	}, nil
}

func (ms *MetaStore) Disconnect() error {
	return ms.client.Disconnect(context.TODO())
}

func (ms *MetaStore) CheckMachine(machine Machine_ID, found error,
	not_found error) error {
	return checkHelper(ms.machines().FindOne(context.TODO(),
		bson.M{"mid": machine}), found, not_found)
}

func (ms *MetaStore) CheckUser(user User_ID, found error,
	not_found error) error {
	return checkHelper(ms.users().FindOne(context.TODO(), bson.M{"uid": user}),
		found, not_found)
}

func (ms *MetaStore) CheckNamespace(nsid Namespace_ID, member User_ID,
	found error, not_found error) error {
	return checkHelper(ms.namespaces().FindOne(context.TODO(),
		bson.M{"nsid": nsid, "users": member}), found, not_found)
}

func checkHelper(res *mongo.SingleResult, found error, not_found error) error {
	if res.Err() == nil {
		return found
	}

	if res.Err() == mongo.ErrNoDocuments {
		return not_found
	}

	return res.Err()
}

var (
	ErrNotImplemented      = errors.New("Not implemented!")
	ErrMachineExists       = errors.New("Machine already exists!")
	ErrMachineDoesNotExist = errors.New("Machine does not exist!")
	ErrUserExists          = errors.New("User already exists!")
	ErrUserDoesNotExist    = errors.New("User does not exist!")
	ErrNoAccess            = errors.New("Unable to access namespace!")
	ErrAccess              = errors.New("Namespace is accessible!")
)

// Virtual File Manager Interface Methods Below .............

func (ms *MetaStore) CreateMachine(new_machine Machine_ID) error {
	err := ms.machines().FindOne(context.TODO(),
		bson.M{"mid": new_machine}).Err()

	if err == nil {
		return ErrMachineExists
	}

	if err != mongo.ErrNoDocuments {
		return err
	}

	_, err = ms.machines().InsertOne(context.TODO(),
		Machine{MID: new_machine})

	return err
}

func (ms *MetaStore) DeleteMachine(old_machine Machine_ID) error {
	// First we must make sure the machine actually exists.
	filter := bson.M{"mid": old_machine}
	res, err := ms.machines().DeleteOne(context.TODO(), filter)

	if res.DeletedCount == 0 {
		return ErrMachineDoesNotExist
	}

	// TODO ----------------------------
	// Must make sure to delete this machine number from all File Objects!
	// ---------------------------------

	return err
}

func (ms *MetaStore) CreateUser(new_user User_ID) error {
	if err := ms.CheckUser(new_user, ErrUserExists, nil); err != nil {
		return err
	}

	_, err := ms.users().InsertOne(context.TODO(), User{UID: new_user})

	return err
}

func (ms *MetaStore) DeleteUser(old_user User_ID) error {
	filter := bson.M{"uid": old_user}
	res, err := ms.users().DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrUserDoesNotExist
	}

	updateRequest := bson.M{"$pull": bson.M{"users": old_user}}

	// Delete user from all namespaces.
	_, err = ms.namespaces().UpdateMany(context.TODO(),
		bson.D{}, updateRequest)

	return err

	// TODO -------------
	// NOTE, we will also need some serverside
	// script which deletes all namespaces which have no
	// users. This could be done here, but might be
	// slow.
	// Deleting a Namespace requires some work...
	// ------------------
}

// TODO -------------- Get File Location
// TODO -------------- Set File Location

func (ms *MetaStore) CreateNamespace(name string,
	member User_ID) (Namespace_ID, error) {
	if err := ms.CheckUser(member, nil, ErrUserDoesNotExist); err != nil {
		return Nil_Namespace_ID, err
	}

	// After confirming the user exists, we must create
	// a unique Namespace_ID for the new namespace.

	var nsid Namespace_ID
	for {
		nsid = Namespace_ID(uuid.New())

		// We need a non nil ID.
		if nsid == Nil_Namespace_ID {
			continue
		}

		// Now check if the ID exists already.
		err := ms.namespaces().FindOne(context.TODO(),
			bson.M{"nsid": nsid}).Err()

		if err == mongo.ErrNoDocuments {
			break // Unique ID has been found.
		}

		if err != nil {
			return Nil_Namespace_ID, err
		}

		// Otherwise, err == nil... Non-Unique ID was made.
	}

	namespace := Namespace{
		NSID:  nsid,
		Name:  name,
		Users: []User_ID{member},
	}

	// Finally, upload the Namespace to Mongo.
	_, err := ms.namespaces().InsertOne(context.TODO(), namespace)

	if err != nil {
		return Nil_Namespace_ID, err
	}

	return nsid, nil
}

func (ms *MetaStore) DeleteNamespace(nsid Namespace_ID, member User_ID) error {
	// First, does the given user exist.
	err := ms.users().FindOne(context.TODO(),
		bson.M{"uid": member}).Err()

	if err == mongo.ErrNoDocuments {
		return ErrUserDoesNotExist
	}

	if err != nil {
		return err
	}

	// Check if the namespace exists and if the user has access.
	// NOTE, the user will not be told if a namespace does or does
	// not exist.
	filter := bson.M{"nsid": nsid, "users": member}
	res, err := ms.namespaces().DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNoAccess
	}

	return nil

	// TODO -----------------------
	// Add removal of tags and garbage colleciton
	// of file objects.
	// ----------------------------
}

// func (ms *MetaStore) AddUserToNamespace(nsid Namespace_ID, recruiter User_ID,
// 	recruit User_ID) error {
// 	// First off, does the recruiter exist?
// 	err :=
// }
