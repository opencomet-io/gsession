package gsession

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/opencomet-io/gsession/memstore"
)

type Service struct {
	AccessManager  *Manager
	RefreshManager *Manager
}

func NewService() *Service {
	store := memstore.New()

	srvc := &Service{
		AccessManager: &Manager{
			Store:       store,
			Codec:       GobCodec{},
			TokenPrefix: "",
			TokenLength: 32,
			Lifetime:    20 * time.Minute,
		},
		RefreshManager: &Manager{
			Store:       store,
			Codec:       GobCodec{},
			TokenPrefix: "refresh:",
			TokenLength: 48,
			Lifetime:    20 * 24 * time.Hour,
		},
	}

	return srvc
}

func (srvc *Service) InitSessionSet(ctx context.Context, initial, proof map[string]any) (string, string, error) {
	refreshValues := map[string]any{
		"data":  initial,
		"proof": proof,
	}
	refreshToken, err := srvc.RefreshManager.InitSession(ctx, refreshValues)
	if err != nil {
		return "", "", err
	}

	accessValues := map[string]any{
		"data": initial,
		"from": refreshToken,
	}
	accessToken, err := srvc.AccessManager.InitSession(ctx, accessValues)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (srvc *Service) InvalidateSessionSet(ctx context.Context, accessToken, refreshToken string) error {
	if err := srvc.RefreshManager.InvalidateSession(ctx, refreshToken); err != nil {
		return err
	}

	if err := srvc.AccessManager.InvalidateSession(ctx, accessToken); err != nil {
		return err
	}

	return nil
}

func (srvc *Service) RenewTokens(ctx context.Context, accessToken, refreshToken string) (string, string, error) {
	newRefreshToken, err := srvc.RefreshManager.RenewToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := srvc.AccessManager.RenewToken(ctx, accessToken)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (srvc *Service) RequestAccessSession(ctx context.Context, refreshToken string, proof map[string]any) (string, bool, error) {
	refreshVals, err := srvc.RefreshManager.GetSessionValues(ctx, refreshToken)
	if err != nil {
		return "", false, err
	}

	expectedProof, ok := refreshVals["proof"]
	if !ok {
		return "", false, nil
	}

	if !reflect.DeepEqual(expectedProof, proof) {
		return "", false, nil
	}

	dataAny, ok := refreshVals["data"]
	if !ok {
		return "", false, nil
	}

	data, ok := dataAny.(map[string]any)
	if !ok {
		return "", false, nil
	}

	accessVals := map[string]any{
		"data": data,
		"from": refreshToken,
	}
	accessToken, err := srvc.AccessManager.InitSession(ctx, accessVals)
	if err != nil {
		return "", false, err
	}

	return accessToken, true, nil
}

func (srvc *Service) GetSessionValues(ctx context.Context, accessToken string) (map[string]any, error) {
	vals, err := srvc.AccessManager.GetSessionValues(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	dataAny, ok := vals["data"]
	if !ok {
		return nil, errors.New("gsession: corrupted session data")
	}

	data, ok := dataAny.(map[string]any)
	if !ok {
		return nil, errors.New("gsession: corrupted session data")
	}

	return data, nil
}

func (srvc *Service) SetSessionValues(ctx context.Context, accessToken string, vals map[string]any) error {
	accessVals, err := srvc.AccessManager.GetSessionValues(ctx, accessToken)
	if err != nil {
		return err
	}

	refreshTokenAny, ok := accessVals["from"]
	if !ok {
		return errors.New("gsession: corrupted session data")
	}

	refreshToken, ok := refreshTokenAny.(string)
	if !ok {
		return errors.New("gsession: corrupted session data")
	}

	newAccessVals := map[string]any{
		"data": vals,
		"from": refreshToken,
	}
	if err := srvc.AccessManager.SetSessionValues(ctx, accessToken, newAccessVals); err != nil {
		return err
	}

	refreshVals, expiry, found, err := srvc.RefreshManager.retrieveSession(ctx, refreshToken)
	if err != nil {
		return err
	}
	if !found {
		return ErrNotFound
	}

	proofAny, ok := refreshVals["proof"]
	if !ok {
		return errors.New("gsession: corrupted session data")
	}

	newRefreshVals := map[string]any{
		"data":  vals,
		"proof": proofAny,
	}
	if err := srvc.RefreshManager.insertSession(ctx, refreshToken, newRefreshVals, expiry); err != nil {
		return err
	}

	return nil
}
