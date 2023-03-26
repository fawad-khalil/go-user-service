package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	db "github.com/user-service/services"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	id        string `gorm:"type:number;not null;primaryKey;autoIncrement" json:"id"`
	uid       string `gorm:"type:text;not null" json:"uid"`
	firstName string `gorm:"type:text;not null" json:"firstName"`
	lastName  string `gorm:"type:text;not null" json:"lastName"`
	password  string `gorm:"type:text;not null" json:"password"`
	role      string `gorm:"type:text;not null" json:"role"`
	createdAt time.Time
	updatedAt time.Time
	deletedAt gorm.DeletedAt
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var user User
	result := db.DB.First(&user, mux.Vars(r)["uid"])

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(user)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, "Failed to save.", http.StatusBadRequest)
		return
	}

	resErr := json.NewEncoder(w).Encode(result)

	if resErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	uid := mux.Vars(r)["uid"]

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.DB.Model(&User{}).Where("uid =", uid).Save(&user)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// Return updated user in response
	json.NewEncoder(w).Encode(&user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func login(w http.ResponseWriter, r *http.Request) {

}

func logout(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("Starting Service...")

	dbErr := db.NewSQLiteDB("development.db")

	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })

	// RedisPong, err := client.Ping(client.Context()).Result()

	if dbErr != nil {
		log.Fatal("Failed to connect to DB: ", dbErr)
		return
	}

	fmt.Println("Connected to SQLite DB")
	// fmt.Println("Connected to Redis", RedisPong)

	// err := db.DB.AutoMigrate(&User{})

	if db.DB.Migrator().HasTable(&User{}) {
		if err := db.DB.First(&User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			user := User{firstName: "John", lastName: "Doe", password: "admin", role: "admin"}

			result := db.DB.Create(&user)
			fmt.Println("Admin user Seeding Error: ", result.Error)
			fmt.Println("Admin user seed, total rows affected: ", result.RowsAffected)
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	apiSubRouter := r.PathPrefix("/api").Subrouter()
	usersSubRouter := apiSubRouter.PathPrefix("/users").Subrouter()

	usersSubRouter.HandleFunc("/", getUser).Methods("GET")
	usersSubRouter.HandleFunc("/", addUser).Methods("POST")
	usersSubRouter.HandleFunc("/api/users", updateUser).Methods("PUT")
	usersSubRouter.HandleFunc("/api/users", deleteUser).Methods("DELETE")

	authSubRouter := apiSubRouter.PathPrefix("/auth").Subrouter()

	authSubRouter.HandleFunc("/", login).Methods("POST")
	authSubRouter.HandleFunc("/", logout).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}
