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

// Confirm the given virtual object exists.
func (ms *MetaStore) CheckVObject(void VO_ID, found error,
	not_found error) error {
	return checkHelper(ms.vobjects().FindOne(context.TODO(),
		bson.M{"void": void}), found, not_found)
}

// Confirm a virtual object is a root of a namespace.
func (ms *MetaStore) CheckRoot(nsid Namespace_ID, void VO_ID,
	found error, not_found error) error {
	return checkHelper(ms.namespaces().FindOne(context.TODO(),
		bson.M{"nsid": nsid, "rootobjects": void}), found, not_found)
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

// Checks if a member has access to a specific object.
// This will initially check if the given file and user exist.
func (ms *MetaStore) CheckVObjectAccess(void VO_ID, member User_ID,
	found error, not_found error) error {
	if err := ms.CheckUser(member, nil, ErrUserDoesNotExist); err != nil {
		return err
	}

	// Confirm the virtual object exists.
	if err := ms.CheckVObject(void, nil, not_found); err != nil {
		// Note, not found is returned here instead of VObject Doesn't Exist.
		// We may not be willing to tell the user whether or not the
		// request object exists or not.
		return err
	}

	// Now perform search...
	curr_root := void
	for curr_root != Nil_VO_ID {
		res := ms.vobjects().FindOne(context.TODO(), bson.M{"void": curr_root})

		// This marks an error which is not the user's fault.
		if res.Err() != nil {
			return ErrInternal
		}

		var vobject VObject
		res.Decode(&vobject)

		// We now will seach the given namespaces to see
		// if the user belongs to any of them.
		namespaces := vobject.Namespaces
		res = ms.namespaces().FindOne(context.TODO(),
			bson.M{"nsid": bson.M{"$in": namespaces}, "users": member})

		// If we found a namespace... success!!!
		if res.Err() == nil {
			return found
		}

		// If nothing was found... continue
		if res.Err() == mongo.ErrNoDocuments {
			curr_root = vobject.ClosestRoot
		} else {
			// Normal error case.
			return res.Err()
		}
	}

	// Nothing was found!
	return not_found
}

// TODO Potentially delete this...
// // Retrieve all namespaces a virtual object belongs to.
// func (ms *MetaStore) AccumulateNamespaces(void VO_ID) ([]Namespace_ID, error) {
// 	// Confirm the virtual object exists.
// 	if err := ms.CheckVObject(void, nil, ErrVObjectNotFound); err != nil {
// 		return nil, err
// 	}

// 	var vobject VObject
// 	var namespaces []Namespace_ID

// 	curr_root := void
// 	for curr_root != Nil_VO_ID {
// 		res := ms.vobjects().FindOne(context.TODO(), bson.M{"void": void})

// 		// There should always be an object tied to the IDs in this loop.
// 		// If this is not the case, there must be some error with the
// 		// data in the DB... not the user's fault!
// 		if res.Err() != nil {
// 			return nil, ErrInternal
// 		}

// 		res.Decode(&vobject)

// 		namespaces = append(namespaces, vobject.Namespaces...)
// 		curr_root = vobject.ClosestRoot
// 	}

// 	return namespaces, nil
// }

func (ms *MetaStore) GenerateNewVOID() (VO_ID, error) {
	var new_void VO_ID
	for {
		new_void = VO_ID(uuid.New())

		if new_void == Nil_VO_ID {
			continue
		}

		res := ms.vobjects().FindOne(context.TODO(), bson.M{"void": new_void})

		if res.Err() == mongo.ErrNoDocuments {
			return new_void, nil
		}

		if res.Err() != nil {
			return Nil_VO_ID, res.Err()
		}
	}
}

var (
	ErrInternal            = errors.New("Inconsistent results!")
	ErrNotImplemented      = errors.New("Not implemented!")
	ErrMachineExists       = errors.New("Machine already exists!")
	ErrMachineDoesNotExist = errors.New("Machine does not exist!")
	ErrUserExists          = errors.New("User already exists!")
	ErrUserDoesNotExist    = errors.New("User does not exist!")
	ErrNoAccess            = errors.New("Unable to access namespace!")
	ErrAccess              = errors.New("Namespace is accessible!")
	ErrNoUserInNamespace   = errors.New("Cannot find user in namespace!")
	ErrVObjectNotFound     = errors.New("Virtual object not found!")
	ErrNotADirectory       = errors.New("Not a directory!")
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
		NSID:        nsid,
		Name:        name,
		RootObjects: []VO_ID{},
		Users:       []User_ID{member},
	}

	// Finally, upload the Namespace to Mongo.
	_, err := ms.namespaces().InsertOne(context.TODO(), namespace)

	if err != nil {
		return Nil_Namespace_ID, err
	}

	return nsid, nil
}

func (ms *MetaStore) DeleteNamespace(nsid Namespace_ID, member User_ID) error {
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
	// I.e. remove this namespace from all existing files!
	// ----------------------------
}

