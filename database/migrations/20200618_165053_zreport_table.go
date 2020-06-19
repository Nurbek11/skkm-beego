package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ZreportTable_20200618_165053 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ZreportTable_20200618_165053{}
	m.Created = "20200618_165053"

	migration.Register("ZreportTable_20200618_165053", m)
}

// Run the migrations
func (m *ZreportTable_20200618_165053) Up() {
	m.SQL("CREATE TABLE zreport(id serial PRIMARY KEY,cashier_id int,organization_id int,shift_id int,cash VARCHAR(255),start_sales VARCHAR(255),start_payouts VARCHAR(255),start_sales_return VARCHAR(255),start_refunds VARCHAR(255),shift_sales VARCHAR(255),shift_payouts VARCHAR(255),shift_sales_return VARCHAR(255),shift_refunds VARCHAR(255), time_of_creation timestamp without time zone)")

}

// Reverse the migrations
func (m *ZreportTable_20200618_165053) Down() {
	m.SQL("DROP TABLE zreport")

}
