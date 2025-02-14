package integration

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Suite) TestBuy() {
	res := s.Auth("sasha", "gaiti")

	var resp struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(res.Body).Decode(&resp)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.Token)

	//-------
	item := "book"

	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/api/buy/"+item, nil)
	s.Require().NoError(err)
	var bearer = "Bearer " + resp.Token
	req.Header.Add("Authorization", bearer)

	res, err = s.server.Client().Do(req)
	s.Require().NoError(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	s.Require().Equal(http.StatusOK, res.StatusCode)
	s.Require().Empty(body)
}
