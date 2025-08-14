package booster

type Strategy interface {
	Execute(ctx context.Context, config Config) (*Result, error)
	Validate(config Config) error
	GetPlatformRequirements() Requirements
}