package main

var dbConfig string
var connectionPool int

var categoryTable = "categories"
var customerTable = "customers"
var customerAmountTable = "customer_amounts"
var productTable = "products"
var saleTable = "sales"
var saleProductTable = "sale_products"
var productStockTable = "product_stocks"
var subcategoryTable = "sub_categories"
var storeTable = "stores"
var userTable = "users"

var test bool
var migrate bool

var categoryDigits = 5
var customerDigits = 15
var productDigits = 10
var saleDigits = 15
var subcategoryDigits = 6
var storeDigits = 5
var userDigits = 6

var defaultLimit = "25"
var defaultOffset = "0"

var androidLive = "a2t3K5Y8e2W7Z5T2"
var androidTest = "E8y6S5H5T4e9q7q7"
var iOSLive = "b4E6U9K8j6b5E9W3"
var iOSTest = "R4n7N8G4m9B4S5n2"
var cron = "ZNPZTTDEVAStYczW"

// for checking unauth request
var apikeys = map[string]string{
	androidLive: "1", // android live
	androidTest: "1", // android test
	iOSLive:     "1", // ios live
	iOSTest:     "1", // ios test
	cron:        "1", // cron
}

var dialogType = "1"
var toastType = "2"
var appUpdateAvailable = "3"
var appUpdateRequired = "4"

// required fields
var categoryRequiredFields = []string{
	"store_u_id", "name",
}
var customerRequiredFields = []string{
	"store_u_id",
}
var customerAmountRequiredFields = []string{
	"store_u_id", "customer_u_id", "amount", "type",
}
var productRequiredFields = []string{
	"store_u_id", "name",
}
var saleRequiredFields = []string{
	"store_u_id",
}
var saleProductRequiredFields = []string{
	"store_u_id",
}
var productStockRequiredFields = []string{
	"store_u_id", "product_u_id", "quantity", "type", "total",
}
var subcategoryRequiredFields = []string{
	"store_u_id", "name",
}
var storeRequiredFields = []string{}
var userRequiredFields = []string{}

// server codes
var statusCodeOk = "200"
var statusCodeCreated = "201"
var statusCodeBadRequest = "400"
var statusCodeForbidden = "403"
var statusCodeServerError = "500"
var statusCodeDuplicateEntry = "1000"

// versions
var iOSVersionCode = 1.0
var iOSForceVersionCode = 1.0

var androidVersionCode = 1.0
var androidForceVersionCode = 1.0
