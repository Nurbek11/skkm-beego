package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200508_203949 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200508_203949{}
	m.Created = "20200508_203949"

	migration.Register("KkmTable_20200508_203949", m)
}

// Run the migrations
func (m *KkmTable_20200508_203949) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY, title VARCHAR(255),cash bigint, organization_id bigint)")
}

// Reverse the migrations
func (m *KkmTable_20200508_203949) Down() {
	m.SQL("drop table kkm")

}
