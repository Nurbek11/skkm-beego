package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ChequeTable_20200703_085057 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ChequeTable_20200703_085057{}
	m.Created = "20200703_085057"

	migration.Register("ChequeTable_20200703_085057", m)
}

// Run the migrations
func (m *ChequeTable_20200703_085057) Up() {
	m.SQL("CREATE TABLE cheque(id serial PRIMARY KEY,kkm_id int,shift_id int,pos_transaction_id varchar(255),taken int,markup int,ticket_type varchar(255),money_placement_type varchar(255),domain varchar(255),tax_type varchar(255),taxation_type varchar(255),tax_percent varchar(255),tax_sum varchar(255),tax_is_in_total_sum bool,operator_name varchar(255),operator_code int,payment_cash int,payment_card int,payment_credit int,customer_email varchar(255),customer_phone varchar(255),ticket_number varchar(255),offline_ticket_number varchar(255),printed_ticket_number varchar(255),shift_document_number int,qr_code varchar(255),is_canceled bool,total_sum VARCHAR(255),change_money VARCHAR(255), total_discount VARCHAR(255), total_charge VARCHAR(255), n_d_s VARCHAR(255),payment_type VARCHAR(255),operation_type VARCHAR(255),date_time timestamp without time zone)")

}

// Reverse the migrations
func (m *ChequeTable_20200703_085057) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
