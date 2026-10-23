package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"gamification/internal/app"
	"gamification/internal/repository/memory"
	"gamification/internal/transport/http/apigen"
	"golang.org/x/crypto/bcrypt"
)

func TestTaskCRUD(t *testing.T) {
	handler, store := testStack(t)
	owner := register(t, handler, "tasks-owner@example.com", "longenough")
	other := register(t, handler, "tasks-other@example.com", "longenough")
	auth := "Bearer " + owner.Access
	foreign := "Bearer " + other.Access

	created := call(handler, http.MethodPost, "/api/v1/tasks", `{"title":"  first  ","reward_coins":10}`, auth, "")
	if created.Code != http.StatusCreated {
		t.Fatalf("create %d %s", created.Code, created.Body.String())
	}
	var task apigen.Task
	if err := json.Unmarshal(created.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}
	if task.Title != "first" || task.Description != "" || task.RewardCoins != 10 || task.IsDone || task.CompletedAt != nil {
		t.Fatalf("task = %+v", task)
	}
	if strings.Contains(created.Body.String(), "password") {
		t.Fatal("response has password")
	}

	second := call(handler, http.MethodPost, "/api/v1/tasks", `{"title":"second","description":"note","reward_coins":0}`, auth, "")
	if second.Code != http.StatusCreated {
		t.Fatalf("second %d %s", second.Code, second.Body.String())
	}

	list := call(handler, http.MethodGet, "/api/v1/tasks?limit=1", "", auth, "")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), `"title":"second"`) || strings.Contains(list.Body.String(), `"title":"first"`) {
		t.Fatalf("page %d %s", list.Code, list.Body.String())
	}
	open := call(handler, http.MethodGet, "/api/v1/tasks?status=open", "", auth, "")
	if open.Code != http.StatusOK || strings.Count(open.Body.String(), `"title"`) != 2 {
		t.Fatalf("open %s", open.Body.String())
	}
	otherList := call(handler, http.MethodGet, "/api/v1/tasks", "", foreign, "")
	if otherList.Code != http.StatusOK || otherList.Body.String() == "" || strings.Contains(otherList.Body.String(), "first") {
		t.Fatalf("other list %d %s", otherList.Code, otherList.Body.String())
	}

	got := call(handler, http.MethodGet, "/api/v1/tasks/"+task.Id.String(), "", auth, "")
	if got.Code != http.StatusOK {
		t.Fatalf("get %d %s", got.Code, got.Body.String())
	}
	hidden := call(handler, http.MethodGet, "/api/v1/tasks/"+task.Id.String(), "", foreign, "")
	missing := call(handler, http.MethodGet, "/api/v1/tasks/11111111-1111-1111-1111-111111111111", "", auth, "")
	if hidden.Code != http.StatusNotFound || hidden.Body.String() != missing.Body.String() {
		t.Fatalf("hidden %d %s missing %d %s", hidden.Code, hidden.Body.String(), missing.Code, missing.Body.String())
	}

	patched := call(handler, http.MethodPatch, "/api/v1/tasks/"+task.Id.String(), `{"title":"renamed"}`, auth, "")
	if patched.Code != http.StatusOK || !strings.Contains(patched.Body.String(), `"title":"renamed"`) || !strings.Contains(patched.Body.String(), `"reward_coins":10`) {
		t.Fatalf("patch %d %s", patched.Code, patched.Body.String())
	}

	store.MarkDone(task.Id.String())
	locked := call(handler, http.MethodPatch, "/api/v1/tasks/"+task.Id.String(), `{"reward_coins":11}`, auth, "")
	if locked.Code != http.StatusConflict {
		t.Fatalf("reward %d %s", locked.Code, locked.Body.String())
	}
	titleOK := call(handler, http.MethodPatch, "/api/v1/tasks/"+task.Id.String(), `{"title":"after"}`, auth, "")
	if titleOK.Code != http.StatusOK || !strings.Contains(titleOK.Body.String(), `"reward_coins":10`) {
		t.Fatalf("title after done %d %s", titleOK.Code, titleOK.Body.String())
	}
	removed := call(handler, http.MethodDelete, "/api/v1/tasks/"+task.Id.String(), "", auth, "")
	if removed.Code != http.StatusConflict {
		t.Fatalf("delete done %d %s", removed.Code, removed.Body.String())
	}

	var secondTask apigen.Task
	if err := json.Unmarshal(second.Body.Bytes(), &secondTask); err != nil {
		t.Fatal(err)
	}
	foreignDelete := call(handler, http.MethodDelete, "/api/v1/tasks/"+secondTask.Id.String(), "", foreign, "")
	if foreignDelete.Code != http.StatusNotFound {
		t.Fatalf("foreign delete %d %s", foreignDelete.Code, foreignDelete.Body.String())
	}
	deleted := call(handler, http.MethodDelete, "/api/v1/tasks/"+secondTask.Id.String(), "", auth, "")
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete %d %s", deleted.Code, deleted.Body.String())
	}
	gone := call(handler, http.MethodGet, "/api/v1/tasks/"+secondTask.Id.String(), "", auth, "")
	if gone.Code != http.StatusNotFound {
		t.Fatalf("gone %d", gone.Code)
	}

	complete := call(handler, http.MethodPost, "/api/v1/tasks/"+task.Id.String()+"/complete", "", auth, "")
	if complete.Code != http.StatusOK {
		t.Fatalf("complete %d %s", complete.Code, complete.Body.String())
	}
	profile := call(handler, http.MethodGet, "/api/v1/profile", "", auth, "")
	if !strings.Contains(profile.Body.String(), `"coins":0`) {
		t.Fatalf("coins changed: %s", profile.Body.String())
	}
}

func testStack(t *testing.T) (http.Handler, *memory.Store) {
	t.Helper()
	store := memory.New()
	svc := app.NewForTest(store, []byte(testSecret), "gamification", 15*time.Minute, time.Hour, bcrypt.MinCost)
	return Handler(discardLogger(), svc), store
}
