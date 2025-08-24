import { ExecutionOrchestrator } from '../orchestrators/execution.orchestrator';
import { EventOrchestrator } from '../orchestrators/event.orchestrator';
import { ManageStagingUseCase } from '../use-cases/manage-staging.use-case';
import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';
import { StagedOperations } from '../../core/types/execution.types';
import { ExecutionEntity } from '../../core/entities/execution.entity';
import { BatchEntity } from '../../core/entities/batch.entity';

export interface ExecutionFacadeCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: string, progress?: number, error?: string) => void;
  onSyncRequired: () => void | Promise<void>;
  onBatchStarted: (batchId: string) => void;
  onBatchCompleted: (batchId: string) => void;
  onBatchError: (batchId: string, error: string) => void;
}

export interface ExecutionFacade {
  // Lifecycle
  initialize(): void;
  destroy(): void;
  isRunning(): boolean;

  // Execution
  executeStagedBatch(): Promise<{ success: boolean; batchId?: string; error?: string }>;
  executeBatch(operations: StagedOperations): Promise<{ success: boolean; batchId?: string; error?: string }>;
  syncWithBackend(): Promise<void>;

  // Staging
  stageOperation(boosterId: string, operation: 'apply' | 'revert'): void;
  stageBatch(operations: StagedOperations): void;
  clearStaging(): void;
  getStagedOperations(): StagedOperations;
  getStagedCount(): number;

  // State
  getExecutions(): ExecutionEntity[];
  getCurrentBatch(): BatchEntity | undefined;
  getStats(): any;

  // Callbacks
  setCallbacks(callbacks: ExecutionFacadeCallbacks): void;
}

export function createExecutionFacade(
  executionOrchestrator: ExecutionOrchestrator,
  eventOrchestrator: EventOrchestrator,
  manageStagingUseCase: ManageStagingUseCase,
  executionRepository: ExecutionStateRepository
): ExecutionFacade {
  let isInitialized = false;

  return {
    // Lifecycle
    initialize() {
      if (isInitialized) return;

      eventOrchestrator.start();
      isInitialized = true;
    },

    destroy() {
      if (!isInitialized) return;

      eventOrchestrator.stop();
      executionRepository.clearExecutions();
      executionRepository.clearStagedOperations();
      isInitialized = false;
    },

    isRunning() {
      return isInitialized && eventOrchestrator.isRunning();
    },

    // Execution
    async executeStagedBatch() {
      if (!isInitialized) {
        throw new Error('Execution facade not initialized');
      }
      return executionOrchestrator.executeStagedBatch();
    },

    async executeBatch(operations: StagedOperations) {
      if (!isInitialized) {
        throw new Error('Execution facade not initialized');
      }
      return executionOrchestrator.executeBatch(operations);
    },

    async syncWithBackend() {
      if (!isInitialized) return;
      return executionOrchestrator.syncWithBackend();
    },

    // Staging
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

    getStagedCount() {
      return manageStagingUseCase.getStagedCount();
    },

    // State
    getExecutions() {
      return executionRepository.getAllExecutions();
    },

    getCurrentBatch() {
      return executionRepository.getCurrentBatch();
    },

    getStats() {
      return {
        executions: executionRepository.getExecutionStats(),
        events: eventOrchestrator.getStats(),
      };
    },

    // Callbacks
    setCallbacks(callbacks: ExecutionFacadeCallbacks) {
      executionOrchestrator.setCallbacks({
        onExecutionStatusChanged: callbacks.onExecutionStatusChanged,
        onBatchStarted: callbacks.onBatchStarted,
        onBatchCompleted: callbacks.onBatchCompleted,
        onBatchError: callbacks.onBatchError,
      });

      eventOrchestrator.setCallbacks({
        onExecutionStatusChanged: callbacks.onExecutionStatusChanged,
        onSyncRequired: callbacks.onSyncRequired,
      });
    },
  };
}
