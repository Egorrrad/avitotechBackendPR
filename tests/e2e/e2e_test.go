// nolint
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

// --- Вспомогательные DTO ---

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type PullRequestCreateReq struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestMergeReq struct {
	PullRequestID string `json:"pull_request_id"`
}

type PullRequestReassignReq struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type PullRequestResponse struct {
	PR struct {
		ID                string   `json:"pull_request_id"`
		Name              string   `json:"pull_request_name"`
		AuthorID          string   `json:"author_id"`
		Status            string   `json:"status"`
		AssignedReviewers []string `json:"assigned_reviewers"`
		CreatedAt         *string  `json:"createdAt"`
		MergedAt          *string  `json:"mergedAt"`
	} `json:"pr"`
	ReplacedBy string `json:"replaced_by,omitempty"`
}

type UserReviewsResponse struct {
	UserID       string `json:"user_id"`
	PullRequests []struct {
		ID       string `json:"pull_request_id"`
		Name     string `json:"pull_request_name"`
		AuthorID string `json:"author_id"`
		Status   string `json:"status"`
	} `json:"pull_requests"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func randomString(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), rand.Intn(100000))
}

func sendRequest(t *testing.T, method, endpoint string, payload interface{}) (int, []byte) {
	var bodyReader io.Reader
	if payload != nil {
		jsonBody, err := json.Marshal(payload)
		require.NoError(t, err)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bodyReader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, respBody
}

// --- ТЕСТЫ ---

func TestE2E_AutoAssign(t *testing.T) {
	// Сценарий: В команде 3 человека. Автор создает PR. Остальные двое должны стать ревьюверами.
	teamName := randomString("team")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	rev2 := randomString("u_rev2")
	prID := randomString("pr")

	// 1. Создаем команду
	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: rev1, Username: "Rev1", IsActive: true},
			{UserID: rev2, Username: "Rev2", IsActive: true},
		},
	}
	code, _ := sendRequest(t, "POST", "/team/add", teamReq)
	require.Equal(t, http.StatusCreated, code)

	// 2. Создаем PR
	prReq := PullRequestCreateReq{
		PullRequestID:   prID,
		PullRequestName: "Feat 1",
		AuthorID:        author,
	}
	code, body := sendRequest(t, "POST", "/pullRequest/create", prReq)
	require.Equal(t, http.StatusCreated, code)

	var resp PullRequestResponse
	err := json.Unmarshal(body, &resp)
	require.NoError(t, err)

	// 3. Проверки
	assert.Equal(t, "OPEN", resp.PR.Status)
	assert.Len(t, resp.PR.AssignedReviewers, 2, "Should have 2 reviewers")
	assert.NotContains(t, resp.PR.AssignedReviewers, author, "Author cannot be reviewer")

	// Проверяем, что оба ревьювера из команды
	for _, reviewer := range resp.PR.AssignedReviewers {
		assert.Contains(t, []string{rev1, rev2}, reviewer, "Reviewer must be from team")
	}

	// Проверяем, что назначены оба доступных
	assert.Contains(t, resp.PR.AssignedReviewers, rev1)
	assert.Contains(t, resp.PR.AssignedReviewers, rev2)

	// Проверяем наличие createdAt
	assert.NotNil(t, resp.PR.CreatedAt, "CreatedAt should be set")
	assert.Nil(t, resp.PR.MergedAt, "MergedAt should be nil for OPEN PR")
}

func TestE2E_NotEnoughCandidates(t *testing.T) {
	// Сценарий: В команде только Автор и 1 коллега. Должен назначиться только 1 ревьювер.
	teamName := randomString("team_small")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	prID := randomString("pr")

	// 1. Создаем команду (2 человека)
	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "Author2", IsActive: true},
			{UserID: rev1, Username: "Rev1", IsActive: true},
		},
	}
	sendRequest(t, "POST", "/team/add", teamReq)

	// 2. Создаем PR
	prReq := PullRequestCreateReq{
		PullRequestID:   prID,
		PullRequestName: "Small Team PR",
		AuthorID:        author,
	}
	code, body := sendRequest(t, "POST", "/pullRequest/create", prReq)
	require.Equal(t, http.StatusCreated, code)

	var resp PullRequestResponse
	json.Unmarshal(body, &resp)

	assert.Len(t, resp.PR.AssignedReviewers, 1, "Should have only 1 reviewer")
	assert.Equal(t, rev1, resp.PR.AssignedReviewers[0])
}

func TestE2E_InactiveUsersIgnored(t *testing.T) {
	// Сценарий: В команде Автор, Активный ревьювер и Неактивный. Неактивный не должен попасть в ревью.
	teamName := randomString("team_inactive")
	author := randomString("u_auth")
	activeRev := randomString("u_active")
	inactiveRev := randomString("u_inactive")
	prID := randomString("pr")

	// 1. Создаем команду
	teamReq := TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: activeRev, Username: "Active", IsActive: true},
			{UserID: inactiveRev, Username: "Inactive", IsActive: false},
		},
	}
	sendRequest(t, "POST", "/team/add", teamReq)

	// 2. Создаем PR
	prReq := PullRequestCreateReq{
		PullRequestID:   prID,
		PullRequestName: "Inactive Check",
		AuthorID:        author,
	}
	code, body := sendRequest(t, "POST", "/pullRequest/create", prReq)
	require.Equal(t, http.StatusCreated, code)

	var resp PullRequestResponse
	json.Unmarshal(body, &resp)

	assert.Len(t, resp.PR.AssignedReviewers, 1)
	assert.Equal(t, activeRev, resp.PR.AssignedReviewers[0])
	assert.NotContains(t, resp.PR.AssignedReviewers, inactiveRev)
}

func TestE2E_MergeIdempotency(t *testing.T) {
	// Сценарий: Мерджим PR два раза. Второй раз не должен быть ошибкой.
	teamName := randomString("team_merge")
	author := randomString("u_auth")
	prID := randomString("pr_merge")

	// 1. Setup
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members:  []TeamMember{{UserID: author, Username: "A", IsActive: true}},
	})
	sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "Merge me", AuthorID: author,
	})

	// 2. Первый Merge
	reqMerge := PullRequestMergeReq{PullRequestID: prID}
	code1, body1 := sendRequest(t, "POST", "/pullRequest/merge", reqMerge)
	require.Equal(t, http.StatusOK, code1)

	var resp1 PullRequestResponse
	json.Unmarshal(body1, &resp1)
	assert.Equal(t, "MERGED", resp1.PR.Status)
	assert.NotNil(t, resp1.PR.MergedAt, "MergedAt should be set after merge")

	// Парсим время из первого ответа
	time1, err := time.Parse(time.RFC3339Nano, *resp1.PR.MergedAt)
	require.NoError(t, err)

	// 3. Второй Merge (Идемпотентность)
	time.Sleep(1000 * time.Millisecond) // Небольшая задержка
	code2, body2 := sendRequest(t, "POST", "/pullRequest/merge", reqMerge)
	require.Equal(t, http.StatusOK, code2, "Second merge should return 200 OK")

	var resp2 PullRequestResponse
	json.Unmarshal(body2, &resp2)
	assert.Equal(t, "MERGED", resp2.PR.Status)
	assert.NotNil(t, resp2.PR.MergedAt)

	// Парсим время из второго ответа
	time2, err := time.Parse(time.RFC3339Nano, *resp2.PR.MergedAt)
	require.NoError(t, err)

	// Время мерджа должно остаться прежним (идемпотентность)
	assert.WithinDuration(t, time1, time2, time.Microsecond, "MergedAt timestamp should not change")
}

func TestE2E_ReassignReviewer(t *testing.T) {
	// Сценарий: Команда из 4 человек (Автор + 3 коллеги).
	// Назначено 2 ревьювера. Меняем одного из них. Должен выбраться третий (свободный).
	teamName := randomString("team_reassign")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	rev2 := randomString("u_rev2")
	rev3 := randomString("u_rev3")
	prID := randomString("pr_reassign")

	// 1. Создаем команду
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "A", IsActive: true},
			{UserID: rev1, Username: "R1", IsActive: true},
			{UserID: rev2, Username: "R2", IsActive: true},
			{UserID: rev3, Username: "R3", IsActive: true},
		},
	})

	// 2. Создаем PR
	code, body := sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "Reassign", AuthorID: author,
	})
	require.Equal(t, http.StatusCreated, code)

	var createResp PullRequestResponse
	json.Unmarshal(body, &createResp)

	currentReviewers := createResp.PR.AssignedReviewers
	require.Len(t, currentReviewers, 2, "Should have 2 reviewers initially")

	// Определяем, кого будем менять (первого из списка)
	toReplace := currentReviewers[0]
	toStay := currentReviewers[1]

	// 3. Делаем Reassign
	reassignReq := PullRequestReassignReq{
		PullRequestID: prID,
		OldUserID:     toReplace,
	}
	codeRe, bodyRe := sendRequest(t, "POST", "/pullRequest/reassign", reassignReq)
	require.Equal(t, http.StatusOK, codeRe)

	var reassignResp PullRequestResponse
	json.Unmarshal(bodyRe, &reassignResp)

	// 4. Проверки
	newReviewers := reassignResp.PR.AssignedReviewers
	assert.Len(t, newReviewers, 2, "Should still have 2 reviewers")
	assert.Contains(t, newReviewers, toStay, "Existing reviewer should stay")
	assert.NotContains(t, newReviewers, toReplace, "Old reviewer should be removed")

	// Новый ревьювер должен быть из команды и не быть автором
	newReviewer := ""
	for _, r := range newReviewers {
		if r != toStay {
			newReviewer = r
			break
		}
	}
	assert.NotEmpty(t, newReviewer, "Should have a new reviewer")
	assert.NotEqual(t, author, newReviewer, "New reviewer cannot be author")
	assert.Contains(t, []string{rev1, rev2, rev3}, newReviewer, "New reviewer must be from team")

	// Проверяем поле replaced_by
	assert.Equal(t, newReviewer, reassignResp.ReplacedBy, "ReplacedBy should contain new reviewer ID")
}

func TestE2E_CannotReassignMerged(t *testing.T) {
	// Сценарий: Нельзя менять ревьювера, если PR уже смержен.
	teamName := randomString("team_fail_reassign")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	rev2 := randomString("u_rev2")
	prID := randomString("pr_fail_reassign")

	// 1. Setup
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "A", IsActive: true},
			{UserID: rev1, Username: "R1", IsActive: true},
			{UserID: rev2, Username: "R2", IsActive: true},
		},
	})
	sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "PR", AuthorID: author,
	})

	// 2. Merge
	sendRequest(t, "POST", "/pullRequest/merge", PullRequestMergeReq{PullRequestID: prID})

	// 3. Try Reassign
	code, body := sendRequest(t, "POST", "/pullRequest/reassign", PullRequestReassignReq{
		PullRequestID: prID, OldUserID: rev1,
	})

	assert.Equal(t, http.StatusConflict, code)

	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "PR_MERGED", errResp.Error.Code)
}

func TestE2E_GetUsersReview(t *testing.T) {
	// Сценарий: Проверить, что PR отображается в списке ревью пользователя.
	teamName := randomString("team_get")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	prID := randomString("pr_get")

	// 1. Setup
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "A", IsActive: true},
			{UserID: rev1, Username: "R1", IsActive: true},
		},
	})
	sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "PR", AuthorID: author,
	})

	// 2. Get Reviews
	code, body := sendRequest(t, "GET", fmt.Sprintf("/users/getReview?user_id=%s", rev1), nil)
	require.Equal(t, http.StatusOK, code)

	var resp UserReviewsResponse
	json.Unmarshal(body, &resp)

	assert.Equal(t, rev1, resp.UserID)
	require.Len(t, resp.PullRequests, 1)
	assert.Equal(t, prID, resp.PullRequests[0].ID)
}

// --- ДОПОЛНИТЕЛЬНЫЕ ТЕСТЫ ---

func TestE2E_TeamAlreadyExists(t *testing.T) {
	// Тест на создание дубликата команды
	teamName := randomString("team_dup")
	author := randomString("u_auth")

	teamReq := TeamRequest{
		TeamName: teamName,
		Members:  []TeamMember{{UserID: author, Username: "A", IsActive: true}},
	}

	// Первое создание - успешно
	code, _ := sendRequest(t, "POST", "/team/add", teamReq)
	require.Equal(t, http.StatusCreated, code)

	// Второе создание - ошибка
	code, body := sendRequest(t, "POST", "/team/add", teamReq)
	assert.Equal(t, http.StatusBadRequest, code)

	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "TEAM_EXISTS", errResp.Error.Code)
}

func TestE2E_PRAlreadyExists(t *testing.T) {
	// Тест на создание дубликата PR
	teamName := randomString("team_pr_dup")
	author := randomString("u_auth")
	prID := randomString("pr_dup")

	// Setup команды
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members:  []TeamMember{{UserID: author, Username: "A", IsActive: true}},
	})

	prReq := PullRequestCreateReq{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:        author,
	}

	// Первое создание - успешно
	code, _ := sendRequest(t, "POST", "/pullRequest/create", prReq)
	require.Equal(t, http.StatusCreated, code)

	// Второе создание - ошибка
	code, body := sendRequest(t, "POST", "/pullRequest/create", prReq)
	assert.Equal(t, http.StatusConflict, code)

	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "PR_EXISTS", errResp.Error.Code)
}

func TestE2E_ReassignNotAssigned(t *testing.T) {
	// Тест на попытку переназначить пользователя, который не назначен
	teamName := randomString("team_not_assigned")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	rev2 := randomString("u_rev2")
	rev3 := randomString("u_rev3")
	notAssigned := randomString("u_not_assigned")
	prID := randomString("pr_not_assigned")

	// Setup команды
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "A", IsActive: true},
			{UserID: rev1, Username: "R1", IsActive: true},
			{UserID: rev2, Username: "R2", IsActive: true},
			{UserID: rev3, Username: "R3", IsActive: true},
		},
	})

	// Создаем отдельную команду для notAssigned
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: randomString("other_team"),
		Members:  []TeamMember{{UserID: notAssigned, Username: "N", IsActive: true}},
	})

	// Создаем PR
	sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "PR", AuthorID: author,
	})

	// Пытаемся переназначить пользователя из другой команды
	code, body := sendRequest(t, "POST", "/pullRequest/reassign", PullRequestReassignReq{
		PullRequestID: prID,
		OldUserID:     notAssigned,
	})

	assert.Equal(t, http.StatusConflict, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_ASSIGNED", errResp.Error.Code)
}

func TestE2E_ReassignNoCandidates(t *testing.T) {
	// Тест на попытку переназначить, когда нет доступных кандидатов
	teamName := randomString("team_no_cand")
	author := randomString("u_auth")
	rev1 := randomString("u_rev1")
	prID := randomString("pr_no_cand")

	// Setup команды только с автором и одним ревьювером
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: author, Username: "A", IsActive: true},
			{UserID: rev1, Username: "R1", IsActive: true},
		},
	})

	// Создаем PR (назначится только rev1)
	sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID: prID, PullRequestName: "PR", AuthorID: author,
	})

	// Пытаемся переназначить rev1, но других кандидатов нет
	code, body := sendRequest(t, "POST", "/pullRequest/reassign", PullRequestReassignReq{
		PullRequestID: prID,
		OldUserID:     rev1,
	})

	assert.Equal(t, http.StatusConflict, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NO_CANDIDATE", errResp.Error.Code)
}

func TestE2E_GetTeamNotFound(t *testing.T) {
	// Тест на получение несуществующей команды
	nonExistentTeam := randomString("team_not_exists")

	code, body := sendRequest(t, "GET", fmt.Sprintf("/team/get?team_name=%s", nonExistentTeam), nil)
	assert.Equal(t, http.StatusNotFound, code)

	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestE2E_SetIsActiveUserNotFound(t *testing.T) {
	// Тест на установку isActive для несуществующего пользователя
	nonExistentUser := randomString("u_not_exists")

	code, body := sendRequest(t, "POST", "/users/setIsActive", map[string]interface{}{
		"user_id":   nonExistentUser,
		"is_active": false,
	})

	assert.Equal(t, http.StatusNotFound, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestE2E_MergeNotFound(t *testing.T) {
	// Тест на мердж несуществующего PR
	nonExistentPR := randomString("pr_not_exists")

	code, body := sendRequest(t, "POST", "/pullRequest/merge", PullRequestMergeReq{
		PullRequestID: nonExistentPR,
	})

	assert.Equal(t, http.StatusNotFound, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestE2E_CreatePRAuthorNotFound(t *testing.T) {
	// Тест на создание PR с несуществующим автором
	nonExistentAuthor := randomString("u_not_exists")
	prID := randomString("pr_no_author")

	code, body := sendRequest(t, "POST", "/pullRequest/create", PullRequestCreateReq{
		PullRequestID:   prID,
		PullRequestName: "Test",
		AuthorID:        nonExistentAuthor,
	})

	assert.Equal(t, http.StatusNotFound, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestE2E_ReassignPRNotFound(t *testing.T) {
	// Тест на переназначение несуществующего PR
	nonExistentPR := randomString("pr_not_exists")
	someUser := randomString("u_some")

	code, body := sendRequest(t, "POST", "/pullRequest/reassign", PullRequestReassignReq{
		PullRequestID: nonExistentPR,
		OldUserID:     someUser,
	})

	assert.Equal(t, http.StatusNotFound, code)
	var errResp ErrorResponse
	json.Unmarshal(body, &errResp)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestE2E_SetIsActiveSuccess(t *testing.T) {
	// Тест на изменение активности пользователя
	teamName := randomString("team_active")
	userID := randomString("u_user")

	// Создаем команду с пользователем
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: userID, Username: "User", IsActive: true},
		},
	})

	// Деактивируем пользователя
	code, body := sendRequest(t, "POST", "/users/setIsActive", map[string]interface{}{
		"user_id":   userID,
		"is_active": false,
	})

	require.Equal(t, http.StatusOK, code)

	var resp struct {
		User struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			TeamName string `json:"team_name"`
			IsActive bool   `json:"is_active"`
		} `json:"user"`
	}
	json.Unmarshal(body, &resp)

	assert.Equal(t, userID, resp.User.UserID)
	assert.Equal(t, false, resp.User.IsActive)
	assert.Equal(t, teamName, resp.User.TeamName)
}

func TestE2E_GetTeamSuccess(t *testing.T) {
	// Тест на получение команды
	teamName := randomString("team_get")
	user1 := randomString("u1")
	user2 := randomString("u2")

	// Создаем команду
	sendRequest(t, "POST", "/team/add", TeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{UserID: user1, Username: "User1", IsActive: true},
			{UserID: user2, Username: "User2", IsActive: false},
		},
	})

	// Получаем команду
	code, body := sendRequest(t, "GET", fmt.Sprintf("/team/get?team_name=%s", teamName), nil)
	require.Equal(t, http.StatusOK, code)

	var team TeamRequest
	json.Unmarshal(body, &team)
	assert.Equal(t, teamName, team.TeamName)
	assert.Len(t, team.Members, 2)

	// Проверяем членов команды
	memberIDs := []string{team.Members[0].UserID, team.Members[1].UserID}
	assert.Contains(t, memberIDs, user1)
	assert.Contains(t, memberIDs, user2)
}
