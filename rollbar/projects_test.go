package rollbar

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestListProjects(t *testing.T) {
	teardown := setup()
	defer teardown()

	handURL := "/projects"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/list.json"))
	})

	response, err := client.ListProjects()

	if err != nil {
		t.Fatal(err)
	}

	expected := ListProjectsResponse{
		Error: 0,
		Result: []struct {
			AccountID    int    `json:"account_id"`
			ID           int    `json:"id"`
			DateCreated  int    `json:"date_created"`
			DateModified int    `json:"date_modified"`
			Name         string `json:"name"`
		}{
			{
				ID:           12112,
				AccountID:    8608,
				DateCreated:  1407933721,
				DateModified: 1457475137,
				Name:         "",
			},
			{
				ID:           106671,
				AccountID:    8608,
				DateCreated:  1489139046,
				DateModified: 1549293583,
				Name:         "Client-Config",
			},
			{
				ID:           12116,
				AccountID:    8608,
				DateCreated:  1407933922,
				DateModified: 1556814300,
				Name:         "My",
			},
		},
	}

	if !reflect.DeepEqual(*response, expected) {
		t.Errorf("expected response %v, got %v.", *response, expected)
	}
}
