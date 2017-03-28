package routing

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type organization struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type organizations []organization

var dbRef *sql.DB

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/organizations", http.StatusSeeOther)
}

func organizationIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := dbRef.Query("SELECT * FROM organizations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var orgs organizations
	for rows.Next() {
		var org organization
		err := rows.Scan(&org.Id, &org.Name, &org.Created_at, &org.Updated_at)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		orgs = append(orgs, org)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func organizationInstance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var org organization
	row := dbRef.QueryRow("SELECT * FROM organizations WHERE id=$1;", ps.ByName("id"))

	err := row.Scan(&org.Id, &org.Name, &org.Created_at, &org.Updated_at)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(&org)
}

func organizationCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var org organization
	var lastID int
	json.NewDecoder(r.Body).Decode(&org)

	err := dbRef.QueryRow("INSERT INTO organizations (name) VALUES ($1) RETURNING id;", org.Name).Scan(&lastID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Status", "201")
	w.Write([]byte("Create Organization: " + strconv.Itoa(lastID)))
}

func organizationUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var reqOrg organization
	var lastID int
	err := json.NewDecoder(r.Body).Decode(&reqOrg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = dbRef.QueryRow("UPDATE organizations SET name=$1 WHERE id=$2 RETURNING id;", &reqOrg.Name, ps.ByName("id")).Scan(&lastID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Status", "201")
	w.Write([]byte("Updated Organization: " + strconv.Itoa(lastID)))
}

func organizationDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var lastID int
	err := dbRef.QueryRow("DELETE FROM organizations WHERE id=$1 RETURNING id;", ps.ByName("id")).Scan(&lastID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Status", "204")
	w.Write([]byte("Deleted Organization: " + strconv.Itoa(lastID)))
}

func Router(db *sql.DB) *httprouter.Router {
	dbRef = db

	mux := httprouter.New()
	mux.GET("/", index)
	mux.GET("/organizations", organizationIndex)
	mux.GET("/organizations/:id", organizationInstance)
	mux.POST("/organizations", organizationCreate)
	mux.PUT("/organizations/:id", organizationUpdate)
	mux.DELETE("/organizations/:id", organizationDelete)

	return mux
}
