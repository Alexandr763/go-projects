package httpapi

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"gamification/internal/transport/http/apigen"
)

func TestCompleteAwardsOnce(t *testing.T) {
	handler, _ := testStack(t)
	owner := register(t, handler, "award-owner@example.com", "longenough")
	other := register(t, handler, "award-other@example.com", "longenough")
	auth := "Bearer " + owner.Access
	foreign := "Bearer " + other.Access

	created := call(handler, http.MethodPost, "/api/v1/tasks", `{"title":"ship","reward_coins":15}`, auth, "")
	if created.Code != http.StatusCreated {
		t.Fatalf("create %d %s", created.Code, created.Body.String())
	}
	var task apigen.Task
	if err := json.Unmarshal(created.Body.Bytes(), &task); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	codes := make(chan int, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rec := call(handler, http.MethodPost, "/api/v1/tasks/"+task.Id.String()+"/complete", "", auth, "")
			codes <- rec.Code
		}()
	}
	wg.Wait()
	close(codes)
	for code := range codes {
		if code != http.StatusOK {
			t.Fatalf("parallel complete %d", code)
		}
	}

	again := call(handler, http.MethodPost, "/api/v1/tasks/"+task.Id.String()+"/complete", "", auth, "")
	if again.Code != http.StatusOK {
		t.Fatalf("repeat %d %s", again.Code, again.Body.String())
	}
	var done apigen.Task
	if err := json.Unmarshal(again.Body.Bytes(), &done); err != nil {
		t.Fatal(err)
	}
	if !done.IsDone || done.CompletedAt == nil || done.RewardCoins != 15 {
		t.Fatalf("done = %+v", done)
	}

	profile := call(handler, http.MethodGet, "/api/v1/profile", "", auth, "")
	var got apigen.Profile
	if err := json.Unmarshal(profile.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Coins != 15 {
		t.Fatalf("coins = %d", got.Coins)
	}

	history := call(handler, http.MethodGet, "/api/v1/profile/transactions", "", auth, "")
	var list apigen.TransactionList
	if err := json.Unmarshal(history.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Transactions) != 1 || list.Transactions[0].Amount != 15 || list.Transactions[0].TaskId != task.Id {
		t.Fatalf("history = %+v", list.Transactions)
	}
	otherHistory := call(handler, http.MethodGet, "/api/v1/profile/transactions", "", foreign, "")
	if otherHistory.Code != http.StatusOK || otherHistory.Body.String() == "" || len(mustTransactions(t, otherHistory.Body.Bytes())) != 0 {
		t.Fatalf("other history %d %s", otherHistory.Code, otherHistory.Body.String())
	}

	hidden := call(handler, http.MethodPost, "/api/v1/tasks/"+task.Id.String()+"/complete", "", foreign, "")
	missing := call(handler, http.MethodPost, "/api/v1/tasks/11111111-1111-1111-1111-111111111111/complete", "", auth, "")
	if hidden.Code != http.StatusNotFound || hidden.Body.String() != missing.Body.String() {
		t.Fatalf("hidden %d %s missing %d %s", hidden.Code, hidden.Body.String(), missing.Code, missing.Body.String())
	}

	zero := call(handler, http.MethodPost, "/api/v1/tasks", `{"title":"free","reward_coins":0}`, auth, "")
	var free apigen.Task
	if err := json.Unmarshal(zero.Body.Bytes(), &free); err != nil {
		t.Fatal(err)
	}
	if rec := call(handler, http.MethodPost, "/api/v1/tasks/"+free.Id.String()+"/complete", "", auth, ""); rec.Code != http.StatusOK {
		t.Fatalf("zero %d %s", rec.Code, rec.Body.String())
	}
	profile = call(handler, http.MethodGet, "/api/v1/profile", "", auth, "")
	if err := json.Unmarshal(profile.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Coins != 15 {
		t.Fatalf("coins after zero = %d", got.Coins)
	}
	history = call(handler, http.MethodGet, "/api/v1/profile/transactions?limit=1", "", auth, "")
	list = apigen.TransactionList{}
	if err := json.Unmarshal(history.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Transactions) != 1 || list.Transactions[0].Amount != 0 {
		t.Fatalf("newest = %+v", list.Transactions)
	}
}

func mustTransactions(t *testing.T, body []byte) []apigen.Transaction {
	t.Helper()
	var list apigen.TransactionList
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatal(err)
	}
	return list.Transactions
}
