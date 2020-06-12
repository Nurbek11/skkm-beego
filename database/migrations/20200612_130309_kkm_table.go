package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200612_130309 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200612_130309{}
	m.Created = "20200612_130309"

	migration.Register("KkmTable_20200612_130309", m)
}

// Run the migrations
func (m *KkmTable_20200612_130309) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY, title VARCHAR(255),cash VARCHAR(255), organization_bin VARCHAR(255))")

}

// Reverse the migrations
func (m *KkmTable_20200612_130309) Down() {
	m.SQL("DROP TABLE kkm")

}
