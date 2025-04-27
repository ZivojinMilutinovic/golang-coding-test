# golang-coding-test
Golang coding test
DOCKER BUILD
To build a docker image run: docker build -t go-gin-app .  
To run the program in Docker run docker run -p 8080:8080 go-gin-app

Client

The Client struct contains the base URL for the server and methods for interacting with the API.

Fields:

    BaseURL (string): The base URL of the server.

Client Methods
NewClient(baseURL string) *Client

Creates a new client with the given base URL.

Parameters:

    baseURL (string): The base URL of the server.

Returns:

    *Client: A pointer to a new Client instance.

Set(key string, value interface{}, ttl int) error

Sets a value for the given key with an optional TTL (time-to-live).

Parameters:

    key (string): The key to associate with the value.

    value (interface{}): The value to associate with the key.

    ttl (int): The TTL (in seconds) for the value.

Returns:

    error: An error if the request fails.

Get(key string) (interface{}, error)

Gets the value for the given key.

Parameters:

    key (string): The key to retrieve the value for.

Returns:

    (interface{}, error): The value associated with the key and an error if the request fails.

Update(key string, value interface{}) error

Updates the value associated with the given key.

Parameters:

    key (string): The key to update the value for.

    value (interface{}): The new value to associate with the key.

Returns:

    error: An error if the request fails.

Remove(key string) error

Removes the value associated with the given key.

Parameters:

    key (string): The key to remove.

Returns:

    error: An error if the request fails.

Push(key, item string) error

Pushes a new item to the list stored at the given key.

Parameters:

    key (string): The key for the list.

    item (string): The item to push to the list.

Returns:

    error: An error if the request fails.

Pop(key string) (interface{}, error)

Pops the last item from the list stored at the given key.

Parameters:

    key (string): The key for the list.

Returns:

    (interface{}, error): The popped item and an error if the request fails.

post(path string, body map[string]interface{}) error

A helper method that performs a POST request to the given path with the provided body.

Parameters:

    path (string): The API endpoint path to post to.

    body (map[string]interface{}): The body to send with the POST request.

Returns:

    error: An error if the request fails.

TestClient Function
TestClient()

A testing function that demonstrates usage of the Client methods. It performs a sequence of operations like setting values, pushing and popping items from a list, and deleting keys.

Steps:

    Sets a string value "Hello, World!" with the key "stringKey".

    Sets a list ([]string{}) with the key "myList".

    Pushes multiple values to "myList".

    Retrieves and prints the current values in "myList".

    Pops all values from "myList".

    Deletes the key "stringKey".

    Confirms the deletion of "stringKey".

REST API command list:
# Set a value for a key (Set a string value)
curl -X POST http://localhost:8080/set/stringKey -H "Content-Type: application/json" -d "{\"value\": \"Hello, World!\", \"ttl\": 60}"

#  Get a value by key (Key exists)
curl http://localhost:8080/get/stringKey

#  Get a value by key (Key does not exist)
curl http://localhost:8080/get/nonExistentKey

#  Update a value for a key (Key exists)
curl -X POST http://localhost:8080/update/stringKey -H "Content-Type: application/json" -d "{\"value\": \"Updated Value\"}"

#  Update a value for a key (Key does not exist)
curl -X POST http://localhost:8080/update/nonExistentKey -H "Content-Type: application/json" -d "{\"value\": \"New Value\"}"

#  Remove a key (Key exists)
curl -X DELETE http://localhost:8080/remove/stringKey

#  Push a value to a list (Key exists and is a list)
curl -X POST http://localhost:8080/push/myList -H "Content-Type: application/json" -d "{\"value\": \"Item1\"}"

#  Push a value to a list (Key does not exist, will create a new list)
curl -X POST http://localhost:8080/push/myNewList -H "Content-Type: application/json" -d "{\"value\": \"Item2\"}"

#  Pop a value from a list (List exists and has values)
curl -X POST http://localhost:8080/pop/myList

#  Pop a value from a list (List exists but is empty)
curl -X POST http://localhost:8080/pop/myEmptyList

#  Pop a value from a list (List does not exist)
curl -X POST http://localhost:8080/pop/nonExistentList

#  Set a value for a key (Set a list)
curl -X POST http://localhost:8080/set/myList -H "Content-Type: application/json" -d "{\"value\": [\"Item1\", \"Item2\", \"Item3\"], \"ttl\": 60}"

#  Get a list value (Key exists and is a list)
curl http://localhost:8080/get/myList

#  Get a value that is not set (Key does not exist)
curl http://localhost:8080/get/nonExistentKey

#  Set a value with no TTL (without expiration)
curl -X POST http://localhost:8080/set/withoutTTL -H "Content-Type: application/json" -d "{\"value\": \"No TTL value\", \"ttl\": 0}"

#  Update an existing list
curl -X POST http://localhost:8080/update/myList -H "Content-Type: application/json" -d "{\"value\": [\"Item4\", \"Item5\"]}"

# check the updated list
curl http://localhost:8080/get/myList

#  Remove a list
curl -X DELETE http://localhost:8080/remove/myList

# check the list after deleting(should not be present)
curl http://localhost:8080/get/myList

