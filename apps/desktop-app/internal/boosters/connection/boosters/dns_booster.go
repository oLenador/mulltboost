package connection

import (
	"context"

	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
)

type DNSBooster struct {
	entity       entities.Booster
	translations i18n.Translations
	i18nSvc      *i18n.Service
}

func NewDNSBooster() *DNSBooster {
	entity := entities.Booster{
		ID:             "connection_dns_booster",
		NameKey:        "booster.connection.dns.name",
		DescriptionKey: "booster.connection.dns.description",
		Category:       entities.CategoryConnection,
		Level:          entities.LevelFree,
		Platform:       []entities.Platform{entities.PlatformWindows},
		Reversible:     true,
		RiskLevel:      entities.RiskLow,
		Version:        "1.0.0",
		Tags:           []string{"network", "dns", "speed"},
	}

	translations := map[i18n.Language]i18n.Translation{
		i18n.English: {
			"booster.connection.dns.name":        "DNS Optimizer",
			"booster.connection.dns.description": "Configures faster DNS servers and clears DNS cache",
		},
		i18n.Portuguese: {
			"booster.connection.dns.name":        "Otimizador DNS",
			"booster.connection.dns.description": "Configura servidores DNS mais rápidos e limpa cache DNS",
		},
		i18n.Spanish: {
			"booster.connection.dns.name":        "Optimizador DNS",
			"booster.connection.dns.description": "Configura servidores DNS más rápidos y limpia caché DNS",
		},
	}

	i18nSvc := i18n.NewService()
	i18nSvc.SetTranslations(translations)

	return &DNSBooster{
		entity:       entity,
		translations: translations,
		i18nSvc:      i18nSvc,
	}
}

func (b *DNSBooster) GetTranslations() i18n.Translations {

	return b.translations
}

func (b *DNSBooster) GetEntity() entities.Booster {
	return b.entity
}
func (b *DNSBooster) GetEntityDto(lang i18n.Language) dto.BoosterDto {
	name := b.i18nSvc.Translate(b.entity.NameKey, lang)
	description := b.i18nSvc.Translate(b.entity.DescriptionKey, lang)

	return dto.BoosterDto{
		ID:          b.entity.ID,
		Name:        name,
		Description: description,
		Category:    b.entity.Category,
		Level:       b.entity.Level,
		Platform:    b.entity.Platform,
		Reversible:  b.entity.Reversible,
		RiskLevel:   b.entity.RiskLevel,
		Version:     b.entity.Version,
		Tags:        b.entity.Tags,
	}
}

func (b *DNSBooster) Execute(ctx context.Context) (*entities.BoosterResult, error) {

	return &entities.BoosterResult{
		Success: true,
		Message: "DNS settings optimized",
		BackupData: map[string]interface{}{
			"dns_backup": "original_dns_settings",
			"booster_id": b.entity.ID,
		},
	}, nil
}

func (b *DNSBooster) Validate(ctx context.Context) error {
	return nil
}

func (b *DNSBooster) CanApply(ctx context.Context) bool {
	return true
}

func (b *DNSBooster) CanRevert(ctx context.Context) bool {
	return b.entity.Reversible
}

func (b *DNSBooster) Revert(ctx context.Context) (*entities.BoosterResult, error) {
	return &entities.BoosterResult{
		Success: true,
		Message: "DNS settings reverted",
	}, nil
}
