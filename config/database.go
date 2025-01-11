package config

import(
	"fmt",
	"log",
	"os",

	"gorm.io/driver/postgres",
	"gorm.io/gorm"
)

var DB *gorm.db