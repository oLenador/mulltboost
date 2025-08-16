package windows

import (
	"context"
)

type AudioOptimizationService interface {
	// Configuração de áudio
	OptimizeAudioForGaming(ctx context.Context) error
	OptimizeAudioForMusic(ctx context.Context) error
	RestoreDefaultAudioSettings(ctx context.Context) error
	
	// Configurações específicas
	SetAudioLatency(ctx context.Context, latency int) error
	SetSampleRate(ctx context.Context, sampleRate int) error
	SetBitDepth(ctx context.Context, bitDepth int) error
	
	// Status
	GetCurrentAudioConfiguration(ctx context.Context) (*entities.AudioConfiguration, error)
	IsAudioOptimized(ctx context.Context) (bool, error)
}
