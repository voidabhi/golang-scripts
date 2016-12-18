package main

import (
    "encoding/json"
    "net/http"
    "github.com/graphql-go/graphql"
    "github.com/graphql-go/handler"
)

func customHandler(schema *graphql.Schema) func(http.ResponseWriter, *http.Request) {
    return func(rw http.ResponseWriter, r *http.Request) {
        opts := handler.NewRequestOptions(r)

        rootValue := map[string]interface{}{
            "response": rw,
            "request":  r,
        }

        params := graphql.Params{
            Schema:         *schema,
            RequestString:  opts.Query,
            VariableValues: opts.Variables,
            OperationName:  opts.OperationName,
            RootObject:     rootValue,
        }

        result := graphql.Do(params)

        jsonStr, err := json.Marshal(result)

        if err != nil {
            panic(err)
        }

        rw.Header().Set("Content-Type", "application/json")

        rw.Header().Set("Access-Control-Allow-Credentials", "true")
        rw.Header().Set("Access-Control-Allow-Methods", "POST")
        rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        rw.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

        rw.Write(jsonStr)
    }
}

func main() {
    queryType := graphql.ObjectConfig{
        Name: "Ping",
        Description: "Ping to get pong...",

        Fields: graphql.Fields{
            "ping": &graphql.Field{
                Type: graphql.String,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    return "pong", nil
                },
            },
        },
    }

    schema, err := graphql.NewSchema(graphql.SchemaConfig{
        Query: graphql.NewObject(queryType),
    })

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/graphql", customHandler(&schema))
    http.ListenAndServe(":8888", nil)
}
