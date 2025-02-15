//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
)

func (s *Suite) TestInfo() {
	res := s.Auth("sasha", "gaiti")

	var resp struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(res.Body).Decode(&resp)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.Token)

	//-------
	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/api/info", nil)
	s.Require().NoError(err)
	var bearer = "Bearer " + resp.Token
	req.Header.Add("Authorization", bearer)

	res, err = s.server.Client().Do(req)
	s.Require().NoError(err)

	defer res.Body.Close()
	s.Require().Equal(http.StatusOK, res.StatusCode)
}
