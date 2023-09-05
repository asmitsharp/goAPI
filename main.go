package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

var (
	users      = make(map[int]User)
	usersMutex sync.RWMutex
	nextUserID = 1
)

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	userSlice := make([]User, 0, len(users))
	for _, user := range users {
		userSlice = append(userSlice, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userSlice)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	newUser.ID = nextUserID
	users[nextUserID] = newUser
	nextUserID++

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, found := users[userID]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	_, found := users[userID]
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updatedUser.ID = userID
	users[userID] = updatedUser

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func main() {
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/getAllUser", getAllUsers)
	http.HandleFunc("/getUser", getUser)
	http.HandleFunc("/updateUser", updateUser)

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}
