package handlers

import "net/http"

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// (POST /pullRequest/create)
func (h *HTTPHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {

}

// Пометить PR как MERGED (идемпотентная операция)
// (POST /pullRequest/merge)
func (h *HTTPHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {

}

// Переназначить конкретного ревьювера на другого из его команды
// (POST /pullRequest/reassign)
func (h *HTTPHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {

}
