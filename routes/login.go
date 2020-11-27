package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/powerqueue/fitque-users-api/models"
)

//RetrieveLoginHandler - handler definition
func (ls *LoginServer) RetrieveLoginHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	serviceRequest := &models.LoginDefinition{}
	err := json.NewDecoder(r.Body).Decode(&serviceRequest)
	if err != nil {
		// return nil, rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}

	loginDef, err := ls.loginService.GetLogin(serviceRequest)
	if err != nil {
		// err = rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", loginDef)
}

//CreateLogin - handler method definition
func (ls *LoginServer) CreateLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside CreateLogin Route Handler")
	serviceRequest := &models.LoginDefinition{}
	err := json.NewDecoder(r.Body).Decode(&serviceRequest)
	if err != nil {
		// return nil, rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}

	loginDef, err := ls.loginService.CreateLogin(serviceRequest)
	if err != nil {
		// err = rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", loginDef)
}

//TermLogin - handler method definition
func (ls *LoginServer) TermLogin(w http.ResponseWriter, r *http.Request) {
	serviceRequest := &models.LoginDefinition{}
	err := json.NewDecoder(r.Body).Decode(&serviceRequest)
	if err != nil {
		// return nil, rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}

	loginDef, err := ls.loginService.TermLogin(serviceRequest)
	if err != nil {
		// err = rest.ErrorCode{Code: 500, Err: err}
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", loginDef)
}
