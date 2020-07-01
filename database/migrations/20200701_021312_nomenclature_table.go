package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type NomenclatureTable_20200701_021312 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &NomenclatureTable_20200701_021312{}
	m.Created = "20200701_021312"

	migration.Register("NomenclatureTable_20200701_021312", m)
}

// Run the migrations
func (m *NomenclatureTable_20200701_021312) Up() {
	m.SQL("CREATE TABLE nomenclature(id serial PRIMARY KEY,organization_bin VARCHAR(255) NOT NULL,title VARCHAR(255),price VARCHAR(255),discount VARCHAR(255), extra_charge VARCHAR(255),number VARCHAR(255),sum VARCHAR(255), is_dis_price boolean, is_dis_discount boolean, is_dis_ex_charge boolean, is_dis_number boolean)")

}

// Reverse the migrations
func (m *NomenclatureTable_20200701_021312) Down() {

}
