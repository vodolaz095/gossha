package models

import (
	"time"
)

// Message represents message in chat, it is persisted in relational
// database via https://github.com/jinzhu/gorm object relational mapper
type Message struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	Message   string `sql:"size:255"`
	IP        string `sql:"size:65"`
	Hostname  string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Key stories the hash of public key, used by User to authorize, it is
// persisted in relational database via https://github.com/jinzhu/gorm
// object relational mapper
type Key struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	Content   string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session stories Ip and Hostname, used by User to authorize, it is
// persisted in relational database via https://github.com/jinzhu/gorm
// object relational mapper
type Session struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	IP        string `sql:"size:65"`
	Hostname  string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Notification represents message or system notification emitted by one of Handler's,
// this struct is not persisted in database
type Notification struct {
	User             User
	Message          Message
	IsSystem         bool
	IsChat           bool
	IsPrivateMessage bool
}
