package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ProductTable_20200523_113610 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ProductTable_20200523_113610{}
	m.Created = "20200523_113610"

	migration.Register("ProductTable_20200523_113610", m)
}

// Run the migrations
func (m *ProductTable_20200523_113610) Up() {
	m.SQL("CREATE TABLE product(id serial PRIMARY KEY,cheque_id int NOT NULL,title VARCHAR(255),price VARCHAR(255), discount VARCHAR(255), extra_charge VARCHAR(255), number VARCHAR(255), sum VARCHAR(255), is_dis_price boolean, is_dis_discount boolean, is_dis_ex_charge boolean, is_dis_number boolean)")

}

// Reverse the migrations
func (m *ProductTable_20200523_113610) Down() {
	m.SQL("DROP TABLE product")

}
