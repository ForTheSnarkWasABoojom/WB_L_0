package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	lib "Golang/WB_0/receiver/lib"

	"github.com/nats-io/nats.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	str := make([]rune, n)
	for i := range str {
		str[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(str)
}

func GenerateOrder() lib.Order {
	orderId := RandString(20)
	delivery := lib.Delivery{
		Delivery_id: RandString(20),
		Name:        RandString(20),
		Phone:       "+" + strconv.Itoa(rand.Intn(1000000000000)),
		Zip:         rand.Intn(10000),
		City:        RandString(20),
		Address:     RandString(20),
		Region:      RandString(20),
		Email:       RandString(20),
	}
	payment := lib.Payment{
		Transaction:   RandString(20),
		Request_id:    RandString(20),
		Currency:      RandString(20),
		Provider:      RandString(20),
		Amount:        rand.Intn(10000),
		Payment_dt:    rand.Intn(10000),
		Bank:          RandString(20),
		Delivery_cost: rand.Intn(10000),
		Goods_total:   rand.Intn(10000),
		Custom_fee:    rand.Intn(10000),
	}
	item := lib.Item{
		Chrt_id:      rand.Intn(10000),
		Track_number: RandString(20),
		Price:        rand.Intn(10000),
		Rid:          RandString(20),
		Name:         RandString(20),
		Sale:         rand.Intn(10000),
		Size:         rand.Intn(10000),
		Total_price:  rand.Intn(10000),
		Nm_id:        rand.Intn(10000),
		Brand:        RandString(20),
		Status:       rand.Intn(10000),
		Order_uid:    orderId,
	}
	var items []lib.Item
	items = append(items, item)
	order := lib.Order{
		Order_uid:          orderId,
		Track_number:       RandString(10),
		Entry:              RandString(10),
		Delivery_id:        delivery.Delivery_id,
		Payment_id:         payment.Transaction,
		Locale:             RandString(10),
		Internal_signature: RandString(10),
		Customer_id:        RandString(10),
		Delivery_service:   RandString(10),
		Shardkey:           rand.Intn(100),
		Sm_id:              rand.Intn(100),
		Date_created:       time.Now().Format(time.RFC3339),
		Oof_shard:          rand.Intn(100),
		Delivery:           delivery,
		Payment:            payment,
		Items:              items,
	}
	return order
}

func main() {
	order := GenerateOrder()

	nc, err := nats.Connect("demo.nats.io", nats.Name("API"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error marshaling struct: %v", err)
	}
	if err := nc.Publish("updates", jsonBytes); err != nil {
		log.Fatal(err)
	}
	log.Println("Send an order")
	log.Println(order.Date_created)
}
