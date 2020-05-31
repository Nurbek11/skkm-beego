package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ShiftTable_20200531_113324 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ShiftTable_20200531_113324{}
	m.Created = "20200531_113324"

	migration.Register("ShiftTable_20200531_113324", m)
}

// Run the migrations
func (m *ShiftTable_20200531_113324) Up() {
	m.SQL("CREATE TABLE shift(id serial PRIMARY KEY, kkm_id int NOT NULL,payouts VARCHAR(255),making VARCHAR(255),income VARCHAR(255), is_open_shift boolean, shift_opening timestamp without time zone, shift_closing timestamp without time zone)")

}

// Reverse the migrations
func (m *ShiftTable_20200531_113324) Down() {
	  m.SQL("DROP TABLE shift")

}
