package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"sqat.yerzhanat.net/internal/data"
)

func newTestApplication(t *testing.T) *application {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &application{
		logger: logger,
		models: data.NewMockModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) deleteReq(t *testing.T, urlPath string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, urlPath string, data []byte) (int, http.Header, string) {
	reader := bytes.NewReader(data)
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", reader)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) patchReq(t *testing.T, urlPath string, data []byte) (int, http.Header, string) {
	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPatch, ts.URL+urlPath, reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer 12345678912345678912345678")

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) putReq(t *testing.T, urlPath string, data []byte) (int, http.Header, string) {
	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, ts.URL+urlPath, reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer 11111111111111111111111111")

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
