package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ShiftTable_20200523_220530 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ShiftTable_20200523_220530{}
	m.Created = "20200523_220530"

	migration.Register("ShiftTable_20200523_220530", m)
}

// Run the migrations
func (m *ShiftTable_20200523_220530) Up() {
	m.SQL("CREATE TABLE shift(id serial PRIMARY KEY, kkm_id int NOT NULL,payouts VARCHAR(255),making VARCHAR(255),income VARCHAR(255), is_open_shift boolean, shift_opening timestamp without time zone, shift_closing timestamp without time zone, FOREIGN KEY (kkm_id) REFERENCES kkm(id) ON DELETE CASCADE)")

}

// Reverse the migrations
func (m *ShiftTable_20200523_220530) Down() {
	m.SQL("drop table shift")

}
