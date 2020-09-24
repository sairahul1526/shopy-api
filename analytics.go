package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Analytics .
func Analytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// pies
	pies := []map[string]interface{}{}

	// user active and expired
	result, status, ok := selectProcess("select sum(cash) as cash_total, sum(card) as card_total, sum(`check`) as check_total, sum(voucher) as voucher_total, sum(store_credit) as store_credit_total, sum(paytm) as paytm_total, sum(pay_later) as pay_later_total, sum(other) as other_total from " + saleTable + " where date(created_date_time) >= '" + r.FormValue("from") + "' and date(created_date_time) <= '" + r.FormValue("to") + "' and store_u_id = '" + r.FormValue("store_u_id") + "' and status = 1")
	if !ok {
		SetReponseStatus(w, r, status, "", dialogType, response)
		return
	}
	if len(result) > 0 {
		pies = append(pies, map[string]interface{}{
			"title": "Payment methods",
			"type":  "1",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Payments",
					"type":  "1",
					"data": []map[string]string{
						map[string]string{
							"title": "Cash",
							"shown": result[0]["cash_total"] + " (Cash Credit, Store)",
							"value": strings.Split(result[0]["cash_total"], ".")[0],
							"color": "#AED6F1",
						},
						map[string]string{
							"title": "Card",
							"shown": result[0]["card_total"],
							"value": strings.Split(result[0]["card_total"], ".")[0],
							"color": "#A2D9CE",
						},
						map[string]string{
							"title": "Check",
							"shown": result[0]["check_total"],
							"value": strings.Split(result[0]["check_total"], ".")[0],
							"color": "#F5B7B1",
						},
						map[string]string{
							"title": "Voucher",
							"shown": result[0]["voucher_total"],
							"value": strings.Split(result[0]["voucher_total"], ".")[0],
							"color": "#878a88",
						},
						map[string]string{
							"title": "Store\nCredit",
							"shown": result[0]["store_credit_total"],
							"value": strings.Split(result[0]["store_credit_total"], ".")[0],
							"color": "#b55bc9",
						},
						map[string]string{
							"title": "Paytm",
							"shown": result[0]["paytm_total"],
							"value": strings.Split(result[0]["paytm_total"], ".")[0],
							"color": "#4bd1e3",
						},
						map[string]string{
							"title": "Pay\nLater",
							"shown": result[0]["pay_later_total"],
							"value": strings.Split(result[0]["pay_later_total"], ".")[0],
							"color": "#e3974b",
						},
						map[string]string{
							"title": "Other",
							"shown": result[0]["other_total"],
							"value": strings.Split(result[0]["other_total"], ".")[0],
							"color": "#000000",
						},
					},
				},
			},
		})
	}

	// result, status, ok = selectProcess("select sum(cash) as cash_total, sum(card) as card_total, sum(`check`) as check_total, sum(voucher) as voucher_total, sum(store_credit) as store_credit_total, sum(paytm) as paytm_total, sum(other) as other_total from " + saleTable + " where date(created_date_time) >= '" + r.FormValue("from") + "' and date(created_date_time) <= '" + r.FormValue("to") + "' and store_u_id = '" + r.FormValue("store_u_id") + "' and status = 1")
	// if !ok {
	// 	SetReponseStatus(w, r, status, "", dialogType, response)
	// 	return
	// }
	// if len(result) > 0 {
	// 	// check bar
	// 	pies = append(pies, map[string]interface{}{
	// 		"title": "Revenue",
	// 		"type":  "2",
	// 		"data": []map[string]interface{}{
	// 			map[string]interface{}{
	// 				"title": "",
	// 				"type":  "2",
	// 				"data": []map[string]string{
	// 					map[string]string{
	// 						"title": "Cash Credit, Store",
	// 						"shown": "Cash Credit, Store",
	// 						"value": "5158",
	// 						"color": "#AED6F1",
	// 					},
	// 					map[string]string{
	// 						"title": "Net Sales",
	// 						"shown": "",
	// 						"value": "7718",
	// 						"color": "#AED6F1",
	// 					},
	// 					map[string]string{
	// 						"title": "Total Order Conunt",
	// 						"shown": "Count",
	// 						"value": "2626",
	// 						"color": "#AED6F1",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	})
	// }

	// check time series
	// monthly profit
	result, status, ok = selectProcess("select sum(sale-cost) as profit, date(created_date_time) as monthly_date from " + saleProductTable + " where date(created_date_time) >= '" + r.FormValue("from") + "' and date(created_date_time) <= '" + r.FormValue("to") + "' and store_u_id = '" + r.FormValue("store_u_id") + "' and cost > 0 and status = 1 group by month(created_date_time)")
	if !ok {
		SetReponseStatus(w, r, status, "", dialogType, response)
		return
	}
	if len(result) > 0 {
		profits := []map[string]string{}
		for _, profit := range result {
			profits = append(profits, map[string]string{
				"title": profit["monthly_date"],
				"shown": "",
				"value": profit["profit"],
				"color": "#000000",
			})
		}
		pies = append(pies, map[string]interface{}{
			"title": "Monthly Profit",
			"type":  "3",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Monthly Profit",
					"type":  "3",
					"data":  profits,
				},
			},
		})
	}
	// pies = append(pies, map[string]interface{}{
	// 	"title": "Monthly Profit",
	// 	"type":  "3",
	// 	"data": []map[string]interface{}{
	// 		map[string]interface{}{
	// 			"title": "Monthly Profit",
	// 			"type":  "3",
	// 			"data": []map[string]string{
	// 				map[string]string{
	// 					"title": "2020-01-01",
	// 					"shown": "",
	// 					"value": "51",
	// 					"color": "#000000",
	// 				},
	// 				map[string]string{
	// 					"title": "2020-02-01",
	// 					"shown": "",
	// 					"value": "40",
	// 					"color": "#000000",
	// 				},
	// 				map[string]string{
	// 					"title": "2020-03-01",
	// 					"shown": "",
	// 					"value": "120",
	// 					"color": "#AED6F1",
	// 				},
	// 				map[string]string{
	// 					"title": "2020-04-01",
	// 					"shown": "",
	// 					"value": "89",
	// 					"color": "#A2D9CE",
	// 				},
	// 				map[string]string{
	// 					"title": "2020-05-01",
	// 					"shown": "",
	// 					"value": "51",
	// 					"color": "#A2D9CE",
	// 				},
	// 			},
	// 		},
	// 	},
	// })

	// // grid
	// revenue breakdown
	result, status, ok = selectProcess("select count(*) as order_count, round(sum(subtotal)) as total_gross, round(sum(discount)) as total_discount, round(sum(total)) as total_net from " + saleTable + " where date(created_date_time) >= '" + r.FormValue("from") + "' and date(created_date_time) <= '" + r.FormValue("to") + "' and store_u_id = '" + r.FormValue("store_u_id") + "' and status = 1")
	if !ok {
		SetReponseStatus(w, r, status, "", dialogType, response)
		return
	}
	if len(result) > 0 {
		pies = append(pies, map[string]interface{}{
			"title": "Revenue Breakdown",
			"type":  "4",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Revenue Breakdown",
					"type":  "4",
					"data": []map[string]string{
						map[string]string{
							"title": "No. Orders",
							"shown": "No. Orders",
							"value": result[0]["order_count"],
							"color": "#000000",
						},
						map[string]string{
							"title": "Gross Sales",
							"shown": "Gross Sales",
							"value": result[0]["total_gross"],
							"color": "#000000",
						},
						map[string]string{
							"title": "Discount",
							"shown": "Discount",
							"value": result[0]["total_discount"],
							"color": "#000000",
						},
						map[string]string{
							"title": "Net Sales",
							"shown": "Net Sales",
							"value": result[0]["total_net"],
							"color": "#000000",
						},
					},
				},
			},
		})
	}

	// check list
	// top products
	result, status, ok = selectProcess("select name, round(sum(quantity)) as total_quantity, round(sum(quantity*sale)) as total_amount from " + saleProductTable + " where date(created_date_time) >= '" + r.FormValue("from") + "' and date(created_date_time) <= '" + r.FormValue("to") + "' and store_u_id = '" + r.FormValue("store_u_id") + "' and status = 1 group by product_u_id order by total_amount desc limit 5")
	if !ok {
		SetReponseStatus(w, r, status, "", dialogType, response)
		return
	}
	if len(result) > 0 {
		products := []map[string]string{}
		for _, product := range result {
			products = append(products, map[string]string{
				"title": product["name"] + " X " + product["total_quantity"],
				"shown": product["name"],
				"value": product["total_amount"],
				"color": "#000000",
			})
		}
		pies = append(pies, map[string]interface{}{
			"title": "Top Selling Products",
			"type":  "5",
			"data": []map[string]interface{}{
				map[string]interface{}{
					"title": "Revenue Breakdown",
					"type":  "5",
					"data":  products,
				},
			},
		})
	}

	response["graphs"] = pies
	response["meta"] = setMeta(statusCodeOk, "ok", "")

	fmt.Println(pies)
	w.WriteHeader(getHTTPStatusCode(response["meta"].(map[string]string)["status"]))
	meta, required := checkAppUpdate(r)
	if required {
		response["meta"] = meta
	}
	json.NewEncoder(w).Encode(response)
}
