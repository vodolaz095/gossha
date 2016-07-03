// +build !nosqlite3

package models

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite" //it is ok
)
