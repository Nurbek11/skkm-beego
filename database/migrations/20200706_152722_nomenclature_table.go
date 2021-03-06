package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type NomenclatureTable_20200706_152722 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &NomenclatureTable_20200706_152722{}
	m.Created = "20200706_152722"

	migration.Register("NomenclatureTable_20200706_152722", m)
}

// Run the migrations
func (m *NomenclatureTable_20200706_152722) Up() {
	m.SQL("CREATE TABLE nomenclature(id serial PRIMARY KEY,organization_bin VARCHAR(255) NOT NULL,title VARCHAR(255),price VARCHAR(255),discount VARCHAR(255), extra_charge VARCHAR(255),quantity_in_stock VARCHAR(255),sum VARCHAR(255), is_dis_price boolean, is_dis_discount boolean, is_dis_ex_charge boolean, is_dis_number boolean)")

}

// Reverse the migrations
func (m *NomenclatureTable_20200706_152722) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
