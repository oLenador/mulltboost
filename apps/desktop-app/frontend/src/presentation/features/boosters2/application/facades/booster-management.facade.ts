import { LoadBoostersUseCase } from '../use-cases/load-boosters.use-case';
import { ManageStagingUseCase } from '../use-cases/manage-staging.use-case';
import { BoosterDataRepository } from '../../core/repositories/booster-data.repository';
import { BoosterEntity } from '../../core/entities/booster.entity';
import { StagedOperations } from '../../core/types/execution.types';

export interface BoosterPageConfig {
  category: string;
  language: string;
}

export interface BoosterManagementFacade {
  // Data loading
  loadBoosters(config: BoosterPageConfig): Promise<{ success: boolean; error?: string }>;
  getBoosters(): BoosterEntity[];
  getBooster(id: string): BoosterEntity | undefined;

  // Filtering
  searchBoosters(searchTerm: string): BoosterEntity[];
  getBoostersByStatus(applied: boolean): BoosterEntity[];
  getBoostersByRiskLevel(riskLevel: 'low' | 'medium' | 'high'): BoosterEntity[];

  // Staging operations
  stageOperation(boosterId: string, operation: 'apply' | 'revert'): void;
  stageBatch(operations: StagedOperations): void;
  clearStaging(): void;
  getStagedOperations(): StagedOperations;
  validateStaging(): { valid: boolean; conflicts: string[]; warnings: string[] };

  // Computed values
  getAppliedCount(): number;
  getStagedCount(): number;
  hasChanges(): boolean;
}

export function createBoosterManagementFacade(
  loadBoostersUseCase: LoadBoostersUseCase,
  manageStagingUseCase: ManageStagingUseCase,
  boosterRepository: BoosterDataRepository
): BoosterManagementFacade {
  return {
    // Data loading
    async loadBoosters(config: BoosterPageConfig) {
      const result = await loadBoostersUseCase.execute(config);
      return { success: result.success, error: result.error };
    },

    getBoosters() {
      return boosterRepository.getAllBoosters();
    },

    getBooster(id: string) {
      return boosterRepository.getBooster(id);
    },

    // Filtering
    searchBoosters(searchTerm: string) {
      return boosterRepository.searchBoosters(searchTerm);
    },

    getBoostersByStatus(applied: boolean) {
      return boosterRepository.getBoostersByStatus(applied);
    },

    getBoostersByRiskLevel(riskLevel: 'low' | 'medium' | 'high') {
      return boosterRepository.getBoostersByRiskLevel(riskLevel);
    },

    // Staging operations
    stageOperation(boosterId: string, operation: 'apply' | 'revert') {
      manageStagingUseCase.stageOperation({ boosterId, operation });
    },

    stageBatch(operations: StagedOperations) {
      manageStagingUseCase.stageBatch({ operations });
    },

    clearStaging() {
      manageStagingUseCase.clearStaging();
    },

    getStagedOperations() {
      return manageStagingUseCase.getStagedOperations();
    },

    validateStaging() {
      return manageStagingUseCase.validateStaging();
    },

    // Computed values
    getAppliedCount() {
      return boosterRepository.getBoostersByStatus(true).length;
    },

    getStagedCount() {
      return manageStagingUseCase.getStagedCount();
    },

    hasChanges() {
      return manageStagingUseCase.getStagedCount() > 0;
    },
  };
}