package outbound

import (
    "context"

    "github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type GPUInfoRepository interface {
    GetGPUInfo(ctx context.Context) ([]entities.GPUInfo, error)
}