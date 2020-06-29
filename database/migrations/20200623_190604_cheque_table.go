package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ChequeTable_20200623_190604 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ChequeTable_20200623_190604{}
	m.Created = "20200623_190604"

	migration.Register("ChequeTable_20200623_190604", m)
}

// Run the migrations
func (m *ChequeTable_20200623_190604) Up() {
	m.SQL("CREATE TABLE cheque(id serial PRIMARY KEY,kkm_id int,shift_id int, total_sum VARCHAR(255),change_money VARCHAR(255), total_discount VARCHAR(255), total_charge VARCHAR(255), n_d_s VARCHAR(255),payment_type VARCHAR(255),operation_type VARCHAR(255) )")

}

// Reverse the migrations
func (m *ChequeTable_20200623_190604) Down() {
	m.SQL("DROP TABLE cheque")

}