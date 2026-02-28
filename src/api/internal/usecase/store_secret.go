package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type StoreSecretUseCase struct {
	secretManager domain.SecretManager
}

func NewStoreSecretUseCase(sm domain.SecretManager) *StoreSecretUseCase {
	return &StoreSecretUseCase{
		secretManager: sm,
	}
}

func (uc *StoreSecretUseCase) Execute(ctx context.Context, key, value string) error {
	return uc.secretManager.StoreSecret(ctx, key, value)
}
