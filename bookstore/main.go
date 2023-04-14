package main

import (
    "encoding/json"
    // "fmt"
    "log"
    "net/http"
    // "strconv"
    "github.com/gorilla/mux"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Book struct {
    gorm.Model
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Cost        float64 `json:"cost"`
}

var db *gorm.DB
var err error

func main() {
    // Set up database connection
    dsn := "user:password@tcp(database:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    // Migrate the schema
    db.AutoMigrate(&Book{})
    // Set up router
    router := mux.NewRouter()
    router.HandleFunc("/books", getBooks).Methods("GET")
    router.HandleFunc("/books/{id}", getBook).Methods("GET")
    router.HandleFunc("/books", addBook).Methods("POST")
    router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
    router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
    router.HandleFunc("/books/search", searchBook).Methods("GET")
    router.HandleFunc("/books/sort", sortBooks).Methods("GET")
    // Start server
    log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    var books []Book
    db.Find(&books)
    json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var book Book
    db.First(&book, params["id"])
    json.NewEncoder(w).Encode(&book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    json.NewDecoder(r.Body).Decode(&book)
    db.Create(&book)
    json.NewEncoder(w).Encode(&book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var book Book
    db.First(&book, params["id"])
    json.NewDecoder(r.Body).Decode(&book)
    db.Save(&book)
    json.NewEncoder(w).Encode(&book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var book Book
    db.Delete(&book, params["id"])
    json.NewEncoder(w).Encode("The book is deleted successfully")
}

func searchBook(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Query().Get("title")
    var book Book
    db.Where("title = ?", title).First(&book)
    json.NewEncoder(w).Encode(&book)
}

func sortBooks(w http.ResponseWriter, r *http.Request) {
    order := r.URL.Query().Get("order")
    var books []Book
    if order == "asc" {
        db.Order("cost asc").Find(&books)
    } else {
        db.Order("cost desc").Find(&books)
    }
    json.NewEncoder(w).Encode(&books)
}