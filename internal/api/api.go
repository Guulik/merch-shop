package api

import (
	"github.com/go-playground/validator/v10"
)

type Api struct {
	authorizer   AuthorizerService
	buyer        Buyer
	infoProvider InfoProvider
	coinSender   CoinSender
}

func New(
	authorizer AuthorizerService,
	buyer Buyer,
	infoProvider InfoProvider,
	coinSender CoinSender,
) *Api {
	return &Api{
		authorizer:   authorizer,
		buyer:        buyer,
		infoProvider: infoProvider,
		coinSender:   coinSender,
	}
}

func validate(request interface{}) error {
	valid := validator.New(validator.WithRequiredStructEnabled())

	err := valid.Struct(request)
	if err != nil {
		return err
	}
	return nil
}
