package vfm

import (
	"fmt"

	"github.com/google/uuid"
)

// This will be the UUID used for all VFM Objects.
type VO_ID uuid.UUID

// This represents the Nil VFM UUID. (All 0s)
var Nil_VO_ID VO_ID

// This will be the UUID used for Namespaces.
type Namespace_ID uuid.UUID

var Nil_Namespace_ID Namespace_ID

// User ID will probably be a username/email. (Consult Auth0 for this)
type User_ID string

// The machine IDs will not be created by the VFM, thus an integer is
// probably the best option.
type Machine_ID uint64

// Small portion of namespace information.
type Namespace_Stub struct {
	NSID Namespace_ID
	Name string
}

// All namespace information accessible by a user.
type Namespace struct {
	NSID        Namespace_ID
	Name        string
	RootObjects []VO_ID

	// Users which have access to this Namespace.
	Users []User_ID
}

// VFM_Obj_Type will classify objects as either files or directorys.
type VFM_Object_Type int

const (
	VFM_File_Type VFM_Object_Type = iota // Placeholder type.
	VFM_Dir_Type
)

// This struct is a very condensed way of representing
// a VFM Object. It will be used when the contents of a
// directory are requested by a user.
type VFM_Object_Stub struct {
	Id   VO_ID
	Name string

	// NOTE, This stub does not contain a parentID
	// field. We are assuming the parentID of this object
	// will already be known by the user.
	// Similarly, we are assuming the Namespaces of the parent
	// will also be known by the user. This field should only
	// list Namespaces which are not known to the user.
	// I.e. Namespaces this object belongs to, but its parent
	// does not.
	Namespaces []Namespace_Stub

	Type VFM_Object_Type
}

// A VFM_Object is the Go representation of a file or folder's
// metadata. Forseeably, this will be something which is returned to
// a user. Thus, the information held will be specific to the user
// requesting the information.
// Most notably, Object and Namespace IDs given in this object
// will only be those which the requesting user has access to.
type VFM_Object interface {
	GetID() VO_ID

	// This Object may not have a parent.
	// or... Its parent may not be accesibly by the
	// requesting user. In these cases, an error is
	// returned here.
	GetParentID() (VO_ID, error)

	GetName() string

	// NOTE, the details of an object will usually be requested
	// by a user. This call should not return the IDs of every namespace
	// this Object belongs to.
	GetNamespaces() []Namespace_Stub

	GetObjectType() VFM_Object_Type
	GetSubObjects() ([]VFM_Object_Stub, error)
}

type VFM_Header struct {
	id         VO_ID
	parent_id  VO_ID
	name       string
	namespaces []Namespace_Stub
}

func (h VFM_Header) GetID() VO_ID {
	return h.id
}

func (h VFM_Header) GetParentID() (VO_ID, error) {
	if h.parent_id == Nil_VO_ID {
		return Nil_VO_ID,
			fmt.Errorf(
				"Parent is inaccessbile or does not exist! (%s)\n", h.name)
	}

	return h.parent_id, nil
}

func (h VFM_Header) GetName() string {
	return h.name
}

func (h VFM_Header) GetNamespaces() []Namespace_Stub {
	return h.namespaces
}

type VFM_File struct {
	VFM_Header
}

func (f VFM_File) GetObjectType() VFM_Object_Type {
	return VFM_File_Type
}

func (f VFM_File) GetSubObjects() ([]VFM_Object_Stub, error) {
	return nil, fmt.Errorf("Files do not have sub objects! (%s)\n", f.name)
}

type VFM_Directory struct {
	VFM_Header
	sub_objects []VFM_Object_Stub
}

func (f VFM_Directory) GetObjectType() VFM_Object_Type {
	return VFM_Dir_Type
}

func (f VFM_Directory) GetSubObjects() ([]VFM_Object_Stub, error) {
	return f.sub_objects, nil
}

