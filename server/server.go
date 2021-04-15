package server

import (
	"encoding/json"
	"html/template"
	"http-server/definitions"
	"http-server/repository"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router

	*datastore.Client
}

func NewServer(c *datastore.Client) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		Client: c,
	}
	s.routes()
	return s
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*"))
}

func (s *Server) routes() {
	s.HandleFunc("/", s.home()).Methods("GET")
	s.HandleFunc("/teams/{teamID}/members", s.listMembers()).Methods("GET")
	s.HandleFunc("/members", s.createMember()).Methods("POST")
	s.HandleFunc("/teams/{teamID}/members/{id}", s.removeMember()).Methods("DELETE")
	s.HandleFunc("members/{id}", s.getMember()).Methods("GET")
	s.HandleFunc("/teams/{teamID}/members/{email}", s.getMemberByEmail()).Methods("GET")
	s.HandleFunc("/teams", s.createTeam()).Methods("POST")
	s.HandleFunc("/teams/{id}", s.getTeam()).Methods("GET")
	s.HandleFunc("/teams", s.getAllTeams()).Methods("GET")
}

func (s *Server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := repository.GetAllTeams(s.Client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err1 := tpl.ExecuteTemplate(w, "index.gohtml", teams)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) createMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i definitions.Member
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key, err := repository.AddMember(s.Client, i)
		if err != nil {
			log.Printf("Failed to create member: %v", err)
			return
		}
		i.ID = key.ID
		log.Printf("Created new member with ID %d\n", key.ID)

		team, err := repository.AddMemberToTeam(s.Client, i, i.TeamID)

		w.Header().Set("Content-Type", "application/")
		if err := json.NewEncoder(w).Encode(team); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listMembers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamIDStr, _ := mux.Vars(r)["teamID"]
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse id: %v", err)
			return
		}
		log.Println(teamID)
		members, err := repository.GetMembers(s.Client, teamID)
		if err != nil {
			log.Printf("Failed to fetch member: %v", err)
			return
		}

		err1 := tpl.ExecuteTemplate(w, "index.gohtml", members)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse id: %v", err)
			return
		}
		teamIDStr, _ := mux.Vars(r)["teamID"]
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse id: %v", err)
			return
		}

		if err := repository.RemoveMember(s.Client, id, teamID); err != nil {
			log.Printf("Failed to delete: %v", err)
		}

		return
	}
}

func (s *Server) getMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Falied to parse id: %v", err)
		}
		member, err := repository.GetMember(s.Client, id)
		if err != nil {
			log.Printf("Falied to fetch member: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(member); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) createTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i definitions.Team
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key, err := repository.AddTeam(s.Client, &i)
		if err != nil {
			log.Printf("Failed to create team: %v", err)
			return
		}
		i.TeamID = key.ID
		log.Printf("Created new team with ID %d\n", key.ID)
		w.Header().Set("Content-Type", "application/")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Falied to parse id: %v", err)
		}
		team, err := repository.GetTeam(s.Client, id)
		if err != nil {
			log.Printf("Falied to fetch team: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(team); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) getMemberByEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailID, _ := mux.Vars(r)["email"]
		teamIDStr, _ := mux.Vars(r)["teamID"]
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse team ID: %v", err)
		}
		member, err := repository.GetMemberByEmail(s.Client, teamID, emailID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(member); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) getAllTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := repository.GetAllTeams(s.Client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err1 := tpl.ExecuteTemplate(w, "index.gohtml", teams)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
