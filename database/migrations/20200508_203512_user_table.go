package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserTable_20200508_203512 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTable_20200508_203512{}
	m.Created = "20200508_203512"

	migration.Register("UserTable_20200508_203512", m)
}

// Run the migrations
func (m *UserTable_20200508_203512) Up() {
	m.SQL("CREATE TABLE users(id serial PRIMARY KEY, username VARCHAR(255), password VARCHAR(255),role VARCHAR(255))")


}

// Reverse the migrations
func (m *UserTable_20200508_203512) Down() {
	m.SQL("DROP TABLE users")

}
