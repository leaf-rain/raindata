package domain

import (
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"go.uber.org/zap"
)

//go:generate wire

type Domain struct {
	logger        *zap.Logger
	repoMetadata  interface_repo.InterfaceMetadataRepo
	repoWriterMsg interface_repo.InterfaceWriterRepo
	repoId        interface_repo.InterfaceIdRepo
}

func NewDomain(
	logger *zap.Logger,
	repoMetadata interface_repo.InterfaceMetadataRepo,
	repoWriterMsg interface_repo.InterfaceWriterRepo,
	repoId interface_repo.InterfaceIdRepo,
) *Domain {
	return &Domain{
		logger:        logger,
		repoMetadata:  repoMetadata,
		repoWriterMsg: repoWriterMsg,
		repoId:        repoId,
	}
}
