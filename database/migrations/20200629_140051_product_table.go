package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ProductTable_20200629_140051 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ProductTable_20200629_140051{}
	m.Created = "20200629_140051"

	migration.Register("ProductTable_20200629_140051", m)
}

// Run the migrations
func (m *ProductTable_20200629_140051) Up() {
	m.SQL("CREATE TABLE product(good_id serial PRIMARY KEY,cheque_id int NOT NULL,good_title VARCHAR(255),good_price VARCHAR(255), good_discount VARCHAR(255), good_extra_charge VARCHAR(255), good_number VARCHAR(255), good_sum VARCHAR(255), is_dis_price boolean, is_dis_discount boolean, is_dis_ex_charge boolean, is_dis_number boolean)")
}

// Reverse the migrations
func (m *ProductTable_20200629_140051) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
