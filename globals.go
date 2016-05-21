package gossha

import (
	"github.com/jinzhu/gorm"
)

// DB is a singleton for database connection, it is initialized by
// InitDatabase function
var DB *gorm.DB

// RuntimeConfig is a current configuration being used
var RuntimeConfig *Config
