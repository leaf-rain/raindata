package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	v1 "github.com/leaf-rain/raindata/admin/api/admin/v1"
	"github.com/leaf-rain/raindata/admin/api/common"
	"github.com/leaf-rain/raindata/admin/internal/conf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strconv"
	"time"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type IProto interface {
	GetHead() *common.BaseHead
	ProtoReflect() protoreflect.Message
}

type Auth struct {
	Hello string
}

// AuthRepoI is a Greater repo.
type AuthRepoI interface{}

func NewAuthBiz(cdata *conf.Data, repo AuthRepoI, logger log.Logger) *AuthBiz {
	return &AuthBiz{cdata: cdata, repo: repo, log: log.NewHelper(log.With(logger, "module", "biz/auth"))}
}

type AuthBiz struct {
	cdata *conf.Data
	repo  AuthRepoI
	log   *log.Helper
}

// NewWhiteListMatcher 创建jwt白名单
func (biz AuthBiz) NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/admin.v1.Auth/Login"] = struct{}{}
	whiteList["/admin.v1.Auth/Logout"] = struct{}{}
	whiteList["/admin.v1.Auth/Register"] = struct{}{}
	whiteList["/admin.v1.Auth/GetPublicContent"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

func (biz AuthBiz) AuthUser(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		header, ok := transport.FromServerContext(ctx)
		if !ok || header == nil {
			return nil, status.Error(codes.InvalidArgument, "argument error")
		}
		whiteList := biz.NewWhiteListMatcher()
		if whiteList(ctx, header.Operation()) {
			// 校验jst是否合法
			su := &AuthUser{}
			err := su.ParseFromContext(ctx, []byte(biz.cdata.GetJwt().GetSecret()))
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "token error")
			}
			// todo:判断jwt是否过期  su.Token
			proto, ok := req.(IProto)
			if !ok {
				return nil, status.Error(codes.InvalidArgument, "argument error")
			}
			head := proto.GetHead()
			if head == nil {
				//反射方式实现
				val := reflect.ValueOf(req).Elem().FieldByName("Head")
				if val.Type().Kind() == reflect.Ptr && val.IsNil() {
					elemType := val.Type().Elem()
					newVal := reflect.New(elemType)
					val.Set(newVal)
				}
				head = proto.GetHead()
			}
			head.Userid, _ = strconv.ParseInt(su.AuthorityId, 10, 64)
			ctx = context.WithValue(ctx, "auth_user", su)
		}
		var now = time.Now()
		reply, err := handler(ctx, req)
		log.Debugf("url is %v, req data %+v, rsp data %+v, time since:%v", header.Operation(), req, reply, time.Since(now))
		return reply, err
	}
}

func SetHead(proto IProto) {
	newVal := &common.BaseHead{} // 注意这里改为指针类型
	reflectMst := proto.ProtoReflect()
	fd := reflectMst.Descriptor().Fields().ByName("Head")
	if fd != nil && reflectMst.Get(fd).IsValid() {
		reflectMst.Set(fd, protoreflect.ValueOfMessage(newVal.ProtoReflect()))
	}
}

func (biz AuthBiz) Login(ctx context.Context, req *v1.LoginReq) (*v1.User, error) {
	var result = new(v1.User)
	var authModel = new(AuthUser)
	authModel.AuthorityId = "123456"
	token := authModel.CreateAccessJwtToken([]byte(biz.cdata.Jwt.GetSecret()))
	result.Token = token
	return result, nil
}

func (biz AuthBiz) Logout(ctx context.Context, req *v1.LogoutReq) (*v1.LogoutReply, error) {
	//TODO implement me
	panic("implement me")
}

func (biz AuthBiz) GetPermissions(ctx context.Context, req *v1.PermissionsReq) (*v1.PermissionsReply, error) {
	return nil, status.Error(codes.Unimplemented, "测试错误")
}

func (biz AuthBiz) GetPublicContent(ctx context.Context, req *v1.PublicContentReq) (*v1.PublicContentReply, error) {
	return &v1.PublicContentReply{Content: "111"}, nil
}
