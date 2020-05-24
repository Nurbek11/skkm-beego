package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ChequeTable_20200521_173552 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ChequeTable_20200521_173552{}
	m.Created = "20200521_173552"

	migration.Register("ChequeTable_20200521_173552", m)
}

// Run the migrations
func (m *ChequeTable_20200521_173552) Up() {
	m.SQL("CREATE TABLE cheque(id serial PRIMARY KEY, total_sum VARCHAR(255),change_money VARCHAR(255), total_discount VARCHAR(255), total_charge VARCHAR(255), n_d_s VARCHAR(255),payment_type VARCHAR(255))")

}

// Reverse the migrations
func (m *ChequeTable_20200521_173552) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
