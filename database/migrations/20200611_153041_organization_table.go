package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type OrganizationTable_20200611_153041 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &OrganizationTable_20200611_153041{}
	m.Created = "20200611_153041"
	migration.Register("OrganizationTable_20200611_153041", m)
}

// Run the migrations
func (m *OrganizationTable_20200611_153041) Up() {
	m.SQL("CREATE TABLE organization(bin serial PRIMARY KEY, title VARCHAR(255),address VARCHAR(255), user_id INT)")
}

// Reverse the migrations
func (m *OrganizationTable_20200611_153041) Down() {
	m.SQL("DROP TABLE organization")

}
