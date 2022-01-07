package database

import "gorm.io/gorm"

// TODO: withRole("admin|guest") => return Select()

func byUserName(username string) func(db *gorm.DB) *gorm.DB {
	// TODO: support get deleted ones too
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("username = ? AND deleted_at IS NULL", username)
	}
}

func byUserID(id int64) func(db *gorm.DB) *gorm.DB {
	// TODO: support get deleted ones too
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("id = ? AND deleted_at IS NULL", id)
	}
}

func byEmail(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("email = ? AND deleted_at IS NULL", email)
	}
}
