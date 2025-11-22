package prservice

import (
	"github.com/Egorrrad/avitotechBackendPR/internal/handlers"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	h := handlers.NewHTTPHandler(nil)
	r := mux.NewRouter()

	// pullRequest
	r.HandleFunc("/pullRequest/create", h.PostPullRequestCreate).Methods("POST")
	r.HandleFunc("/pullRequest/merge", h.PostPullRequestMerge).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", h.PostPullRequestReassign).Methods("POST")

	// team
	r.HandleFunc("/team/add", h.PostTeamAdd).Methods("POST")
	r.HandleFunc("/team/get", h.GetTeamGet).Methods("GET")

	// users
	r.HandleFunc("/users/getReview", h.GetUsersGetReview).Methods("GET")
	r.HandleFunc("/users/setIsActive", h.PostUsersSetIsActive).Methods("POST")

	return r
}
