package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// customer

// CustomerGet .
func CustomerGet(w http.ResponseWriter, r *http.Request) {
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
	SQLQuery := " from `" + customerTable + "`"
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
			count, _, _ := selectProcess("select count(*) as ctn from `" + customerTable + "` where " + where)
			pagination["total_count"] = count[0]["ctn"]
		} else {
			count, _, _ := selectProcess("select count(*) as ctn from `" + customerTable + "`")
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

// CustomerAdd .
func CustomerAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	fieldCheck := requiredFiledsCheck(body, customerRequiredFields)
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
		body["customer_u_id"] = RandStringBytes(customerDigits)
		status, ok = insertSQL(customerTable, body)
		if !strings.EqualFold(status, statusCodeDuplicateEntry) {
			break
		}
	}
	w.Header().Set("Status", status)
	if ok {
		response["customer_u_id"] = body["customer_u_id"]
		response["meta"] = setMeta(status, "Customer Row added", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}

// CustomerUpdate .
func CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	body := map[string]string{}

	// file upload
	r.ParseMultipartForm(32 << 20)
	for key, value := range r.Form {
		body[key] = value[0]
	}

	for _, field := range customerRequiredFields {
		if _, ok := body[field]; ok {
			if len(body[field]) == 0 {
				SetReponseStatus(w, r, statusCodeBadRequest, field+" required", dialogType, response)
				return
			}
		}
	}

	body["modified_date_time"] = time.Now().UTC().String()

	status, ok := updateSQL(customerTable, r.URL.Query(), body)
	w.Header().Set("Status", status)
	if ok {
		response["meta"] = setMeta(status, "Customer Row updated", dialogType)
	} else {
		response["meta"] = setMeta(status, "", dialogType)
	}

	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	json.NewEncoder(w).Encode(response)
}
