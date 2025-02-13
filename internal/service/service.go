package service

import "merch/configure"

type Service struct {
	cfg          *configure.Config
	coinTransfer CoinTransfer
	userProvider UserProvider
	authorizer   Authorizer
}

func New(
	cfg *configure.Config,
	sender CoinTransfer,
	provider UserProvider,
	authorizer Authorizer,
) *Service {
	return &Service{
		cfg:          cfg,
		coinTransfer: sender,
		userProvider: provider,
		authorizer:   authorizer,
	}
}
