package service

import "time"

type Service struct {
	tokenTTL time.Duration
	//cfg          *configure.Config
	coinTransfer CoinTransfer
	userProvider UserProvider
	authorizer   Authorizer
}

func New(
	//прямая передача - это костыль
	tokenTTL time.Duration,
	//cfg *configure.Config,
	sender CoinTransfer,
	provider UserProvider,
	authorizer Authorizer,
) *Service {
	return &Service{
		tokenTTL: tokenTTL,
		//cfg:          cfg,
		coinTransfer: sender,
		userProvider: provider,
		authorizer:   authorizer,
	}
}
