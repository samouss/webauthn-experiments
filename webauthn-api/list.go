package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// List returns a HTTP function that reply with a static list of places.
func List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type place struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			ImageURL string `json:"imageURL"`
		}

		type response struct {
			Places []place `json:"places"`
		}

		data := response{
			Places: []place{
				{
					ID:       "5477500",
					Name:     "Au bon endroit",
					ImageURL: "https://images.unsplash.com/photo-1414235077428-338989a2e8c0?auto=format&fit=crop&w=1950&q=80",
				},
				{
					ID:       "4397400",
					Name:     "Le citizen",
					ImageURL: "https://images.unsplash.com/photo-1514933651103-005eec06c04b?auto=format&fit=crop&w=1567&q=80",
				},
				{
					ID:       "5588602",
					Name:     `L'enoteca`,
					ImageURL: "https://images.unsplash.com/photo-1552566626-52f8b828add9?auto=format&fit=crop&w=1950&q=80",
				},
				{
					ID:       "5578851",
					Name:     "Ogata",
					ImageURL: "https://images.unsplash.com/photo-1549488344-1f9b8d2bd1f3?auto=format&fit=crop&w=1950&q=80",
				},
				{
					ID:       "4863102",
					Name:     `Limmat`,
					ImageURL: "https://images.unsplash.com/photo-1522336572468-97b06e8ef143?auto=format&fit=crop&w=1955&q=80",
				},
				{
					ID:       "5578849",
					Name:     "Maafim",
					ImageURL: "https://images.unsplash.com/photo-1555396273-367ea4eb4db5?auto=format&fit=crop&w=1567&q=80",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
			Write500JSONError(w)
			return
		}
	}
}
