package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// product

// ProductGet .
func ProductGet(w http.ResponseWriter, r *http.Request) {
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
	SQLQuery := " from `" + productTable + "`"
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
			count, _, _ := selectProcess("select count(*) as ctn from `" + productTable + "` where " + where)
			pagination["total_count"] = count[0]["ctn"]
		} else {
			count, _, _ := selectProcess("select count(*) as ctn from `" + productTable + "`")
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

// ProductAdd .
func ProductAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	fieldCheck := requiredFiledsCheck(body, productRequiredFields)
	if len(fieldCheck) > 0 {
		SetReponseStatus(w, r, statusCodeBadRequest, fieldCheck+" required", dialogType, response)
		return
	}

	body["created_date_time"] = time.Now().UTC().String()
	body["modified_date_time"] = body["created_date_time"]
	body["status"] = "1"

	var status string
	var ok bool
	for true {
		body["product_u_id"] = RandStringBytes(productDigits)
		status, ok = insertSQL(productTable, body)
		if !strings.EqualFold(status, statusCodeDuplicateEntry) {
			break
		}
	}
	w.Header().Set("Status", status)
	if ok {
		if strings.EqualFold(body["track"], "1") && !strings.EqualFold(body["stock"], "0") {
			insertSQL(productStockTable, map[string]string{
				"store_u_id":         body["store_u_id"],
				"product_u_id":       body["product_u_id"],
				"quantity":           body["stock"],
				"type":               "1",
				"status":             "1",
				"created_date_time":  body["created_date_time"],
				"modified_date_time": body["created_date_time"],
				"created_by":         body["created_by"],
				"modified_by":        body["modified_by"],
			})
		}
		response["product_u_id"] = body["product_u_id"]
		response["meta"] = setMeta(status, "Product Row added", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

// ProductUpdate .
func ProductUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range productRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	body["modified_date_time"] = time.Now().UTC().String()

	if !strings.EqualFold(body["added"], "0") {
		if strings.Contains(body["added"], "-") {
			insertSQL(productStockTable, map[string]string{
				"store_u_id":         body["store_u_id"],
				"product_u_id":       body["product_u_id"],
				"quantity":           body["added"],
				"type":               "3",
				"status":             "1",
				"created_date_time":  body["modified_date_time"],
				"modified_date_time": body["modified_date_time"],
				"created_by":         body["modified_by"],
				"modified_by":        body["modified_by"],
			})
		} else {
			insertSQL(productStockTable, map[string]string{
				"store_u_id":         body["store_u_id"],
				"product_u_id":       body["product_u_id"],
				"quantity":           body["added"],
				"type":               "1",
				"status":             "1",
				"created_date_time":  body["modified_date_time"],
				"modified_date_time": body["modified_date_time"],
				"created_by":         body["modified_by"],
				"modified_by":        body["modified_by"],
			})
		}
	}
	delete(body, "added")

	status, ok := updateSQL(productTable, r.URL.Query(), body)
	w.Header().Set("Status", status)
	if ok {
		response["meta"] = setMeta(status, "Product Row updated", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}
