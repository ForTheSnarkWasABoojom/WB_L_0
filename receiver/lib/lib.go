package lib

import "fmt"

type Delivery struct {
	Delivery_id string
	Name        string
	Phone       string
	Zip         int
	City        string
	Address     string
	Region      string
	Email       string
}

type Payment struct {
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        int
	Payment_dt    int
	Bank          string
	Delivery_cost int
	Goods_total   int
	Custom_fee    int
}

type Item struct {
	Chrt_id      int
	Track_number string
	Price        int
	Rid          string
	Name         string
	Sale         int
	Size         int
	Total_price  int
	Nm_id        int
	Brand        string
	Status       int
	Order_uid    string
}

type Order struct {
	Order_uid          string
	Track_number       string
	Entry              string
	Delivery_id        string
	Payment_id         string
	Locale             string
	Internal_signature string
	Customer_id        string
	Delivery_service   string
	Shardkey           int
	Sm_id              int
	Date_created       string
	Oof_shard          int
	Delivery           Delivery
	Payment            Payment
	Items              []Item
}

func (order Order) GetOrder() string {
	return fmt.Sprintf(`'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%d' `,
		order.Order_uid, order.Track_number, order.Entry, order.Delivery_id, order.Payment_id,
		order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service,
		order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
}
func (order Order) GetDelivery() string {
	return fmt.Sprintf(`'%s', '%s', '%s', '%d', '%s', '%s', '%s', '%s'`,
		order.Delivery.Delivery_id, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
}
func (order Order) GetPayment() string {
	return fmt.Sprintf(`'%s', '%s', '%s', '%s', '%d', '%d', '%s', '%d', '%d', '%d' `,
		order.Payment.Transaction, order.Payment.Request_id, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.Payment_dt, order.Payment.Bank, order.Payment.Delivery_cost,
		order.Payment.Goods_total, order.Payment.Custom_fee)
}
func (order Order) GetItem(id int) string {
	return fmt.Sprintf(`'%d', '%s', '%d', '%s', '%s', '%d', '%d', '%d', '%d', '%s', '%d', '%s' `,
		order.Items[id].Chrt_id, order.Items[id].Track_number, order.Items[id].Price, order.Items[id].Rid,
		order.Items[id].Name, order.Items[id].Sale, order.Items[id].Size, order.Items[id].Total_price,
		order.Items[id].Nm_id, order.Items[id].Brand, order.Items[id].Status, order.Items[id].Order_uid)
}
