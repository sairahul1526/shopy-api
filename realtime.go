package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	if clients[r.FormValue("store_u_id")] == nil {
		clients[r.FormValue("store_u_id")] = map[*websocket.Conn]bool{}
	}
	clients[r.FormValue("store_u_id")][ws] = true
}

func realtime() {
	defer func() {
		go realtime()
	}()
	for {
		// products
		if len(clients) > 0 {
			products, _, ok := selectProcessNoLogging("select * from " + productTable + " where modified_date_time > '" + time.Now().Add(-10*time.Second).UTC().String() + "'")

			if ok && len(products) > 0 {
				for _, product := range products {
					for client := range clients[product["store_u_id"]] {
						client.WriteJSON(map[string]map[string]string{
							"product": product,
						})
					}
				}
			}

			for storeUID := range clients {
				for client := range clients[storeUID] {
					err := client.WriteMessage(websocket.TextMessage, []byte("1"))
					if err != nil {
						fmt.Println("realtime", err)
						client.Close()
						delete(clients[storeUID], client)
					}
				}
				if len(clients[storeUID]) == 0 {
					delete(clients, storeUID)
				}
			}
		}

		// customers
		if len(clients) > 0 {
			customers, _, ok := selectProcessNoLogging("select * from " + customerTable + " where modified_date_time > '" + time.Now().Add(-10*time.Second).UTC().String() + "'")

			if ok && len(customers) > 0 {
				for _, customer := range customers {
					for client := range clients[customer["store_u_id"]] {
						client.WriteJSON(map[string]map[string]string{
							"customer": customer,
						})
					}
				}
			}

			for storeUID := range clients {
				for client := range clients[storeUID] {
					err := client.WriteMessage(websocket.TextMessage, []byte("1"))
					if err != nil {
						fmt.Println("realtime", err)
						client.Close()
						delete(clients[storeUID], client)
					}
				}
				if len(clients[storeUID]) == 0 {
					delete(clients, storeUID)
				}
			}
		}
	}
}
