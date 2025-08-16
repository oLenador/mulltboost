package booster

import (
	"context"
	
	"github.com/oLenador/mulltbost/internal/core/domain/dto"
	"github.com/oLenador/mulltbost/internal/core/domain/entities"
	"github.com/oLenador/mulltbost/internal/core/domain/services/i18n"
	"github.com/oLenador/mulltbost/internal/core/ports/inbound"
)

// BaseBooster implementa funcionalidades comuns
type BaseBooster struct {
	entity       entities.Booster
	translations i18n.Translations
	i18nSvc      *i18n.Service
	executor     inbound.PlatformExecutor
}

func NewBaseBooster(entity entities.Booster, translations i18n.Translations, executor inbound.PlatformExecutor) *BaseBooster {
	i18nSvc := i18n.NewService()
	i18nSvc.SetTranslations(translations)
	
	return &BaseBooster{
		entity:       entity,
		translations: translations,
		i18nSvc:      i18nSvc,
		executor:     executor,
	}
}

func (b *BaseBooster) GetEntity() entities.Booster {
	return b.entity
}

func (b *BaseBooster) GetEntityDto(lang i18n.Language) dto.BoosterDto {
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

func (b *BaseBooster) Execute(ctx context.Context) (*entities.BoosterResult, error) {
	return b.executor.Execute(ctx, b.entity.ID)
}

func (b *BaseBooster) Validate(ctx context.Context) error {
	return b.executor.Validate(ctx)
}

func (b *BaseBooster) CanApply(ctx context.Context) bool {
	return b.executor.CanExecute(ctx)
}

func (b *BaseBooster) CanRevert(ctx context.Context) bool {
	return b.entity.Reversible && b.executor.CanExecute(ctx)
}

func (b *BaseBooster) Revert(ctx context.Context) (*entities.BoosterResult, error) {
	// Em implementação real, BackupData viria do repositório
	return b.executor.Revert(ctx, nil)
}