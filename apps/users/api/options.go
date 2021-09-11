package api

// ListUsersOptions are options that will be used when listing users.
type ListUsersOptions struct {
	// Page to load.
	// Default: 0.
	Page int
	// PerPage defines how many users to load per page.
	// Default: 20.
	PerPage int
}
