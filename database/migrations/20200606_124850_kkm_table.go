package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200606_124850 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200606_124850{}
	m.Created = "20200606_124850"

	migration.Register("KkmTable_20200606_124850", m)
}

// Run the migrations
func (m *KkmTable_20200606_124850) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY, title VARCHAR(255),cash VARCHAR(255), organization_id bigint)")

}

// Reverse the migrations
func (m *KkmTable_20200606_124850) Down() {
	m.SQL("DROP TABLE kkm")

}
