package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200616_124910 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200616_124910{}
	m.Created = "20200616_124910"

	migration.Register("KkmTable_20200616_124910", m)
}

// Run the migrations
func (m *KkmTable_20200616_124910) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY, title VARCHAR(255),cash VARCHAR(255),password VARCHAR(255), organization_bin VARCHAR(255))")

}

// Reverse the migrations
func (m *KkmTable_20200616_124910) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
