package main

import (
	"log"
	"net/http"

	"football-live-score/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// 1. Buat server GraphQL dengan resolver yang sudah kita buat
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// 2. Pasang GraphQL Playground di root url "/"
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// 3. Pasang endpoint GraphQL API di "/query"
	http.Handle("/query", srv)

	log.Printf("GraphQL Playground siap diakses di http://localhost:%s/", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}