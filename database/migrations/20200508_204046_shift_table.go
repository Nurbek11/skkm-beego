package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ShiftTable_20200508_204046 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ShiftTable_20200508_204046{}
	m.Created = "20200508_204046"

	migration.Register("ShiftTable_20200508_204046", m)
}

// Run the migrations
func (m *ShiftTable_20200508_204046) Up() {
	m.SQL("CREATE TABLE shift(id serial PRIMARY KEY, kkm_id int NOT NULL,sales bigint,sales_return bigint,sales_receipts bigint,sale bigint,sales_return_checks bigint,return_of_sale bigint,payouts bigint,making bigint,cash bigint,income bigint, is_open_shift boolean, shift_opening timestamp without time zone, shift_closing timestamp without time zone, FOREIGN KEY (kkm_id) REFERENCES kkm(id) ON DELETE CASCADE)")

}

// Reverse the migrations
func (m *ShiftTable_20200508_204046) Down() {
	m.SQL("drop table shift")

}
