import { BoosterApiService } from '../../core/services/booster-api.service';
import { BoosterEntity, createBoosterEntity } from '../../core/entities/booster.entity';
import { ApiResponse, QueueStateData, createQueueStateData, createApiResponse } from '../../core/types/api.types';
import { 
  InitBoosterApply, 
  InitRevertBooster, 
  GetExecutionQueueState, 
  GetBoostersByCategory 
} from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler';
import { Language } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/services/i18n';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

function mapRiskLevelToImpact(riskLevel: string): 'low' | 'medium' | 'high' {
  const normalizedLevel = (riskLevel || '').toLowerCase();
  return ['low', 'medium', 'high'].includes(normalizedLevel) 
    ? normalizedLevel as 'low' | 'medium' | 'high'
    : 'low';
}

export function createBackendClient(): BoosterApiService {
  return {
    async executeBooster(boosterId: string, operation: 'apply' | 'revert'): Promise<ApiResponse> {
      try {
        const result = operation === 'apply' 
          ? await InitBoosterApply(boosterId) 
          : await InitRevertBooster(boosterId);

        if (!result.Success) {
          const error = result.Error || result.Message || `Failed to ${operation} booster`;
          return createApiResponse(false, undefined, error);
        }

        return createApiResponse(
          true, 
          undefined, 
          undefined, 
          `Successfully initiated ${operation}`,
          result.OperationID
        );
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : `Failed to ${operation} booster`;
        return createApiResponse(false, undefined, errorMessage);
      }
    },

    async getQueueState(): Promise<QueueStateData> {
      try {
        const stateOrArray = await GetExecutionQueueState();
        return createQueueStateData(stateOrArray);
      } catch (error) {
        console.error('Failed to get queue state:', error);
        return createQueueStateData([]);
      }
    },

    async getBoostersByCategory(category: string, language: string): Promise<BoosterEntity[]> {
      try {
        const response = await GetBoostersByCategory(category as BoosterCategory, language as Language);
        return response.map(dto => createBoosterEntity({
          id: dto.ID,
          name: dto.Name,
          description: dto.Description,
          category: dto.Category,
          level: dto.Level,
          platform: dto.Platform.join(', '),
          dependencies: dto.Dependencies || [],
          conflicts: dto.Conflicts || [],
          reversible: dto.Reversible,
          riskLevel: mapRiskLevelToImpact(dto.RiskLevel),
          version: dto.Version,
          isApplied: dto.IsApplied,
          appliedAt: dto.AppliedAt ? new Date(dto.AppliedAt.toString()) : new Date(0),
          revertedAt: dto.RevertedAt ? new Date(dto.RevertedAt.toString()) : new Date(0),
          tags: dto.Tags || [],
        }));
      } catch (error) {
        console.error('Error loading boosters:', error);
        return [];
      }
    },
  };
}