package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/leaf-rain/raindata/admin/api/admin/v1"
	"github.com/leaf-rain/raindata/admin/internal/biz"
	"github.com/leaf-rain/raindata/admin/internal/conf"
)

func NewAuthService(
	logger log.Logger,
	cdata *conf.Data,
	biz *biz.AuthBiz,
) *AuthService {
	return &AuthService{
		log:   log.NewHelper(log.With(logger, "module", "service/auth")),
		cdata: cdata,
		biz:   biz,
	}
}

type AuthService struct {
	pb.UnimplementedAuthServer

	log   *log.Helper
	cdata *conf.Data
	biz   *biz.AuthBiz
}

func (s *AuthService) GetPermissions(ctx context.Context, req *pb.PermissionsReq) (rsp *pb.PermissionsReply, err error) {
	return s.biz.GetPermissions(ctx, req)
}

func (s *AuthService) GetPublicContent(ctx context.Context, req *pb.PublicContentReq) (rsp *pb.PublicContentReply, err error) {
	return s.biz.GetPublicContent(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginReq) (*pb.User, error) {
	return s.biz.Login(ctx, req)
}
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutReply, error) {
	return s.biz.Logout(ctx, req)
}
