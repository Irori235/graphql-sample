package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

// User
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Sample data
var data map[string]User

// Define GraphQL types
var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        userType,
				Description: "Get user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := params.Args["id"].(string)
					if isOK {
						return data[idQuery], nil
					}
					return nil, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, variables map[string]interface{}, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

// Set middleware to allow CORS
func allowCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		fn(w, r)
	}
}

// Handle GraphQL requests
func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	// Decode the request body
	var rBody struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &rBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute the query
	result := executeQuery(rBody.Query, rBody.Variables, schema)
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.Handle("/graphql", allowCORS(handleGraphQL))

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Curl : curl -X POST -H 'Content-Type: application/json' -d '{\"query\": \"query getUser($id: String!){ user(id: $id) { id name } }\", \"variables\": {\"id\": \"1\"}}' http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	data = map[string]User{
		"1": {
			ID:   "1",
			Name: "Alice",
		},
		"2": {
			ID:   "2",
			Name: "Bob",
		},
	}
}
