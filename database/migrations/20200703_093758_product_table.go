package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ProductTable_20200703_093758 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ProductTable_20200703_093758{}
	m.Created = "20200703_093758"

	migration.Register("ProductTable_20200703_093758", m)
}

// Run the migrations
func (m *ProductTable_20200703_093758) Up() {
	m.SQL("CREATE TABLE product(id serial PRIMARY KEY,cheque_id int NOT NULL,title VARCHAR(255),section_code VARCHAR(255), price VARCHAR(255),discount_sum int,discount VARCHAR(255),discount_storno bool, extra_charge VARCHAR(255),number VARCHAR(255),sum VARCHAR(255),markup_sum int,markup_storno bool,total int, is_dis_price boolean, is_dis_discount boolean, is_dis_ex_charge boolean, is_dis_number boolean,is_storno bool)")

}

// Reverse the migrations
func (m *ProductTable_20200703_093758) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
