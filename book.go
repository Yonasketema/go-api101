package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type bookModel struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	query := `
    SELECT id, title, author, published_year
    FROM books
    LIMIT 12`

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)

	}

	books := []*bookModel{}

	for rows.Next() {
		var book bookModel

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.PublishedYear,
		)

		if err != nil {

			log.Fatal(err)
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(books)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func getBook(w http.ResponseWriter, r *http.Request) {

	var book bookModel

	bookId := r.PathValue("id")

	if _, err := uuid.Parse(bookId); err != nil {
		fmt.Fprintf(w, "Invalid Book ID '%s'. Please provide a valid UUID.", bookId)
		return
	}

	query := `
		SELECT id, title, author, published_year
		FROM books
		WHERE id = $1`

	err := db.QueryRow(query, bookId).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			w.Write([]byte("Record Not Found !!"))
			return
		default:
			log.Fatal(err)
		}
	}

	js, err := json.Marshal(book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book bookModel

	query := `
		INSERT INTO books (title, author, published_year)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow(query, book.Title, book.Author, book.PublishedYear).Scan(&book.ID)
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(book)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(js)

}

func updateBook(w http.ResponseWriter, r *http.Request) {

	var book bookModel

	bookId := r.PathValue("id")

	if _, err := uuid.Parse(bookId); err != nil {
		fmt.Fprintf(w, "Invalid Book ID '%s'. Please provide a valid UUID.", bookId)
		return
	}

	getQuery := `
		SELECT id ,title, author, published_year
		FROM books
		WHERE id = $1`

	err := db.QueryRow(getQuery, bookId).
		Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			w.Write([]byte("Record Not Found !!"))
			return
		default:
			log.Fatal(err)
		}
	}

	var inputBook struct {
		Title         *string `json:"title"`
		Author        *string `json:"author"`
		PublishedYear *int    `json:"published_year"`
	}

	err = json.NewDecoder(r.Body).Decode(&inputBook)
	if err != nil {
		log.Fatal(err)
	}

	if inputBook.Title != nil {
		book.Title = *inputBook.Title
	}

	if inputBook.Author != nil {
		book.Author = *inputBook.Author
	}

	if inputBook.PublishedYear != nil {
		book.PublishedYear = *inputBook.PublishedYear
	}

	updateQuery := `
		UPDATE books
		SET title = $1, author = $2, published_year = $3
		WHERE id = $4`

	_, err = db.Exec(updateQuery, book.Title, book.Author, book.PublishedYear, book.ID)
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(book)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	bookId := r.PathValue("id")

	if _, err := uuid.Parse(bookId); err != nil {
		fmt.Fprintf(w, "Invalid Book ID '%s'. Please provide a valid UUID.", bookId)
		return
	}

	query := `
		DELETE FROM books
		WHERE id = $1`

	result, err := db.Exec(query, bookId)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		w.Write([]byte("Record Not Found !"))
		return
	}

	w.Write([]byte("book successfully deleted"))
}
