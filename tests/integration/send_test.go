//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type user struct {
	Username string
	Password string
}

func (s *Suite) TestSend() {
	sasha := user{
		Username: "sasha",
		Password: "gaiti",
	}
	vasa := user{
		Username: "vasa",
		Password: "mayami",
	}

	userFrom := s.Auth(sasha.Username, sasha.Password)
	_ = s.Auth(vasa.Username, vasa.Password)

	var resp struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(userFrom.Body).Decode(&resp)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.Token)

	//-------
	coinAmount := 50
	body, _ := json.Marshal(map[string]interface{}{
		"toUser": vasa.Username,
		"amount": coinAmount,
	})
	req, err := http.NewRequest(http.MethodPost, s.server.URL+"/api/sendCoin", bytes.NewBuffer(body))
	s.Require().NoError(err)
	var bearer = "Bearer " + resp.Token
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.server.Client().Do(req)
	s.Require().NoError(err)

	defer res.Body.Close()
	s.Require().Equal(http.StatusOK, res.StatusCode)
}
