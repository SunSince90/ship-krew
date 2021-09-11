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

// GetUserOptions are options that will be used when getting a single user
// from the board.
type GetUserOptions struct {
	// Name of the user.
	Name string
}
