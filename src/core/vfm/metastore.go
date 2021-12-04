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

// Is an object a root object of a given namespace?
func (vo *VObject) IsRootOf(nsid Namespace_ID) bool {
	for _, ns := range vo.Namespaces {
		if ns == nsid {
			return true
		}
	}

	return false
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

func (ms *MetaStore) Destruct() error {
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

// Confirm a namespace with the given nsid and user exists.
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

// Decode an object which is expected to exist in the database.
func (ms *MetaStore) DecodeExistingVObject(void VO_ID, vobject *VObject) error {
	res := ms.vobjects().FindOne(context.TODO(), bson.M{"void": void})

	if res.Err() == mongo.ErrNoDocuments {
		return ErrInternal
	}

	if res.Err() != nil {
		return res.Err()
	}

	err := res.Decode(vobject)

	return err
}

// Checks if a member has access to a specific object.
// This function assumes the given object is real.
func (ms *MetaStore) CheckVObjectAccess(void VO_ID, member User_ID,
	found error, not_found error) error {
	// Now perform search...
	curr_root := void
	for curr_root != Nil_VO_ID {
		var vobject VObject
		err := ms.DecodeExistingVObject(curr_root, &vobject)
		if err != nil {
			return err
		}

		// We now will seach the given namespaces to see
		// if the user belongs to any of them.
		namespaces := vobject.Namespaces
		res := ms.namespaces().FindOne(context.TODO(),
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

// Check to see if an object belongs to a Namespace.
// This assumes the given object exists.
func (ms *MetaStore) CheckVObjectNamespace(void VO_ID, nsid Namespace_ID,
	found error, not_found error) error {
	curr_root := void
	for curr_root != Nil_VO_ID {
		var vobject VObject
		err := ms.DecodeExistingVObject(curr_root, &vobject)
		if err != nil {
			return err
		}

		if vobject.IsRootOf(nsid) {
			return found
		}

		curr_root = vobject.ClosestRoot
	}

	return not_found
}

func (ms *MetaStore) AccumulateAccessibleNamespaces(void VO_ID,
	member User_ID) ([]Namespace_ID, error) {
	return nil, nil
}

// Check to see if an object exists, and if a user has access to it.
func (ms *MetaStore) CheckVObjectAccessAndExists(void VO_ID, member User_ID,
	found error, not_found error) error {
	if err := ms.CheckVObject(void, nil, not_found); err != nil {
		return err
	}

	return ms.CheckVObjectAccess(void, member, found, not_found)
}

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

func expectSingleUpdate(ures *mongo.UpdateResult, err error) error {
	if err != nil {
		return err
	}

	if ures.MatchedCount == 0 {
		return ErrInternal
	}

	return nil
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
	ErrVObjectFound        = errors.New("Virtual object found!")
	ErrNotADirectory       = errors.New("Not a directory!")
	ErrObjectInNamespace   = errors.New("Already in namespace!")
	ErrNotRoot             = errors.New("Virtual object not a root!")
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

	_, err = ms.vobjects().UpdateMany(context.TODO(),
		bson.M{},
		bson.M{"$pull": bson.M{"machines": old_machine}})

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

func (ms *MetaStore) GetFileLocations(file_id VO_ID) ([]Machine_ID, error) {
	res := ms.vobjects().FindOne(context.TODO(),
		bson.M{"void": file_id})

	if res.Err() == mongo.ErrNoDocuments {
		return nil, ErrVObjectNotFound
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	var vobject VObject
	err := res.Decode(&vobject)
	if err != nil {
		return nil, err
	}

	return vobject.Machines, nil
}

func (ms *MetaStore) SetFileLocations(file_id VO_ID, locs []Machine_ID) error {
	// First confirm all given locations are valid.
	fres, err := ms.machines().Find(context.TODO(),
		bson.M{"mid": bson.M{"$in": locs}})

	if err != nil {
		return err
	}

	if fres.RemainingBatchLength() != len(locs) {
		return ErrMachineDoesNotExist
	}

	ures, err := ms.vobjects().UpdateOne(context.TODO(),
		bson.M{"void": file_id, "isdir": false},
		bson.M{"$set": bson.M{"machines": locs}})

	if err != nil {
		return err
	}

	if ures.MatchedCount == 0 {
		return ErrVObjectNotFound
	}

	return nil
}

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
	res := ms.namespaces().FindOne(context.TODO(), filter)

	if res.Err() == mongo.ErrNoDocuments {
		return ErrNoAccess
	}

	if res.Err() != nil {
		return res.Err()
	}

	var namespace Namespace
	err := res.Decode(&namespace)
	if err != nil {
		return err
	}

	// First remove all Root objects from the namespace.
	for _, rvoid := range namespace.RootObjects {
		err := ms.DirectRemoveObjectFromNamespace(nsid, rvoid, ErrInternal)
		if err != nil {
			return err
		}
	}

	// Finally, delete the namespace itself.
	dres, err := ms.namespaces().DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}

	if dres.DeletedCount == 0 {
		return ErrInternal
	}

	return nil
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

func (ms *MetaStore) AddObjectToNamespace(nsid Namespace_ID,
	member User_ID, void VO_ID) error {
	// A few things need to be checked here...
	// Does the given file exist, and does the user have access to it???
	err := ms.CheckVObjectAccessAndExists(void, member, nil, ErrNoAccess)
	if err != nil {
		return err
	}

	// Next check if the given namespace exists, and does the user have access
	// to it???
	err = ms.CheckNamespace(nsid, member, nil, ErrNoAccess)
	if err != nil {
		return err
	}

	// Finally, check if the given object already belongs to the
	// given namespace.
	err = ms.CheckVObjectNamespace(void, nsid, ErrObjectInNamespace, nil)
	if err != nil {
		return err
	}

	// Addition process...

	// First, tag the given object.
	err = ms.TagObject(nsid, void, true)
	if err != nil {
		return err
	}

	var vobject VObject
	err = ms.DecodeExistingVObject(void, &vobject)

	if err != nil {
		return err
	}

	// Lastly, reroute subobjects if needed.
	if vobject.IsDir {
		for _, svoid := range vobject.SubObjects {
			err = ms.RerouteVObject(nsid, void, svoid, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Recursive helper for routing subobjects to a new closest root.
func (ms *MetaStore) RerouteVObject(nsid Namespace_ID, root VO_ID,
	void VO_ID, root_found bool) error {
	// We do not reroute the starting object.
	err := expectSingleUpdate(ms.vobjects().UpdateOne(context.TODO(),
		bson.M{"void": void},
		bson.M{"$set": bson.M{"closestroot": root}}))

	if err != nil {
		return err
	}

	var vobject VObject
	err = ms.DecodeExistingVObject(void, &vobject)

	if err != nil {
		return err
	}

	num_nsids := len(vobject.Namespaces)
	is_nsid_root := false

	// If a root of nsid is yet to be found and
	// this object is a root of nsid, we must untag
	// it.
	// If a root has already been found, there is no
	// reason to check this subobject for being a root.
	if !root_found && vobject.IsRootOf(nsid) {
		is_nsid_root = true

		err = ms.TagObject(nsid, void, false)
		if err != nil {
			return err
		}

		num_nsids--
	}

	// If our object now is a directory and not a root.
	// we must reroute all of its subobjects.
	if vobject.IsDir && num_nsids == 0 {
		for _, svoid := range vobject.SubObjects {
			// When
			err = ms.RerouteVObject(nsid, root, svoid,
				root_found || is_nsid_root)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Either adds or removes a namespace tag for an object
func (ms *MetaStore) TagObject(nsid Namespace_ID, void VO_ID, add bool) error {
	var operation string
	if add {
		operation = "$push"
	} else {
		operation = "$pull"
	}

	err := expectSingleUpdate(ms.vobjects().UpdateOne(context.TODO(),
		bson.M{"void": void},
		bson.M{operation: bson.M{"namespaces": nsid}}))

	if err != nil {
		return err
	}

	return expectSingleUpdate(ms.namespaces().UpdateOne(context.TODO(),
		bson.M{"nsid": nsid},
		bson.M{operation: bson.M{"rootobjects": void}}))
}

func (ms *MetaStore) RemoveObjectFromNamespace(nsid Namespace_ID,
	member User_ID, void VO_ID) error {
	// First check if this object exists and if the user has access to it.
	err := ms.CheckVObjectAccessAndExists(void, member, nil, ErrNoAccess)
	if err != nil {
		return err
	}

	// Next check if the user has access to the given nsid.
	err = ms.CheckNamespace(nsid, member, nil, ErrNoAccess)
	if err != nil {
		return err
	}

	return ms.DirectRemoveObjectFromNamespace(nsid, void, ErrNotRoot)
}

// Delete an object from a namespace without checking for user.
// If the object is not a root of said namespace, not_found is returned.
func (ms *MetaStore) DirectRemoveObjectFromNamespace(nsid Namespace_ID,
	void VO_ID, not_found error) error {
	var vobject VObject
	if err := ms.DecodeExistingVObject(void, &vobject); err != nil {
		return err
	}

	if !vobject.IsRootOf(nsid) {
		return not_found
	}

	// At this point, we know we are working with a root object of a
	// namespace the user has access to.

	// Untag the object.
	if err := ms.TagObject(nsid, void, false); err != nil {
		return err
	}

	// If nsid was the only nsid vobject was the root of, we must
	// make sure to reroute all subobjects!
	if vobject.IsDir && len(vobject.Namespaces) == 1 {
		for _, svoid := range vobject.SubObjects {
			// We don't want to do any namespace manipulation of these
			// subobjects, so we use the Nil_Namespace_ID here...
			err := ms.RerouteVObject(Nil_Namespace_ID, vobject.ClosestRoot,
				svoid, true)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ms *MetaStore) GetNamespaceDetails(nsid Namespace_ID,
	member User_ID) (*Namespace, error) {
	// First make sure the user has access to the given namespace.
	if err := ms.CheckNamespace(nsid, member, nil, ErrNoAccess); err != nil {
		return nil, err
	}

	res := ms.namespaces().FindOne(context.TODO(),
		bson.M{"nsid": nsid})

	if res.Err() == mongo.ErrNoDocuments {
		return nil, ErrInternal
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	var namespace Namespace
	err := res.Decode(&namespace)
	if err != nil {
		return nil, err
	}

	return &namespace, nil
}

func (ms *MetaStore) CreateObject(parent VO_ID, member User_ID, name string,
	tp VFM_Object_Type) (VO_ID, error) {
	// First off, does the user have access to said parent.
	err := ms.CheckVObjectAccessAndExists(parent, member, nil, ErrNoAccess)
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
	err = ms.DecodeExistingVObject(parent, &p_vobject)
	if err != nil {
		return Nil_VO_ID, err
	}

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
	// Delete an object...
	// First we need to make sure the user has access to this object.
	// Some sort of deletion helper might be helpful here.
	err := ms.CheckVObjectAccessAndExists(void, member, nil, ErrNoAccess)
	if err != nil {
		return err
	}

	// Delegate to helper after access has been confirmed.
	return ms.DirectDeleteObject(void)
}

// Delete an object with no user or existence checking.
// Built as a recursive helper.
func (ms *MetaStore) DirectDeleteObject(void VO_ID) error {
	var vobject VObject
	err := ms.DecodeExistingVObject(void, &vobject)
	if err != nil {
		return err
	}

	ures, err := ms.namespaces().UpdateMany(context.TODO(),
		bson.M{"nsid": bson.M{"$in": vobject.Namespaces}},
		bson.M{"$pull": bson.M{"rootobjects": void}})

	if err != nil {
		return err
	}

	// There must be something wrong!!
	// This object was marked as a root, but is not
	// in the roots slice for its given namespace.
	if ures.MatchedCount != int64(len(vobject.Namespaces)) {
		return ErrInternal
	}

	// Finally do deletion.
	dres, err := ms.vobjects().DeleteOne(context.TODO(), bson.M{"void": void})

	if err != nil {
		return err
	}

	if dres.DeletedCount == 0 {
		return ErrInternal
	}

	// Recursively delete all subobjects.
	if vobject.IsDir {
		for _, svoid := range vobject.SubObjects {
			err = ms.DirectDeleteObject(svoid)
			if err != nil {
				return err
			}
		}
	}

	// Success!
	return nil
}

func (ms *MetaStore) GetObjectDetails(void VO_ID,
	member User_ID) (VFM_Object, error) {
	// First off, does user have access to this object.
	err := ms.CheckVObjectAccessAndExists(void, member, nil, ErrNoAccess)
	if err != nil {
		return nil, err
	}

	// Next let's get the object itself.
	var vobject VObject
	err = ms.DecodeExistingVObject(void, &vobject)
	if err != nil {
		return nil, err
	}

	// Next, let's get all accessible namespaces from void alone.
	ns_stubs, err := ms.GetAccessibleNamespaces(void, member)
	if err != nil {
		return nil, err
	}

	header := VFM_Header{
		id:         void,
		parent_id:  vobject.ParentID,
		name:       vobject.Name,
		namespaces: ns_stubs,
	}

	// If we aren't working with a directory...
	// stop here.
	if !vobject.IsDir {
		return VFM_File{header}, nil
	}

	// Time to find sub object stubs.
	var subs []VFM_Object_Stub
	for _, svoid := range vobject.SubObjects {
		var sobject VObject
		err := ms.DecodeExistingVObject(svoid, &sobject)
		if err != nil {
			return nil, err
		}

		access_roots, err := ms.GetAccessibleRootNamespaces(&sobject, member)
		if err != nil {
			return nil, err
		}

		var ty VFM_Object_Type
		if sobject.IsDir {
			ty = VFM_Dir_Type
		} else {
			ty = VFM_File_Type
		}

		subs = append(subs, VFM_Object_Stub{
			Id:         sobject.VOID,
			Name:       sobject.Name,
			Namespaces: access_roots,
			Type:       ty,
		})
	}

	return VFM_Directory{
		header,
		subs,
	}, nil
}

// Accumulate all namespaces void belongs to which member has access
// to.
func (ms *MetaStore) GetAccessibleNamespaces(void VO_ID,
	member User_ID) ([]Namespace_Stub, error) {
	var accessible []Namespace_Stub
	curr_root := void
	for curr_root != Nil_VO_ID {
		var vobject VObject
		err := ms.DecodeExistingVObject(curr_root, &vobject)
		if err != nil {
			return nil, err
		}

		access_roots, err := ms.GetAccessibleRootNamespaces(&vobject, member)
		if err != nil {
			return nil, err
		}

		accessible = append(accessible, access_roots...)
		curr_root = vobject.ClosestRoot
	}

	return accessible, nil
}

// Given an object and a user. Retrieve all Namespaces for which
// the object is a root of and the user has access to.
func (ms *MetaStore) GetAccessibleRootNamespaces(vobject *VObject,
	member User_ID) ([]Namespace_Stub, error) {
	cursor, err := ms.namespaces().Find(context.TODO(),
		bson.M{"nsid": bson.M{"$in": vobject.Namespaces}, "users": member})
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	var found []Namespace
	if err = cursor.All(context.TODO(), &found); err != nil {
		return nil, err
	}

	var access_roots []Namespace_Stub
	for _, ns := range found {
		access_roots = append(access_roots,
			Namespace_Stub{NSID: ns.NSID, Name: ns.Name},
		)
	}

	return access_roots, nil
}
