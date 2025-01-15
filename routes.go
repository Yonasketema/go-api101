package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", getAllBooks)
	mux.HandleFunc("GET /books/{id}", getBook)

	mux.HandleFunc("POST /books", createBook)

	mux.HandleFunc("PATCH /books/{id}", updateBook)

	mux.HandleFunc("DELETE /books/{id}", deleteBook)

	return mux
}
