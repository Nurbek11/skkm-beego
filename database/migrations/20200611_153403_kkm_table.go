package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200611_153403 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200611_153403{}
	m.Created = "20200611_153403"

	migration.Register("KkmTable_20200611_153403", m)
}

// Run the migrations
func (m *KkmTable_20200611_153403) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY, title VARCHAR(255),cash VARCHAR(255), organization_bin bigint)")

}

// Reverse the migrations
func (m *KkmTable_20200611_153403) Down() {
	m.SQL("DROP TABLE kkm")

}
