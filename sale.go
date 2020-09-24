package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// sale

// SaleGet .
func SaleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	params := r.URL.Query()
	limitOffset := " "

	if _, ok := params["limit"]; ok {
		limitOffset += " limit " + params["limit"][0]
		delete(params, "limit")
	}

	offset := defaultOffset
	if _, ok := params["offset"]; ok {
		limitOffset += " offset " + params["offset"][0]
		offset = params["offset"][0]
		delete(params, "offset")
	}

	orderBy := " "

	if _, ok := params["orderby"]; ok {
		orderBy += " order by " + params["orderby"][0]
		delete(params, "orderby")
		if _, ok := params["sortby"]; ok {
			orderBy += " " + params["sortby"][0] + " "
			delete(params, "sortby")
		} else {
			orderBy += " asc "
		}
	} else {
		orderBy += " order by created_date_time desc "
	}

	resp := " * "
	if _, ok := params["resp"]; ok {
		resp = " " + params["resp"][0] + " "
		delete(params, "resp")
	}

	where := ""
	init := false
	for key, val := range params {
		if init {
			where += " and "
		}
		where += " `" + key + "` = '" + val[0] + "' "
		init = true
	}
	SQLQuery := " from `" + saleTable + "`"
	if strings.Compare(where, "") != 0 {
		SQLQuery += " where " + where
	}
	SQLQuery += orderBy
	SQLQuery += limitOffset

	data, status, ok := selectProcess("select " + resp + SQLQuery)
	w.Header().Set("Status", status)
	if ok {
		response["data"] = data

		pagination := map[string]string{}
		if len(where) > 0 {
			count, _, _ := selectProcess("select count(*) as ctn from `" + saleTable + "` where " + where)
			pagination["total_count"] = count[0]["ctn"]
		} else {
			count, _, _ := selectProcess("select count(*) as ctn from `" + saleTable + "`")
			pagination["total_count"] = count[0]["ctn"]
		}
		pagination["count"] = strconv.Itoa(len(data))
		pagination["offset"] = offset
		response["pagination"] = pagination

		response["meta"] = setMeta(status, "ok", "")
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

// SaleAdd .
func SaleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	fieldCheck := requiredFiledsCheck(body, saleRequiredFields)
	if len(fieldCheck) > 0 {
		SetReponseStatus(w, r, statusCodeBadRequest, fieldCheck+" required", dialogType, response)
		return
	}

	body["created_date_time"] = time.Now().UTC().String()
	body["modified_date_time"] = body["created_date_time"]
	body["status"] = "1"

	productsString := body["products"]
	delete(body, "products")

	var status string
	var ok bool
	for true {
		body["sale_u_id"] = RandStringBytes(saleDigits)
		status, ok = insertSQL(saleTable, body)
		if !strings.EqualFold(status, statusCodeDuplicateEntry) {
			break
		}
	}
	w.Header().Set("Status", status)
	if ok {
		if len(body["customer_u_id"]) > 0 && len(body["pay_later"]) > 0 {
			insertSQL(customerAmountTable, map[string]string{
				"store_u_id":         body["store_u_id"],
				"sale_u_id":          body["sale_u_id"],
				"customer_u_id":      body["customer_u_id"],
				"amount":             body["pay_later"],
				"type":               "1",
				"status":             "1",
				"created_date_time":  body["created_date_time"],
				"modified_date_time": body["created_date_time"],
			})
			db.Exec("update " + customerTable + " set amount = amount - " + body["pay_later"] + ", modified_date_time = '" + body["created_date_time"] + "' where store_u_id = '" + body["store_u_id"] + "' and customer_u_id = '" + body["customer_u_id"] + "'")
		}
		response["sale_u_id"] = body["sale_u_id"]

		//
		products := []map[string]string{}
		json.Unmarshal([]byte(productsString), &products)
		for _, product := range products {
			product["sale_u_id"] = body["sale_u_id"]
			product["created_date_time"] = body["created_date_time"]
			product["modified_date_time"] = body["created_date_time"]
			product["status"] = "1"
			if strings.EqualFold(product["track"], "1") {
				insertSQL(productStockTable, map[string]string{
					"store_u_id":         product["store_u_id"],
					"sale_u_id":          body["sale_u_id"],
					"product_u_id":       product["product_u_id"],
					"quantity":           product["quantity"],
					"type":               "2",
					"status":             "1",
					"created_date_time":  body["created_date_time"],
					"modified_date_time": body["created_date_time"],
				})
				db.Exec("update " + productTable + " set stock = stock - " + product["quantity"] + ", modified_date_time = '" + body["created_date_time"] + "' where store_u_id = '" + product["store_u_id"] + "' and product_u_id = '" + product["product_u_id"] + "'")
			}
			delete(product, "track")
			insertSQL(saleProductTable, product)
		}
		//

		response["meta"] = setMeta(status, "Sale Row added", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

// SaleUpdate .
func SaleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range saleRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	body["modified_date_time"] = time.Now().UTC().String()

	status, ok := updateSQL(saleTable, r.URL.Query(), body)
	w.Header().Set("Status", status)
	if ok {
		response["meta"] = setMeta(status, "Sale Row updated", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}