// The virtual file manager will deal with all meta data about the
// file system structure. I.e. everything up until actually storing
// objects on machines.
// The virtual file manager will combine with the old "SystemMapper".
// I.e. the VFM will hold data about which machines a specific file
// resides in.
type VirtualFileManager interface {

	// Admin Focused Functions -------------------------------------------

	// Add a new machine to the system.
	//		new_machine	- The ID of the machine to add.
	CreateMachine(new_machine Machine_ID) error

	// Delete an existing machine from the system.
	//		old_machine - The ID of the machine to delete.
	DeleteMachine(old_machine Machine_ID) error

	// Add a new user to the system.
	//		new_user	- The ID of the user to add.
	CreateUser(new_user User_ID) error

	// Remove a user from the system.
	//		old_user 	- The ID of the user to remove.
	DeleteUser(old_user User_ID) error

	// Get all the machines a file resides on.
	//		file_id		- The ID of the file in question.
	GetFileLocations(file_id VO_ID) ([]Machine_ID, error)

	// Set the machines a file resides on.
	//		file_id		- The ID of the file in question.
	//		locs		- The new slice of locations.
	SetFileLocations(file_id VO_ID, locs []Machine_ID) error

	// Function for doing all needed cleanup of a VFM.
	Destruct() error

	// User Focused Functions ----------------------------------------------

	GetNamespaces(member User_ID) ([]Namespace_Stub, error)

	// This simply creates a new Namespace.
	// 		name    - the name of the Namespace.
	//		member - the ID of the user creating the Namespace.
	CreateNamespace(name string, member User_ID) (Namespace_ID, error)

	// Delete a Namespace.
	// NOTE, this functionality brings up garbage collection relating
	// issues. For example, when a file no longer belongs to any
	// Namespaces, should it be deleted from the cluster?
	//		nsid	- The Namespace in question.
	//		member	- The user deleting the Namespace.
	DeleteNamespace(nsid Namespace_ID, member User_ID) error

	// Adds a member to a Namespace. One a user is the member of a Namespace,
	// he or she has access to all objects inside the Namespace. He or she
	// also has the ability to add other users to the Namespace.
	//		nsid		- The Namespace in question.
	//		recruiter 	- The user adding the recruit. This user must be
	// 					  a member of the Namespace or else an error will
	//                    be returned.
	//		recruit		- The user to add to the Namespace.
	AddUserToNamespace(nsid Namespace_ID,
		recruiter User_ID, recruit User_ID) error

	// Removes a user from a Namespace.
	//		nsid	- The Namespace in question.
	//		axer	- The user performing the remove.
	//		axed	- The user being removed.
	RemoveUserFromNamespace(nsid Namespace_ID,
		axer User_ID, axed User_ID) error

	// NOTE, a "Root Object" of a namespace N is a directory or file which
	// either (a) has no parent directory, or (b) has a parent directory which
	// does not belong to namespace N.

	// Create a freestanding Object in a namespace.
	// I.e. a directory which belongs to no parent directory.
	// By definition, this will be a root object of the namespace.
	//		nsid	- The Namespace to add to.
	//		member	- The user making the request.
	//		name	- The name of the new object.
	//		tp		- The type of the new object.
	CreateFreeObjectInNamespace(nsid Namespace_ID, member User_ID,
		name string, tp VFM_Object_Type) (VO_ID, error)

	// Adds an object to a Namespace. (Either a folder or file)
	// This file cannot already belong to the given namespace.
	// If added successfully, the added object will be a Root object of
	// the namespace.
	//		nsid	- The Namespace to add to.
	//		member	- The user making the request.
	//		void	- The object to add to the namespace.
	//				  This object cannot already belong to the Namespace.
	AddObjectToNamespace(nsid Namespace_ID,
		member User_ID, void VO_ID) error

	// Remove an object from a Namespace. Removing a folder from a Namespace
	// will remove all of its contents from the Namespace as well.
	// If this object is not a "root object" of the Namespace it is
	// being removed from, an error will be returned.
	//		nsid	- The Namespace to add to.
	//		member	- The user making the request.
	//		void	- The object to remove from the namespace.
	RemoveObjectFromNamespace(nsid Namespace_ID, member User_ID,
		void VO_ID) error

	// Get the IDs of every root object in a given name space.
	// 		nsid	- The Namespace in question.
	//		member	- The user making the request.
	GetNamespaceDetails(nsid Namespace_ID, member User_ID) (*Namespace, error)

	// Create a new object. This object will be added to the namespaces
	// of its parent.
	//		parent	- The ID of the parent folder.
	//		member	- The user making the request.
	//		name	- The name of the object.
	// 		tp		- The type of the object.
	CreateObject(parent VO_ID, member User_ID, name string,
		tp VFM_Object_Type) (VO_ID, error)

	// Delete an object. If a user is in one namespace which contains the given
	// object, he or she has the ability to delete said object.
	// Deleting an object will remove from disk... thus all of its
	// corresponding namespaces.
	//		obj_id	- The ID of the object to remove (file or folder)
	//		member	- The user making the request.
	DeleteObject(void VO_ID, member User_ID) error

	// Get the details of a specific object in the file system.
	//		obj_id	- The ID of the object in question.
	//		member 	- The user making the request.
	GetObjectDetails(void VO_ID, member User_ID) (VFM_Object, error)
}
