package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Collection struct {
	Info struct {
		Name string `json:"name"`
	} `json:"info"`
	Item     []Item     `json:"item"`
	Variable []Variable `json:"variable"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Item struct {
	Name    string   `json:"name"`
	Item    []Item   `json:"item,omitempty"`
	Request *Request `json:"request,omitempty"`
}

type Request struct {
	Auth   *Auth  `json:"auth,omitempty"`
	Method string `json:"method"`
	Header []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
	Body *Body `json:"body,omitempty"`
	URL  URL   `json:"url"`
}

type Auth struct {
	Type   string `json:"type"`
	Bearer []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"bearer"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type URL struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Disabled bool   `json:"disabled"`
	} `json:"query"`
}

func extractRequests(items []Item, path string) {
	for _, item := range items {
		currentPath := path + "/" + item.Name
		if item.Request != nil {
			method := item.Request.Method
			url := item.Request.URL.Raw
			hasBody := "no-body"
			if item.Request.Body != nil && item.Request.Body.Raw != "" {
				hasBody = "has-body"
			}
			fmt.Printf("%s | %s | %s | %s\n", method, url, hasBody, currentPath)
		}
		if len(item.Item) > 0 {
			extractRequests(item.Item, currentPath)
		}
	}
}

func main() {
	data, err := os.ReadFile("NEW FULL CPS ACTION.postman_collection.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var col Collection
	if err := json.Unmarshal(data, &col); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== POSTMAN COLLECTION VARIABLES ===")
	for _, v := range col.Variable {
		fmt.Printf("  %s = %s\n", v.Key, v.Value)
	}

	fmt.Println("\n=== ALL API ENDPOINTS ===")
	extractRequests(col.Item, "")

	// Count
	count := 0
	var countItems func([]Item)
	countItems = func(items []Item) {
		for _, item := range items {
			if item.Request != nil {
				count++
			}
			if len(item.Item) > 0 {
				countItems(item.Item)
			}
		}
	}
	countItems(col.Item)
	fmt.Printf("\n=== TOTAL ENDPOINTS: %d ===\n", count)
}
