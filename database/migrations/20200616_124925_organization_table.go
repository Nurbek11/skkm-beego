package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type OrganizationTable_20200616_124925 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &OrganizationTable_20200616_124925{}
	m.Created = "20200616_124925"

	migration.Register("OrganizationTable_20200616_124925", m)
}

// Run the migrations
func (m *OrganizationTable_20200616_124925) Up() {
	m.SQL("CREATE TABLE organization(bin VARCHAR(255) PRIMARY KEY, title VARCHAR(255),address VARCHAR(255), user_id INT)")

}

// Reverse the migrations
func (m *OrganizationTable_20200616_124925) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
