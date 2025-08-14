package systeminfo

import (
    "context"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    "github.com/oLenador/mulltbost/internal/core/ports/outbound"
)

type Service struct {
    infoRepo outbound.SystemInfoRepository
}

func NewService(infoRepo outbound.SystemInfoRepository) *Service {
    return &Service{
        infoRepo: infoRepo,
    }
}

func (s *Service) GetSystemInfo(ctx context.Context) (*entities.SystemInfo, error) {
    return s.getSystemInfo(ctx)
}

func (s *Service) GetHardwareInfo(ctx context.Context) (*entities.SystemInfo, error) {
    return s.getSystemInfo(ctx)
}

func (s *Service) RefreshSystemInfo(ctx context.Context) error {
    // Implementar cache refresh se necessário
    return nil
}

func (s *Service) getSystemInfo(ctx context.Context) (*entities.SystemInfo, error) {
   // os, err := s.infoRepo.GetOSInfo(ctx)
   // if err != nil {
   //     return nil, err
   // }

    // cpu, err := s.infoRepo.GetCPUInfo(ctx)
    // if err != nil {
    //     return nil, err
    // }
// 
    // memory, err := s.infoRepo.GetMemoryInfo(ctx)
    // if err != nil {
    //     return nil, err
    // }
// 
    // storage, err := s.infoRepo.GetStorageInfo(ctx)
    // if err != nil {
    //     return nil, err
    // }
// 
    // network, err := s.infoRepo.GetNetworkInfo(ctx)
    // if err != nil {
    //     return nil, err
    // }
    return nil, nil
    // return &entities.SystemInfo{
    //     OS:          *os,
    //     CPU:         *cpu,
    //     Memory:      *memory,
    //     Storage:     storage,
    //     Network:     network,
    // }, nil
}
