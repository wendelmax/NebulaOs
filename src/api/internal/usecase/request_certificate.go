package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type RequestCertificateUseCase struct {
	networkProvider domain.NetworkProvider
}

func NewRequestCertificateUseCase(np domain.NetworkProvider) *RequestCertificateUseCase {
	return &RequestCertificateUseCase{
		networkProvider: np,
	}
}

func (uc *RequestCertificateUseCase) Execute(ctx context.Context, domainName string) error {
	return uc.networkProvider.RequestCertificate(ctx, domainName)
}
