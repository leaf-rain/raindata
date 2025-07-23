package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	v1 "github.com/leaf-rain/raindata/admin/api/admin/v1"
	"github.com/leaf-rain/raindata/admin/api/common"
	"github.com/leaf-rain/raindata/admin/internal/conf"
	"github.com/leaf-rain/raindata/common/str"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type IProto interface {
	GetHead() *common.BaseHead
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
			uid := str.Str2Int64(su.AuthorityId)
			head := proto.GetHead()
			//	//反射方式实现 反射还是觉得有点浪费性能，后续打算封装方法从ctx中获取head吧
			//if head == nil {
			//	val := reflect.ValueOf(req).Elem().FieldByName("Head")
			//	if val.Type().Kind() == reflect.Ptr && val.IsNil() {
			//		elemType := val.Type().Elem()
			//		newVal := reflect.New(elemType)
			//		val.Set(newVal)
			//	}
			//	head = proto.GetHead()
			//}
			if head != nil {
				head.Userid = uid
			}
			ctx = context.WithValue(ctx, "auth_user", su)
			// 从param获取用户信息
			h := header.RequestHeader()
			var bh = &common.BaseHead{}
			bh.Userid = uid
			bh.Appid = h.Get("appid")
			bh.Channel = h.Get("channel")
			bh.Version = h.Get("version")
			bh.Region = h.Get("region")
			bh.Ext = h.Get("ext")
			bh.UserType = common.USER_TYPE(str.Str2Int64(h.Get("userType")))
			bh.Gender = common.SEX_TYPE(str.Str2Int64(h.Get("gender")))
			ctx = context.WithValue(ctx, "head", bh)
		}
		var now = time.Now()
		reply, err := handler(ctx, req)
		log.Debugf("url is %v, req data %+v, rsp data %+v, time since:%v", header.Operation(), req, reply, time.Since(now))
		return reply, err
	}
}

func GetHead(ctx context.Context, proto IProto) *common.BaseHead {
	head := proto.GetHead()
	if head != nil {
		return head
	}
	su := ctx.Value("head")
	if su != nil {
		return su.(*common.BaseHead)
	}
	return nil
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
	rsp := new(v1.PermissionsReply)
	head := GetHead(ctx, req)
	js, _ := json.Marshal(head)
	rsp.Permissions = []string{string(js)}
	return rsp, nil
}

func (biz AuthBiz) GetPublicContent(ctx context.Context, req *v1.PublicContentReq) (*v1.PublicContentReply, error) {
	return &v1.PublicContentReply{Content: "111"}, nil
}
