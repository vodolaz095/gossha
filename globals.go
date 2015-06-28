package gossha

import (
	"github.com/jinzhu/gorm"
)

// DB is a singleton for database connection, it is initialized by
// InitDatabase function
var DB *gorm.DB

// Board is a map of Handler's, with SessionId's as keys -
// http://godoc.org/golang.org/x/crypto/ssh#ConnMetadata
var Board map[string]*Handler

// RuntimeConfig is a current configuration being used
var RuntimeConfig *Config
