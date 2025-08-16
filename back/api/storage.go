package api

import (
	"digits_say/storage"
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *Server) GetConscienceText(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	iid, err := strconv.Atoi(id)
	if id == "" || err != nil {
		http.Error(w, "Empty or not correct id", http.StatusBadRequest)
		return
	}

	conscience, exist, err := s.Storage.GetConscienceText(iid)
	if err != nil {
		http.Error(w, "Error getting conscience: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Not Found"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(conscience.Message))
		return
	}
}

func (s *Server) GetCommonDayText(w http.ResponseWriter, r *http.Request) {
	common, exist, err := s.Storage.GetCommonDayText()
	if err != nil {
		http.Error(w, "Error getting common day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Not Found"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(common.Message))
		return
	}
}

func (s *Server) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Empty id", http.StatusBadRequest)
		return
	}
	user, exist, err := s.Storage.GetUserByID(id)
	if err != nil {
		http.Error(w, "Error getting user by Telegram ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Not Found"))
		return
	} else {
		resp, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
			return
		}
	}
}

func (s *Server) GetListOfSubscribers(w http.ResponseWriter, r *http.Request) {
	users, err := s.Storage.GetListOfSubscribers()
	if err != nil {
		http.Error(w, "Error of users listing: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}
}

func (s *Server) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	user := &storage.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "can`t parce json: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err := s.Storage.RegisterNewUser(*user)
	if err != nil {
		http.Error(w, "Error registering a new user: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Ok"))
		return
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &storage.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "can`t parce json: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err := s.Storage.UpdateUser(*user)
	if err != nil {
		http.Error(w, "Error updating a user: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Ok"))
		return
	}
}
