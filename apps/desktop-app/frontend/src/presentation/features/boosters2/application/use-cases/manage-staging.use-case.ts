import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';
import { BoosterDataRepository } from '../../core/repositories/booster-data.repository';
import { StagedOperations } from '../../core/types/execution.types';
import { BoosterEntity } from '../../core/entities/booster.entity';

export interface StageOperationRequest {
  boosterId: string;
  operation: 'apply' | 'revert';
}

export interface StageBatchRequest {
  operations: StagedOperations;
}

export interface ValidateStagingResponse {
  valid: boolean;
  conflicts: string[];
  warnings: string[];
}

export interface ManageStagingUseCase {
  stageOperation(request: StageOperationRequest): void;
  stageBatch(request: StageBatchRequest): void;
  clearStaging(): void;
  validateStaging(): ValidateStagingResponse;
  getStagedOperations(): StagedOperations;
  getStagedCount(): number;
}

export function createManageStagingUseCase(
  executionRepository: ExecutionStateRepository,
  boosterRepository: BoosterDataRepository
): ManageStagingUseCase {
  return {
    stageOperation(request: StageOperationRequest): void {
      const { boosterId, operation } = request;
      const booster = boosterRepository.getBooster(boosterId);
      
      if (!booster) {
        throw new Error(`Booster ${boosterId} not found`);
      }

      const current = executionRepository.getStagedOperations();
      const currentStaged = current[boosterId];

      // If operation would restore original state, remove from staging
      if ((operation === 'apply' && booster.isApplied) || 
          (operation === 'revert' && !booster.isApplied)) {
        executionRepository.removeStagedOperation(boosterId);
        return;
      }

      // Otherwise, stage the operation
      executionRepository.addStagedOperation(boosterId, operation);
    },

    stageBatch(request: StageBatchRequest): void {
      executionRepository.setStagedOperations(request.operations);
    },

    clearStaging(): void {
      executionRepository.clearStagedOperations();
    },

    validateStaging(): ValidateStagingResponse {
      const staged = executionRepository.getStagedOperations();
      const conflicts: string[] = [];
      const warnings: string[] = [];

      Object.entries(staged).forEach(([boosterId, operation]) => {
        const booster = boosterRepository.getBooster(boosterId);
        if (!booster) {
          conflicts.push(`Booster ${boosterId} not found`);
          return;
        }

        // Check reversibility for revert operations
        if (operation === 'revert' && !booster.reversible) {
          conflicts.push(`Booster ${booster.name} is not reversible`);
        }

        // Check dependencies and conflicts
        if (operation === 'apply') {
          booster.dependencies.forEach(depId => {
            const depBooster = boosterRepository.getBooster(depId);
            const depOperation = staged[depId];
            
            if (!depBooster?.isApplied && depOperation !== 'apply') {
              warnings.push(`${booster.name} depends on ${depBooster?.name || depId}`);
            }
          });

          booster.conflicts.forEach(conflictId => {
            const conflictBooster = boosterRepository.getBooster(conflictId);
            const conflictOperation = staged[conflictId];
            
            if (conflictBooster?.isApplied && conflictOperation !== 'revert') {
              conflicts.push(`${booster.name} conflicts with ${conflictBooster.name}`);
            }
          });
        }
      });

      return {
        valid: conflicts.length === 0,
        conflicts,
        warnings,
      };
    },

    getStagedOperations(): StagedOperations {
      return executionRepository.getStagedOperations();
    },

    getStagedCount(): number {
      return Object.keys(executionRepository.getStagedOperations()).length;
    },
  };
}