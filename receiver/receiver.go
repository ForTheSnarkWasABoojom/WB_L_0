package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	lib "Golang/WB_0/receiver/lib"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

var orderSlice []lib.Order

func AddOrderToDB(order lib.Order) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	insertDelivery := "INSERT INTO deliveries VALUES (" + order.GetDelivery() + ");"
	_, err = db.Exec(insertDelivery)
	CheckError(err)

	insertPayment := "INSERT INTO payments VALUES (" + order.GetPayment() + ");"
	_, err = db.Exec(insertPayment)
	CheckError(err)

	insertOrder := "INSERT INTO orders VALUES (" + order.GetOrder() + ");"
	_, err = db.Exec(insertOrder)
	CheckError(err)

	for i := range order.Items {
		insertItem := "INSERT INTO items VALUES (" + order.GetItem(i) + ");"
		_, err = db.Exec(insertItem)
		CheckError(err)
	}

	fmt.Println("Data inserted in db!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	RetrieveDataFromDb()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		HandleNats()
		defer wg.Done()
	}()
	wg.Add(1)
	go func() {
		ProvideService()
		defer wg.Done()
	}()

	wg.Wait()
}

func ProvideService() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	exit_chan := make(chan int)

	go func() {
		for {
			mux := http.NewServeMux()
			mux.HandleFunc("/order/", ProvideGetOrder)
			mux.HandleFunc("/delivery/", ProvideGetDelivery)
			mux.HandleFunc("/payment/", ProvideGetPayment)
			mux.HandleFunc("/item/", ProvideGetItem)
			fmt.Println("Server has started!")
			http.ListenAndServe("localhost:8080", mux)
		}
	}()

	go func() {
		s := <-signalChanel
		switch s {

		default:
			fmt.Println("Server has stopped!")
			exit_chan <- 1
		}
	}()
	exitCode := <-exit_chan
	os.Exit(exitCode)
}

func ProvideGetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		defer r.Body.Close()
		idStr := strings.TrimPrefix(r.URL.Path, "/order/")

		for _, order := range orderSlice {
			if order.Order_uid == idStr {
				w.WriteHeader((http.StatusOK))
				data, err := json.Marshal(order)
				CheckError(err)
				w.Write([]byte(data))
				return
			}
		}
		w.WriteHeader((http.StatusNotFound))
		w.Write([]byte("Order с данным id не найден!"))
		return

	}

	w.WriteHeader(http.StatusBadRequest)
}
func ProvideGetDelivery(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		defer r.Body.Close()
		idStr := strings.TrimPrefix(r.URL.Path, "/delivery/")

		for _, order := range orderSlice {
			if order.Delivery_id == idStr {
				w.WriteHeader((http.StatusOK))
				data, err := json.Marshal(order.Delivery)
				CheckError(err)
				w.Write([]byte(data))
				return
			}
		}
		w.WriteHeader((http.StatusNotFound))
		w.Write([]byte("Delivery с данным id не найден!"))
		return

	}

	w.WriteHeader(http.StatusBadRequest)
}
func ProvideGetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		defer r.Body.Close()
		idStr := strings.TrimPrefix(r.URL.Path, "/payment/")

		for _, order := range orderSlice {
			if order.Payment_id == idStr {
				w.WriteHeader((http.StatusOK))
				data, err := json.Marshal(order.Payment)
				CheckError(err)
				w.Write([]byte(data))
				return
			}
		}
		w.WriteHeader((http.StatusNotFound))
		w.Write([]byte("Payment с данным id не найден!"))
		return

	}

	w.WriteHeader(http.StatusBadRequest)
}
func ProvideGetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		defer r.Body.Close()
		idStr := strings.TrimPrefix(r.URL.Path, "/item/")

		for _, order := range orderSlice {
			i, _ := strconv.Atoi(idStr)
			for _, item := range order.Items {
				if item.Chrt_id == i {
					w.WriteHeader((http.StatusOK))
					data, err := json.Marshal(order.Items[0])
					CheckError(err)
					w.Write([]byte(data))
					return
				}
			}
		}
		w.WriteHeader((http.StatusNotFound))
		w.Write([]byte("Item с данным id не найден!"))
		return

	}

	w.WriteHeader(http.StatusBadRequest)
}

func RetrieveDataFromDb() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	rowsRs, err := db.Query("SELECT order_uid FROM orders")
	defer rowsRs.Close()

	var order_uids []string
	for rowsRs.Next() {
		var order_uid string
		err = rowsRs.Scan(&order_uid)
		CheckError(err)
		order_uids = append(order_uids, order_uid)
	}

	for _, el := range order_uids {
		rowsRs, err := db.Query("SELECT * from orders where order_uid = '" + el + "';")
		CheckError(err)
		defer rowsRs.Close()

		for rowsRs.Next() {
			var order lib.Order

			err = rowsRs.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &order.Delivery_id, &order.Payment_id,
				&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
				&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
			CheckError(err)

			var payment lib.Payment

			paymentRows, err := db.Query("SELECT * from payments where transaction = '" + order.Payment_id + "';")
			CheckError(err)
			defer paymentRows.Close()
			for paymentRows.Next() {
				err = paymentRows.Scan(&payment.Transaction, &payment.Request_id, &payment.Currency, &payment.Provider,
					&payment.Amount, &payment.Payment_dt, &payment.Bank, &payment.Delivery_cost, &payment.Goods_total,
					&payment.Custom_fee)
				CheckError(err)
			}
			order.Payment = payment

			var delivery lib.Delivery

			deliveryRows, err := db.Query("SELECT * from deliveries where delivery_id = '" + order.Delivery_id + "';")
			CheckError(err)
			defer deliveryRows.Close()
			for deliveryRows.Next() {
				err = deliveryRows.Scan(&delivery.Delivery_id, &delivery.Name, &delivery.Phone, &delivery.Zip,
					&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
				CheckError(err)
			}
			order.Delivery = delivery

			var itemSlice []lib.Item
			itemRows, err := db.Query("SELECT * from items where order_uid = '" + order.Order_uid + "';")
			CheckError(err)
			defer itemRows.Close()
			for itemRows.Next() {
				var item lib.Item
				err = itemRows.Scan(&item.Chrt_id, &item.Track_number, &item.Price, &item.Rid, &item.Name,
					&item.Sale, &item.Size, &item.Total_price, &item.Nm_id, &item.Brand, &item.Status, &item.Order_uid)
				CheckError(err)
				itemSlice = append(itemSlice, item)
			}
			order.Items = itemSlice
			orderSlice = append(orderSlice, order)
		}
	}
	fmt.Println("Data from DB was retrieved!")
}
func HandleNats() {
	nc, err := nats.Connect("demo.nats.io")
	CheckError(err)
	defer nc.Close()

	sub, err := nc.SubscribeSync("updates")
	CheckError(err)
	fmt.Println("Connected to Nats! Waiting for the message")

	for {
		msg, err := sub.NextMsg(10 * time.Minute)
		CheckError(err)

		var order lib.Order
		err = json.Unmarshal(msg.Data, &order)
		CheckError(err)

		AddOrderToDB(order)
		orderSlice = append(orderSlice, order)
	}
}
