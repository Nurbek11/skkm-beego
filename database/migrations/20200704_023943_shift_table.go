package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ShiftTable_20200704_023943 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ShiftTable_20200704_023943{}
	m.Created = "20200704_023943"

	migration.Register("ShiftTable_20200704_023943", m)
}

// Run the migrations
func (m *ShiftTable_20200704_023943) Up() {
	m.SQL("CREATE TABLE shift(id serial PRIMARY KEY, kkm_id int NOT NULL,shift_number int,document_number int,payouts VARCHAR(255),depositing VARCHAR(255),withdrawing VARCHAR(255),income VARCHAR(255), is_open_shift boolean, shift_opening varchar(255), shift_closing timestamp without time zone)")

}

// Reverse the migrations
func (m *ShiftTable_20200704_023943) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
