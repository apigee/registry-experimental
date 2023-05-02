package edge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func environmentsTestServer(t *testing.T) *httptest.Server {
	m := http.NewServeMux()

	env := Environment{
		Name: "env-1",
	}

	m.HandleFunc("/environments", (func(w http.ResponseWriter, r *http.Request) {
		envs := []string{env.Name}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(envs); err != nil {
			t.Fatal(err)
		}
	}))

	m.HandleFunc("/environments/", (func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if !strings.Contains(r.URL.Path, "env-1") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(env); err != nil {
				t.Fatalf("want no error %v", err)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	return httptest.NewServer(m)
}

func TestEnvList(t *testing.T) {
	ts := environmentsTestServer(t)
	defer ts.Close()

	baseUrl, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := &EdgeClient{
		client:     http.DefaultClient,
		BaseURLEnv: baseUrl,
		BaseURL:    baseUrl,
	}
	es := &EnvironmentsServiceOp{
		client: client,
	}

	namelist, resp, e := es.ListNames()
	if e != nil {
		t.Errorf("while listing environments, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if len(namelist) <= 0 {
		t.Errorf("no environments found")
		return
	}
}
