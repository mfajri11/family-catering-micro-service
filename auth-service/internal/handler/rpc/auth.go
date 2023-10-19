package rpc

import (
	"context"
	"fmt"

	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/ports"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/handler/rpc/pb"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/log"
)

type AuthRPCHandler struct {
	pb.UnimplementedAuthServer
	svc ports.IAuthService
}

func NewAuthHandler(authService ports.IAuthService) *AuthRPCHandler {
	return &AuthRPCHandler{
		svc: authService,
	}
}

func (h AuthRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	reqMapped := domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	respSvc, err := h.svc.Login(ctx, reqMapped)
	if err != nil {
		err := fmt.Errorf("rpc.AuthRPCHandler.Login: error login: %w", err)
		log.Error(err, "error login")
		return nil, err // TODO: wrap error
	}

	resp = &pb.LoginResponse{
		AccessToken:  respSvc.AccessToken,
		RefreshToken: respSvc.RefreshToken,
	}

	setHeader(ctx, headerOptions{
		"X-Status-Code": "200",
		"X-Sid":         respSvc.SID,
		"X-Path":        "/api/v1/login",
	})

	return resp, nil
}
