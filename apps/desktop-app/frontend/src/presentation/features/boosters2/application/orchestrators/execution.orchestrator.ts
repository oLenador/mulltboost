import { ExecuteBatchUseCase, ExecuteBatchRequest } from '../use-cases/execute-batch.use-case';
import { SyncExecutionsUseCase } from '../use-cases/sync-executions.use-case';
import { ManageStagingUseCase } from '../use-cases/manage-staging.use-case';
import { StagedOperations } from '../../core/types/execution.types';

export interface ExecutionCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: string, progress?: number, error?: string) => void;
  onBatchStarted: (batchId: string) => void;
  onBatchCompleted: (batchId: string) => void;
  onBatchError: (batchId: string, error: string) => void;
}

export interface ExecutionOrchestrator {
  executeStagedBatch(): Promise<{ success: boolean; batchId?: string; error?: string }>;
  executeBatch(operations: StagedOperations): Promise<{ success: boolean; batchId?: string; error?: string }>;
  syncWithBackend(): Promise<void>;
  setCallbacks(callbacks: ExecutionCallbacks): void;
}

export function createExecutionOrchestrator(
  executeBatchUseCase: ExecuteBatchUseCase,
  syncExecutionsUseCase: SyncExecutionsUseCase,
  manageStagingUseCase: ManageStagingUseCase
): ExecutionOrchestrator {
  return {
    async executeStagedBatch() {
      const operations = manageStagingUseCase.getStagedOperations();
      
      if (Object.keys(operations).length === 0) {
        return { success: false, error: 'No operations staged' };
      }

      const result = await executeBatchUseCase.execute({ operations });
      
      if (result.success) {
        manageStagingUseCase.clearStaging();
      }
      
      return result;
    },

    async executeBatch(operations: StagedOperations) {
      return executeBatchUseCase.execute({ operations });
    },

    async syncWithBackend() {
      await syncExecutionsUseCase.execute();
    },

    setCallbacks(callbacks: ExecutionCallbacks) {
      executeBatchUseCase.setCallbacks({
        onBatchStarted: callbacks.onBatchStarted,
        onBatchCompleted: callbacks.onBatchCompleted,
        onBatchError: callbacks.onBatchError,
      });

      syncExecutionsUseCase.setCallbacks({
        onExecutionUpdated: callbacks.onExecutionStatusChanged,
        onSyncCompleted: () => {}, // Could add callback for sync completion
      });
    },
  };
}
