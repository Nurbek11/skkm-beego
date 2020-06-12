package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type OrganizationTable_20200612_122617 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &OrganizationTable_20200612_122617{}
	m.Created = "20200612_122617"

	migration.Register("OrganizationTable_20200612_122617", m)
}

// Run the migrations
func (m *OrganizationTable_20200612_122617) Up() {
	m.SQL("CREATE TABLE organization(bin VARCHAR(255) PRIMARY KEY, title VARCHAR(255),address VARCHAR(255), user_id INT)")

}

// Reverse the migrations
func (m *OrganizationTable_20200612_122617) Down() {
	m.SQL("DROP TABLE organization")

}
