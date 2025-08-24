import { ExecutionService } from '../../core/services/execution.service';
import { BoosterApiService } from '../../core/services/booster-api.service';
import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';
import { StagedOperations } from '../../core/types/execution.types';
import { BatchEntity } from '../../core/entities/batch.entity';

export interface ExecuteBatchRequest {
  operations: StagedOperations;
}

export interface ExecuteBatchResponse {
  success: boolean;
  batchId?: string;
  error?: string;
}

export interface BatchCallbacks {
  onBatchStarted: (batchId: string) => void;
  onBatchCompleted: (batchId: string) => void;
  onBatchError: (batchId: string, error: string) => void;
}

export interface ExecuteBatchUseCase {
  execute(request: ExecuteBatchRequest): Promise<ExecuteBatchResponse>;
  setCallbacks(callbacks: BatchCallbacks): void;
}

export function createExecuteBatchUseCase(
  executionService: ExecutionService,
  apiService: BoosterApiService,
  repository: ExecutionStateRepository
): ExecuteBatchUseCase {
  let callbacks: BatchCallbacks | null = null;

  return {
    async execute(request: ExecuteBatchRequest): Promise<ExecuteBatchResponse> {
      const { operations } = request;

      try {
        // Validate operations
        const validation = executionService.canExecuteOperations(operations);
        if (!validation.canExecute) {
          return { success: false, error: validation.reason };
        }

        // Create executions and batch
        const executions = executionService.createExecutionsFromOperations(operations);
        const batch = executionService.createBatchFromExecutions(executions);

        // Store executions
        executions.forEach(execution => repository.setExecution(execution));
        repository.setCurrentBatch(batch);

        // Notify batch started
        callbacks?.onBatchStarted(batch.id);

        // Execute operations sequentially
        for (const [boosterId, operation] of Object.entries(operations)) {
          try {
            const result = await apiService.executeBooster(boosterId, operation);
            if (!result.success) {
              console.error(`Failed to execute ${operation} on ${boosterId}:`, result.error);
            }
            
            // Small delay to prevent overload
            await new Promise(resolve => setTimeout(resolve, 100));
          } catch (error) {
            console.error(`Error executing ${operation} on ${boosterId}:`, error);
          }
        }

        // Notify batch completed
        callbacks?.onBatchCompleted(batch.id);

        return { success: true, batchId: batch.id };
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Batch execution failed';
        
        // Notify batch error
        const batchId = executionService.generateBatchId();
        callbacks?.onBatchError(batchId, errorMessage);

        return { success: false, error: errorMessage };
      }
    },

    setCallbacks(newCallbacks: BatchCallbacks): void {
      callbacks = newCallbacks;
    },
  };
}