package integration

import (
	"bytes"
	"fmt"
	"net/http"
)

func (s *Suite) TestAuth() {
	res := s.Auth("sasha", "gaiti")
	res.Body.Close()
}

func (s *Suite) Auth(username, password string) *http.Response {

	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	req, err := http.NewRequest(http.MethodPost, s.server.URL+"/api/auth", bytes.NewBuffer([]byte(body)))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	res, err := s.server.Client().Do(req)
	s.Require().NoError(err)
	return res
}
