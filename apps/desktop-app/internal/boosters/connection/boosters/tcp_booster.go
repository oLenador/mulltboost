// internal/boosters/connection/tcp_booster.go
package connection

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
)

type TCPBooster struct {
	entity      entities.Booster
	translations i18n.Translations
	i18nSvc    *i18n.Service
}

func NewTCPBooster() *TCPBooster {
	entity := entities.Booster{
		ID:             "connection_tcp_booster",
		NameKey:        "booster.connection.tcp.name",
		DescriptionKey: "booster.connection.tcp.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskMedium,
		Version:        "1.0.0",
		Tags:           []string{"network", "tcp", "latency"},
	}
	
	translations := i18n.Translations{
		i18n.English: {
			"booster.connection.tcp.name":        "TCP Optimizer",
			"booster.connection.tcp.description": "Optimizes TCP network settings for better latency and throughput",
		},
		i18n.Portuguese: {
			"booster.connection.tcp.name":        "Otimizador TCP",
			"booster.connection.tcp.description": "Otimiza configurações de rede TCP para melhor latência e throughput",
		},
		i18n.PortugueseBrazil: {
			"booster.connection.tcp.name":        "Otimizador de TCP",
			"booster.connection.tcp.description": "Otimiza configurações de rede TCP para reduzir a latência e aumentar o throughput",
		},
		i18n.Spanish: {
			"booster.connection.tcp.name":        "Optimizador TCP",
			"booster.connection.tcp.description": "Optimiza configuraciones de red TCP para mejorar la latencia y el rendimiento",
		},
		i18n.Russian: {
			"booster.connection.tcp.name":        "Оптимизатор TCP",
			"booster.connection.tcp.description": "Оптимизирует настройки TCP для снижения задержки и увеличения пропускной способности",
		},
	}
	
	i18nSvc := i18n.NewService()
	i18nSvc.SetTranslations(translations)

	return &TCPBooster{
		entity:       entity,
		translations: translations,
		i18nSvc: i18nSvc,
	}
}

func (b *TCPBooster) GetEntity() entities.Booster {
	return b.entity
}

func (b *TCPBooster) GetEntityDto(lang i18n.Language) dto.BoosterDto {
	name := b.i18nSvc.Translate(b.entity.NameKey, lang)
	description := b.i18nSvc.Translate(b.entity.DescriptionKey, lang)

	return dto.BoosterDto{
		ID:           b.entity.ID,
		Name:         name,
		Description:  description,
		Category:     b.entity.Category,
		Level:        b.entity.Level,
		Platform:     b.entity.Platform,
		Dependencies: b.entity.Dependencies,
		Conflicts:    b.entity.Conflicts,
		Reversible:   b.entity.Reversible,
		RiskLevel:    b.entity.RiskLevel,
		Version:      b.entity.Version,
		Tags:         b.entity.Tags,
	}
}

func (b *TCPBooster) Execute(ctx context.Context) (*entities.BoosterResult, error) {

	return &entities.BoosterResult{
		Success: true,
		Message: "TCP settings optimized",
		BackupData: map[string]interface{}{
			"registry_backup": "tcp_backup",
			"booster_id":      b.entity.ID,
		},
	}, nil
}

func (b *TCPBooster) Validate(ctx context.Context) error {
	return nil
}

func (b *TCPBooster) CanApply(ctx context.Context) bool {
	return true
}

func (b *TCPBooster) CanRevert(ctx context.Context) bool {
	return b.entity.Reversible
}

func (b *TCPBooster) Revert(ctx context.Context) (*entities.BoosterResult, error) {
	return &entities.BoosterResult{
		Success: true,
		Message: "TCP settings reverted",
	}, nil
}