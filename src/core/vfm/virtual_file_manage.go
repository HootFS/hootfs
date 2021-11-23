package vfm

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

type VFM_Object interface {
	GetID() VFM_UUID
	GetParentID() VFM_UUID
	GetName() string
	GetObjectType() VFM_Object_Type
	GetContents() ([]VFM_UUID, error)
}

type VFM_Header struct {
	id        VFM_UUID
	parent_id VFM_UUID
	name      string
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

type VFM_File struct {
	VFM_Header
}

type VFM_Directory struct {
	VFM_Header
	contents
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

	// Adds an object to a Namespace. (Either a folder or file)
	// When a folder is in a Namespace, all of its subcontents are in the
	// Namespace as well.
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

	// NOTE, a "Root Object" of a Namespace is a directory or file of a
	// Namespace which has no parent directory.
	// The creation of Root Objects is required to create the initial files
	// of the file system.

	// Create a root directory in a Namespace.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		name	- The name of the new directory.
	CreateRootDirInNamespace(ns_id VFM_UUID, member User_ID,
		name string) (VFM_UUID, error)

	// Create a root file in a Namespace.
	//		ns_id	- The Namespace to add to.
	//		member	- The user making the request.
	//		name	- The name of the new file.
	CreateRootFileInNamespace(ns_id VFM_UUID, member User_ID,
		name string) (VFM_UUID, error)

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
}
