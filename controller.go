package main	

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *MuxServer) addUser(w http.ResponseWriter, r *http.Request){
	var userData UserParam 
	var user User

	//Decode the input api json
	json.NewDecoder(r.Body).Decode(&userData)

	user.Name = userData.Name
	user.Email = userData.Email
	user.Age = userData.Age

	err := s.db.Create(&user).Error
	if err!=nil{
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			http.Error(w, pgErr.Message, http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type","application/json")

	// Encode the user struct using the WriteHeader
	json.NewEncoder(w).Encode(user)
}

func (s *MuxServer) updateUser(w http.ResponseWriter, r *http.Request){
	var userData UserParam
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get first matching record that matches the userId
	var user User
	s.db.First(&user,userId)  //Assign it to user

	user.Name = userData.Name
	user.Email = userData.Email
	user.Age = userData.Age 

	s.db.Save(&user)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":"User data updated successfully",
	})
	return
}

func (s *MuxServer) deleteUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
	}

	var user User 
	s.db.First(&user, userId)
	s.db.Delete(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User has been deleted successfully",
	})
	return

}

func (s *MuxServer) listUsers(w http.ResponseWriter, r *http.Request){
	var users []User
	s.db.Find(&users)
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
}