package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/akrylysov/algnhsa"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var db *sql.DB
var err error

var clients = map[string]map[*websocket.Conn]bool{}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func connectDatabase() {
	db, err = sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(connectionPool)
	db.SetMaxIdleConns(connectionPool)
	db.SetConnMaxLifetime(time.Hour)
}

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
}

func inits() {
	rand.Seed(time.Now().UnixNano())
	connectDatabase()
}

func main() {

	dbConfig = os.Getenv("dbConfig")
	connectionPool, _ = strconv.Atoi(os.Getenv("connectionPool"))
	test, _ = strconv.ParseBool(os.Getenv("test"))
	migrate, _ = strconv.ParseBool(os.Getenv("migrate"))

	inits()
	defer db.Close()
	router := mux.NewRouter()

	// analytics
	router.Path("/analytics").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(Analytics)).Methods("GET")

	// category
	router.Path("/category").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(CategoryGet)).Methods("GET")
	router.Path("/category").HandlerFunc(checkHeaders(CategoryAdd)).Methods("POST")
	router.Path("/category").Queries(
		"category_u_id", "{category_u_id}",
	).HandlerFunc(checkHeaders(CategoryUpdate)).Methods("PUT")

	// customer
	router.Path("/customer").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(CustomerGet)).Methods("GET")
	router.Path("/customer").HandlerFunc(checkHeaders(CustomerAdd)).Methods("POST")
	router.Path("/customer").Queries(
		"customer_u_id", "{customer_u_id}",
	).HandlerFunc(checkHeaders(CustomerUpdate)).Methods("PUT")

	// customer amount
	router.Path("/customeramount").Queries(
		"store_u_id", "{store_u_id}",
		"customer_u_id", "{customer_u_id}",
	).HandlerFunc(checkHeaders(CustomerAmountGet)).Methods("GET")
	router.Path("/customeramount").HandlerFunc(checkHeaders(CustomerAmountAdd)).Methods("POST")
	router.Path("/customeramount").Queries(
		"id", "{id}",
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(CustomerAmountUpdate)).Methods("PUT")

	// product
	router.Path("/product").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(ProductGet)).Methods("GET")
	router.Path("/product").HandlerFunc(checkHeaders(ProductAdd)).Methods("POST")
	router.Path("/product").Queries(
		"product_u_id", "{product_u_id}",
	).HandlerFunc(checkHeaders(ProductUpdate)).Methods("PUT")

	// product stock
	router.Path("/productstock").Queries(
		"store_u_id", "{store_u_id}",
		"product_u_id", "{product_u_id}",
	).HandlerFunc(checkHeaders(ProductStockGet)).Methods("GET")
	router.Path("/productstock").HandlerFunc(checkHeaders(ProductStockAdd)).Methods("POST")
	router.Path("/productstock").Queries(
		"id", "{id}",
		"product_u_id", "{product_u_id}",
	).HandlerFunc(checkHeaders(ProductStockUpdate)).Methods("PUT")

	// sale
	router.Path("/sale").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(SaleGet)).Methods("GET")
	router.Path("/sale").HandlerFunc(checkHeaders(SaleAdd)).Methods("POST")
	router.Path("/sale").Queries(
		"sale_u_id", "{sale_u_id}",
	).HandlerFunc(checkHeaders(SaleUpdate)).Methods("PUT")

	// sale
	router.Path("/saleproduct").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(SaleProductGet)).Methods("GET")
	router.Path("/saleproduct").Queries(
		"id", "{id}",
		"sale_u_id", "{sale_u_id}",
	).HandlerFunc(checkHeaders(SaleProductUpdate)).Methods("PUT")

	// subcategory
	router.Path("/subcategory").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(SubcategoryGet)).Methods("GET")
	router.Path("/subcategory").HandlerFunc(checkHeaders(SubcategoryAdd)).Methods("POST")
	router.Path("/subcategory").Queries(
		"subcategory_u_id", "{subcategory_u_id}",
	).HandlerFunc(checkHeaders(SubcategoryUpdate)).Methods("PUT")

	// store
	router.Path("/store").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(StoreGet)).Methods("GET")
	router.Path("/store").Queries(
		"name", "{name}",
	).HandlerFunc(checkHeaders(StoreGet)).Methods("GET")
	router.Path("/store").HandlerFunc(checkHeaders(StoreAdd)).Methods("POST")
	router.Path("/store").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(StoreUpdate)).Methods("PUT")

	// user
	router.Path("/user").Queries(
		"username", "{username}",
	).HandlerFunc(checkHeaders(UserGet)).Methods("GET")
	router.Path("/user").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(checkHeaders(UserGet)).Methods("GET")
	router.Path("/user").HandlerFunc(checkHeaders(UserAdd)).Methods("POST")
	router.Path("/user").Queries(
		"user_u_id", "{user_u_id}",
	).HandlerFunc(checkHeaders(UserUpdate)).Methods("PUT")

	// realtime
	router.Path("/realtime").Queries(
		"store_u_id", "{store_u_id}",
	).HandlerFunc(wsHandler).Methods("GET")

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	// go realtime()

	// fmt.Println(http.ListenAndServe(":5000", &WithCORS{router}))

	algnhsa.ListenAndServe(router, nil)
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,apikey,appversion,pkgname")

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(res, req)
}
