package service

import (
	"context"
	"fmt"

	"github.com/mfajri11/family-catering-micro-service/auth-service/config"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/ports"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/apperrors"
	chttp "github.com/mfajri11/family-catering-micro-service/auth-service/pkg/client/http"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/consts"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/log"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/security"
)

const (
	actionType  = "at"
	refreshType = "rt"
	// passwordType = "pt"
)

type AuthService struct {
	authRepo  ports.IAuthRepository
	authCache ports.IAuthCache
	secure    security.ISecurity
	client    chttp.IRequester
}

func New(authRepo ports.IAuthRepository, authCache ports.IAuthCache, client chttp.IRequester) *AuthService {
	return &AuthService{authRepo: authRepo, authCache: authCache, client: client}
}

func (svc AuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	var (
		sid           string
		err           error
		errRT         error
		refreshToken  string
		accessToken   string
		multipleLogin bool
	)
	multipleLogin = true
	// TODO: GET user basic info using email (email must be in encrypted form because it is PII)
	req.Password, err = svc.secure.Decrypt(req.Password)
	if err != nil {
		err = fmt.Errorf("service.AuthService.Login: error decrypt password: %w", err)
		log.Error(err, "error decrypt password")
		err = apperrors.WrapError(err, apperrors.ErrInternalServer)
		return nil, err
	}

	url := fmt.Sprintf("%s/detail?eid=%s", config.Cfg.App.BaseURLOwner, req.Email)
	var userResp domain.OwnerResponse
	// TODO: find how to secure inter services communication
	err = svc.client.Get(url, &userResp)
	if err != nil {
		err = fmt.Errorf("service.AuthService.Login: error request user data by email: %w", err)
		log.Error(err, "error request user data by email")
		err = apperrors.WrapError(err, apperrors.ErrInternalServer)
		return nil, err
	}

	err = svc.secure.CompareHashPassword(req.Password, userResp.HashedPassword)
	if err != nil {
		err = fmt.Errorf("service.AuthService.Login: error compare password: %w", err)
		log.Error(err, "error compare password")
		err = apperrors.WrapError(err, apperrors.ErrInternalServer)
		return nil, err
	}

	sid, err = svc.authRepo.SessionIDByEmail(ctx, req.Email)
	if err == nil && sid == "" {
		sid, err = svc.secure.GenerateSID()
		// TODO: setex session get email
		refreshToken, errRT = svc.secure.GenerateToken(refreshType)
		multipleLogin = false
	}

	if err != nil {
		return nil, err // TODO: wrap error
	}

	if errRT != nil {
		return nil, errRT
	}

	// token always generated, even the user detected logged in multiple devices
	// and considered as safe because it has short life time
	accessToken, err = svc.secure.GenerateToken(actionType)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	err = svc.authRepo.InsertSession(ctx)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	return &domain.LoginResponse{
		SID:           sid,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		MultipleLogin: multipleLogin,
	}, nil

}

func (svc *AuthService) Deactivate(ctx context.Context, req domain.LogoutRequest) error {
	sid, ok := ctx.Value(consts.SidCtxKey).(string)
	if !ok {
		return fmt.Errorf("invalid sid type") // TODO: wrap error
	}

	token, ok := ctx.Value(consts.TokenCtxKey).(string)
	if !ok {
		return fmt.Errorf("invalid token type") // TODO: wrap error
	}
	_, err := svc.secure.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("error validate token") // TODO: wrap error
	}

	session, err := svc.authRepo.Session(ctx, sid)
	if err != nil {
		return fmt.Errorf("error get session by sid") // TODO: wrap error
	}
	if !session.IsValid {
		return fmt.Errorf("error invalid session") // TODO: wrap error
	}

	encryptedEmail, err := svc.secure.EncryptWithURLEncode(session.Email) // TODO: encrypt email
	if err != nil {
		return err // TODO: wrap error
	}
	url := fmt.Sprintf("%s/detail?eid=%s", config.Cfg.App.BaseURLOwner, encryptedEmail)
	var userResp domain.OwnerResponse
	err = svc.client.Get(url, &userResp)
	if err != nil {
		return err // TODO: wrap error
	}
	req.Password, err = svc.secure.Decrypt(req.Password)
	if err != nil {
		return fmt.Errorf("error decrypt password") // TODO: wrap error
	}
	err = svc.secure.ValidateRequest(&req)
	if err != nil {
		return fmt.Errorf("error validate request") // TODO: wrap error
	}

	err = svc.secure.CompareHashPassword(req.Password, userResp.HashedPassword)
	if err != nil {
		return fmt.Errorf("error compare password") // TODO: wrap error
	}
	// TODO: update user status
	// TODO: delete session (db & cache)
	return nil
}

func (svc *AuthService) Logout(ctx context.Context) error {
	sid, ok := ctx.Value(consts.SidCtxKey).(string)
	if !ok {
		return fmt.Errorf("invalid sid type") // TODO: wrap error
	}

	token, ok := ctx.Value(consts.TokenCtxKey).(string)
	if !ok {
		return fmt.Errorf("invalid token type") // TODO: wrap error
	}
	_, err := svc.secure.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("error validate token") // TODO: wrap error
	}

	session, err := svc.authRepo.Session(ctx, sid)
	if err != nil {
		return fmt.Errorf("error get session by sid") // TODO: wrap error
	}
	if !session.IsValid {
		return fmt.Errorf("error invalid session") // TODO: wrap error
	}

	// TODO: delete session (on cache and on persistance)

	return nil
}

// user does not have session (logged out)
func (svc *AuthService) ForgotPassword(ctx context.Context, req domain.ForgotPasswordRequest) error {
	// TODO: validate access token
	// TODO: decrypt password
	// TODO: validate request
	// TODO: get user data
	// TODO: compare password
	// TODO: send reset link password via email
	panic("not yet implemented")
}

func (svc *AuthService) ResetPassword(ctx context.Context, req domain.ForgotPasswordRequest) {
	// TODO: check session
	// TODO: check token
	// TODO: decrypt password
	// TODO: validate request
	// TODO: get user data
	// TODO: compare password
	// TODO: update password
	panic("not yet implemented")
}

func (svc *AuthService) RenewAccessToken(ctx context.Context) {
	// TODO: check session
	// TODO: get user data
	// TODO: generate access token
	panic("not yet implemented")
}

func (svc *AuthService) GenerateOTP(ctx context.Context) {
	panic("not yet implemented")
}
