package models

type ChequeData struct {
	Password string `orm:"null"`
	Cheque   struct {
		Domain          string
		Markup          uint64
		TotalSum        string
		ChangeMoney     string
		TotalDiscount   string
		TotalCharge     string
		NDS             string
		PaymentType     string
		TaxType         string
		TaxationType    string
		TaxPercent      uint32
		TaxSum          uint64
		TaxIsInTotalSum bool
		CustomerEmail   string
		CustomerPhone   string

		Goods []struct {
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
			DiscountSum     int
			Discount        string
			DiscountStorno  bool
			MarkupSum       uint64
			MarkupStorno    bool
			Total           uint64
			IsStorno        bool
		}
	}
}
