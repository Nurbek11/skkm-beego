package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserTable_20200529_154131 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTable_20200529_154131{}
	m.Created = "20200529_154131"

	migration.Register("UserTable_20200529_154131", m)
}

// Run the migrations
func (m *UserTable_20200529_154131) Up() {
	m.SQL("CREATE TABLE users(id serial PRIMARY KEY, first_name VARCHAR(255),last_name VARCHAR(255),email VARCHAR(255),password VARCHAR(255),role VARCHAR(255))")

}

// Reverse the migrations
func (m *UserTable_20200529_154131) Down() {
	m.SQL("DROP TABLE users")
}
