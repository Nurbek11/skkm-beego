package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type KkmTable_20200702_151236 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &KkmTable_20200702_151236{}
	m.Created = "20200702_151236"

	migration.Register("KkmTable_20200702_151236", m)
}

// Run the migrations
func (m *KkmTable_20200702_151236) Up() {
	m.SQL("CREATE TABLE kkm(id serial PRIMARY KEY,ofd_id int,ofd_req_num int,ofd_token int,shift_number int,shift_open_date timestamp without time zone,shift_closed bool,printed_number int,serial_number VARCHAR(255),fns_kkm_id VARCHAR(255),bin VARCHAR(255),address VARCHAR(255), title VARCHAR(255),cash VARCHAR(255),password VARCHAR(255), organization_bin VARCHAR(255))")

}

// Reverse the migrations
func (m *KkmTable_20200702_151236) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
