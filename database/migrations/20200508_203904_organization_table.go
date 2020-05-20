package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type OrganizationTable_20200508_203904 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &OrganizationTable_20200508_203904{}
	m.Created = "20200508_203904"

	migration.Register("OrganizationTable_20200508_203904", m)
}

// Run the migrations
func (m *OrganizationTable_20200508_203904) Up() {
	m.SQL("CREATE TABLE organization(id serial PRIMARY KEY, title VARCHAR(255),bin BIGINT unique,address VARCHAR(255), user_id INT)")

}

// Reverse the migrations
func (m *OrganizationTable_20200508_203904) Down() {
	m.SQL("DROP TABLE organization")

}
