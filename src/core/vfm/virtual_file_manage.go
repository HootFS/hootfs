package vfm

import "fmt"

// TODO, change this to whichever UUID implementaton we use.
type VFM_UUID struct{}

// TODO, change this to whichever user ID implementation we use.
type User_ID struct{}

// TODO, change this to whichever Machine ID implementation we use.
type Machine_ID struct{}

// VFM_Obj_Type will classify objects as either files or directorys.
type VFM_Object_Type int

const (
	VFM_File_Type VFM_Object_Type = iota // Placeholder type.
	VFM_Dir_Type
)

// A VFM_Object is the Go representation of a file or folder's
// metadata.
type VFM_Object interface {
	GetID() VFM_UUID
	GetParentID() (VFM_UUID, error)
	GetName() string

	// NOTE, the details of an object will usually be requested
	// by a user. This call should not return the IDs of every namespace
	// this Object belongs to. SOMETHING WE NEED TO THINK ABOUT HERE!!!!!
	GetNamespaces() []VFM_UUID

	GetObjectType() VFM_Object_Type
	GetSubObjects() ([]VFM_UUID, error)
}

type VFM_Header struct {
	id         VFM_UUID
	parent_id  VFM_UUID
	name       string
	namespaces []VFM_UUID
}

func (h VFM_Header) GetID() VFM_UUID {
	return h.id
}

func (h VFM_Header) GetParentID() VFM_UUID {
	return h.parent_id
}

func (h VFM_Header) GetName() string {
	return h.name
}

func (h VFM_Header) GetNamespaces() []VFM_UUID {
	return h.namespaces
}

type VFM_File struct {
	VFM_Header
}

func (f VFM_File) GetObjectType() VFM_Object_Type {
	return VFM_File_Type
}

func (f VFM_File) GetSubObjects() ([]VFM_UUID, error) {
	return nil, fmt.Errorf("Files do not have sub objects! (%s)", f.name)
}

type VFM_Directory struct {
	VFM_Header
	sub_objects []VFM_UUID
}

func (f VFM_Directory) GetObjectType() VFM_Object_Type {
	return VFM_Dir_Type
}

func (f VFM_Directory) GetSubObjects() ([]VFM_UUID, error) {
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
	GetFileLocations(file_id VFM_UUID) ([]Machine_ID, error)

	// Set the machines a file resides on.
	//		file_id		- The ID of the file in question.
	//		locs		- The new slice of locations.
	SetFileLocations(file_id VFM_UUID, locs []Machine_ID) error

	// User Focused Functions ----------------------------------------------

	// This simply creates a new Namespace.
	// 		name    - the name of the Namespace.
	//		member - the ID of the user creating the Namespace.
	CreateNamespace(name string, member User_ID) (VFM_UUID, error)

	// Delete a Namespace.
	// NOTE, this functionality brings up garabe collection relating
	// issues. For example, when a file no longer belongs to any
	// Namespaces, should it be deleted from the cluster?
	//		ns_id	- The Namespace in question.
	//		member	- The user deleting the Namespace.
	DeleteNamespace(ns_id VFM_UUID, member User_ID) error

	// Adds a member to a Namespace. One a user is the member of a Namespace,
	// he or she has access to all objects inside the Namespace. He or she
	// also has the ability to add other users to the Namespace.
	//		ns_id		- The Namespace in question.
	//		recruiter 	- The user adding the recruit. This user must be
	// 					  a member of the Namespace or else an error will
	//                    be returned.
	//		recruit		- The user to add to the Namespace.
	AddUserToNamespace(ns_id VFM_UUID, recruiter User_ID, recruit User_ID) error

	// Removes a user from a Namespace.
	//		ns_id	- The Namespace in question.
	//		axer	- The user performing the remove.
	//		axed	- The user being removed.
	RemoveUserFromNamespace(ns_id VFM_UUID, axer User_ID, axed User_ID) error

	// NOTE, a "Root Object" of a namespace N is a directory or file which
	// either (a) has no parent directory, or (b) has a parent directory which
	// does not belong to namespace N.

	// Adds an object to a Namespace. (Either a folder or file)
	// This file cannot already belong to the given namespace.
	// If added successfully, the added object will be a Root object of
	// the namespace.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		object	- The object to add to the namespace.
	//				  This object cannot already belong to the Namespace.
	AddObjectToNamespace(ns_id VFM_UUID, member User_ID, object VFM_UUID) error

	// Remove an object from a Namespace. Removing a folder from a Namespace
	// will remove all of its contents from the Namespace as well.
	// If this object is not a "root object" of the Namespace it is
	// being removed from, an error will be returned.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		object	- The object to remove from the namespace.
	RemoveObjectFromNamespace(ns_id VFM_UUID, member User_ID,
		object VFM_UUID) error

	// Create a freestanding Directory in a namespace.
	// I.e. a directory which belongs to no parent directory.
	// By definition, this will be a root object of the namespace.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		name	- The name of the new directory.
	CreateFreeDirInNamespace(ns_id VFM_UUID, member User_ID,
		name string) (VFM_UUID, error)

	// Create a freestanding file in a namespace.
	// This file will belong to no parent directory.
	// Again, this will be a root object of the namespace.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		name	- The name of the new file.
	CreateFreeFileInNamespace(ns_id VFM_UUID, member User_ID,
		name string) (VFM_UUID, error)

	GetNamespaceRoots(ns_id VFM_UUID, member User_ID) []VFM_UUID

	// NOTE, for the next two functions...
	// Creating a directory or file will add said object to the
	// namespace(s) its parent directory.

	// Create a new directory.
	//		parent	- The ID of the parent folder.
	//		member	- The user making the request.
	//		name	- The name of the directory.
	CreateDir(parent VFM_UUID, member User_ID, name string) (VFM_UUID, error)

	// Create a new file.
	//		parent	- The ID of the parent folder.
	//		member	- The user making the request.
	//		name	- The name of the file.
	CreateFile(parent VFM_UUID, member, User_ID, name string) (VFM_UUID, error)

	// Delete an object. If a user is in one namespace which contains the given
	// object, he or she has the ability to delete said object.
	// Deleting an object will remove from disk... thus all of its
	// corresponding namespaces.
	//		obj_id	- The ID of the object to remove (file or folder)
	//		member	- The user making the request.
	DeleteObject(obj_id VFM_UUID, member User_ID) error

	// Get the details of a specific object in the file system.
	//		obj_id	- The ID of the object in question.
	//		member 	- The user making the request.
	GetObjectDetails(obj_id VFM_UUID, member User_ID) (VFM_Object, error)

	// Get the contents of a directory.
	//		dir_id	- The ID of the directory in question.
	//		member	- The user making the request.
	GetDirectoryContents(dir_id VFM_UUID, member User_ID) ([]VFM_Object, error)
}