func (ms *MetaStore) AddUserToNamespace(nsid Namespace_ID, recruiter User_ID,
	recruit User_ID) error {
	// First, check if the recruit is a real user.
	if err := ms.CheckUser(recruit, nil, ErrUserDoesNotExist); err != nil {
		return err
	}

	// Next, attempt to update the requested namespace.
	res, err := ms.namespaces().UpdateOne(context.TODO(),
		bson.M{"nsid": nsid, "users": recruiter},
		bson.M{"$push": bson.M{"users": recruit}})

	if err != nil {
		return err
	}

	// This is the situation where no document was updated.
	// I.e. no namespace could be found with the given credentials.
	if res.MatchedCount == 0 {
		return ErrNoAccess
	}

	// Success!
	return nil
}

func (ms *MetaStore) RemoveUserFromNamespace(nsid Namespace_ID,
	axer User_ID, axed User_ID) error {
	// First see if the axer has permission to modify the namespace.
	if err := ms.CheckNamespace(nsid, axer, nil, ErrNoAccess); err != nil {
		return err
	}

	// Next make sure the give user being axed exists in the namespace.
	err := ms.CheckNamespace(nsid, axed, nil, ErrNoUserInNamespace)
	if err != nil {
		return err
	}

	_, err = ms.namespaces().UpdateOne(context.TODO(),
		bson.M{"nsid": nsid, "users": axer},
		bson.M{"$pull": bson.M{"users": axed}})

	return err
}

func (ms *MetaStore) CreateFreeObjectInNamespace(nsid Namespace_ID,
	member User_ID, name string, tp VFM_Object_Type) (VO_ID, error) {
	// First confirm we have access.
	if err := ms.CheckNamespace(nsid, member, nil, ErrNoAccess); err != nil {
		return Nil_VO_ID, err
	}

	// Next generate new object ID.
	void, err := ms.GenerateNewVOID()

	if err != nil {
		return Nil_VO_ID, err
	}

	// Create record for new object.
	vobject := VObject{
		VOID:        void,
		ParentID:    Nil_VO_ID,
		ClosestRoot: Nil_VO_ID,
		Name:        name,
		Namespaces:  []Namespace_ID{nsid},
		IsDir:       tp == VFM_Dir_Type,
		Machines:    []Machine_ID{},
		SubObjects:  []VO_ID{},
	}

	_, err = ms.vobjects().InsertOne(context.TODO(), vobject)

	if err != nil {
		return Nil_VO_ID, err
	}

	// Lastly add root to the Namespace roots slice.
	_, err = ms.namespaces().UpdateOne(context.TODO(),
		bson.M{"nsid": nsid}, bson.M{"$push": bson.M{"rootobjects": void}})

	if err != nil {
		return Nil_VO_ID, err
	}

	return void, nil
}

// TODO --------------------------------
// AddObjectToNamespace(nsid Namespace_ID,
// 	member User_ID, object VO_ID) error

// TODO ---------------------------------
// RemoveObjectFromNamespace(nsid Namespace_ID, member User_ID,
// 	object VO_ID) error

// TODO -------------------------------
// GetNamespaceDetails(nsid Namespace_ID, member User_ID) (*Namespace, error)

func (ms *MetaStore) CreateObject(parent VO_ID, member User_ID, name string,
	tp VFM_Object_Type) (VO_ID, error) {
	// First off, does the user have access to said parent.
	err := ms.CheckVObjectAccess(parent, member, nil, ErrNoAccess)
	if err != nil {
		return Nil_VO_ID, err
	}

	// Generate a new ID for the object and add it to the parent's
	// subobjects array.
	void, err := ms.GenerateNewVOID()
	if err != nil {
		return Nil_VO_ID, err
	}

	res, err := ms.vobjects().UpdateOne(context.TODO(),
		bson.M{"void": parent, "isdir": true},
		bson.M{"$push": bson.M{"subobjects": void}})

	if err != nil {
		return Nil_VO_ID, err
	}

	// The requested object was not a directory!!!
	if res.MatchedCount == 0 {
		return Nil_VO_ID, ErrNotADirectory
	}

	var p_vobject VObject
	fres := ms.vobjects().FindOne(context.TODO(), bson.M{"void": parent})

	if fres.Err() == mongo.ErrNoDocuments {
		return Nil_VO_ID, ErrInternal
	}

	if fres.Err() != nil {
		return Nil_VO_ID, fres.Err()
	}

	// Otherwise, the parent was found!
	fres.Decode(&p_vobject)

	// If the parent object is the root object of any namespaces,
	// then the parent object will be the new object's closest root.
	// Otherwise, we will simply pass the parent's closest root to
	// the child.
	var closest_root VO_ID
	if len(p_vobject.Namespaces) == 0 {
		closest_root = p_vobject.ClosestRoot
	} else {
		closest_root = parent
	}

	new_vobject := VObject{
		VOID:        void,
		ParentID:    parent,
		ClosestRoot: closest_root,
		Name:        name,
		Namespaces:  []Namespace_ID{},
		IsDir:       tp == VFM_Dir_Type,
		Machines:    []Machine_ID{},
		SubObjects:  []VO_ID{},
	}

	_, err = ms.vobjects().InsertOne(context.TODO(), new_vobject)

	if err != nil {
		return Nil_VO_ID, err
	}

	return void, nil
}

func (ms *MetaStore) DeleteObject(void VO_ID, member User_ID) error {
	return nil
}
