package types

// Permission contains a field - `Allowed` - that tells if the permission is
// allowed and an array - `Reasons` - with codes with the reason(s) why this
// permission is/is not allowed.
//
// As of now, reasons is non-empty only if allowed is false.
type Permission struct {
	Allowed bool     `json:"allowed" yaml:"allowed"`
	Reasons []string `json:"reasons" yaml:"reasons"`
}

// TODO: this will probabily be called in another way.
type UserSettingsPermissions struct {
	CanModifyOwnProfile Permission `json:"can_modify_own_profile" yaml:"canModifyOwnProfile"`
	CanChangeUsername   Permission `json:"can_change_username" yaml:"canChangeUsername"`
	CanChangeDOB        Permission `json:"can_change_dob" yaml:"canChangeDOB"`
}
