package service

import "merch/configure"

type Service struct {
	//log
	cfg          *configure.Config
	buyer        Buyer
	coinTransfer CoinTransfer
	userProvider UserProvider
	authorizer   Authorizer
}

func New(
	cfg *configure.Config,
	buyer Buyer,
	sender CoinTransfer,
	provider UserProvider,
	authorizer Authorizer,
) *Service {
	return &Service{
		cfg:          cfg,
		buyer:        buyer,
		coinTransfer: sender,
		userProvider: provider,
		authorizer:   authorizer,
	}
}
