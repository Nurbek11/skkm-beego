package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ZreportTable_20200601_121719 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ZreportTable_20200601_121719{}
	m.Created = "20200601_121719"

	migration.Register("ZreportTable_20200601_121719", m)
}

// Run the migrations
func (m *ZreportTable_20200601_121719) Up() {
		m.SQL("CREATE TABLE zreport(id serial PRIMARY KEY,shift int,start_sales VARCHAR(255),start_payouts VARCHAR(255),start_sales_return VARCHAR(255),start_refunds VARCHAR(255),shift_sales VARCHAR(255),shift_payouts VARCHAR(255),shift_sales_return VARCHAR(255),shift_refunds VARCHAR(255), time_of_creation timestamp without time zone)")


}

// Reverse the migrations
func (m *ZreportTable_20200601_121719) Down() {
	m.SQL("DROP TABLE zreport")

}
