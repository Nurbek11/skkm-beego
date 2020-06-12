package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ShiftTable_20200605_221031 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ShiftTable_20200605_221031{}
	m.Created = "20200605_221031"

	migration.Register("ShiftTable_20200605_221031", m)
}

// Run the migrations
func (m *ShiftTable_20200605_221031) Up() {
	m.SQL("CREATE TABLE shift(id serial PRIMARY KEY, kkm_id int NOT NULL,payouts VARCHAR(255),depositing VARCHAR(255),withdrawing VARCHAR(255),income VARCHAR(255), is_open_shift boolean, shift_opening timestamp without time zone, shift_closing timestamp without time zone)")

}

// Reverse the migrations
func (m *ShiftTable_20200605_221031) Down() {
	m.SQL("DROP TABLE shift")

}