package data

import "github.com/leaf-rain/raindata/admin/internal/biz"

type AuthRepo struct {
	data *Data
}

func NewAuthRepo(data *Data) biz.AuthRepoI {
	return &AuthRepo{
		data: data,
	}
}
