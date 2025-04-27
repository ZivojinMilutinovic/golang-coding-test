package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) Set(key string, value interface{}, ttl int) error {
	body := map[string]interface{}{
		"value": value, "ttl": ttl,
	}
	return c.post("/set/"+key, body)
}

func (c *Client) Get(key string) (interface{}, error) {
	resp, err := http.Get(c.BaseURL + "/get/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *Client) Update(key string, value interface{}) error {
	body := map[string]interface{}{
		"value": value,
	}
	return c.post("/update"+key, body)
}

func (c *Client) Remove(key string) error {
	req, _ := http.NewRequest(http.MethodDelete, c.BaseURL+"/remove/"+key, nil)
	_, err := http.DefaultClient.Do(req)
	return err
}

func (c *Client) Push(key, item string) error {
	body := map[string]interface{}{
		"value": item,
	}
	return c.post("/push/"+key, body)
}

func (c *Client) Pop(key string) (interface{}, error) {
	resp, err := http.Post(c.BaseURL+"/pop/"+key, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *Client) post(path string, body map[string]interface{}) error {
	b, _ := json.Marshal(body)
	_, err := http.Post(c.BaseURL+path, "application/json", bytes.NewBuffer(b))
	return err
}

// Testing function for client

func TestClient() {
	// Initialize the client with the base URL
	fmt.Println("Calling test client function")
	baseURL := "http://localhost:8080"
	c := NewClient(baseURL)

	// 1. Set a string value
	err := c.Set("stringKey", "Hello, World!", 0)
	if err != nil {
		log.Fatalf("Error setting string value: %v", err)
	}
	fmt.Println("String value 'Hello, World!' set for 'stringKey'.")

	// 2. Set a list (initialize with an empty list)
	err = c.Set("myList", []string{}, 0)
	if err != nil {
		log.Fatalf("Error setting list: %v", err)
	}
	fmt.Println("List 'myList' initialized.")

	// 3. Push multiple values to the list
	valuesToPush := []string{"value1", "value2", "value3", "value4", "value5"}
	for _, value := range valuesToPush {
		err = c.Push("myList", value)
		if err != nil {
			log.Fatalf("Error pushing value '%s' to list: %v", value, err)
		}
		fmt.Printf("Pushed '%s' to 'myList'.\n", value)
	}

	// 4. Get and print the current list (optional)
	currentList, err := c.Get("myList")
	if err != nil {
		log.Fatalf("Error getting list: %v", err)
	}
	fmt.Printf("Current 'myList' values: %v\n", currentList)

	// 5. Pop values from the list (popping all values one by one)
	for i := 0; i < len(valuesToPush); i++ {
		poppedValue, err := c.Pop("myList")
		if err != nil {
			log.Fatalf("Error popping value from list: %v", err)
		}
		fmt.Printf("Popped value: %s\n", poppedValue)
	}

	// 6. Delete a value (remove the string key)
	err = c.Remove("stringKey")
	if err != nil {
		log.Fatalf("Error deleting value: %v", err)
	}
	fmt.Println("Deleted 'stringKey'.")
}
