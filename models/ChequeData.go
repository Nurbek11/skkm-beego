package models

type ChequeData struct {
	Password string `orm:"null"`
	Cheque struct {
		TotalSum      string
		ChangeMoney   string
		TotalDiscount string
		TotalCharge   string
		NDS           string
		PaymentType   string
		Goods         []struct {
			GoodId          int
			GoodTitle       string
			GoodPrice       string
			GoodDiscount    string
			GoodExtraCharge string
			GoodNumber      string
			GoodSum         string
			IsDisPrice      bool
			IsDisDiscount   bool
			IsDisExCharge   bool
			IsDisNumber     bool
		}
	}
}
