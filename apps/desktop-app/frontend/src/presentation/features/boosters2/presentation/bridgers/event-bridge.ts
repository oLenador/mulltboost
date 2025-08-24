import { ExecutionOrchestrator } from '../../application/orchestrators/execution.orchestrator';
import { EventOrchestrator } from '../../application/orchestrators/event.orchestrator';
import { ExecutionStoreAdapter } from '../stores/adapters/execution-store.adapter';
import { BoosterStoreAdapter } from '../stores/adapters/booster-store.adapter';
import { StagedOperations } from '../../core/types/execution.types';
import { BatchEntity } from '../../core/entities/batch.entity';

export interface EventBridgeCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: string, progress?: number, error?: string) => void;
  onBoosterStatusChanged: (boosterId: string, isApplied: boolean) => void;
  onBatchStarted: (batch: BatchEntity) => void;
  onBatchCompleted: (batchId: string) => void;
  onBatchError: (batchId: string, error: string) => void;
  onSyncRequired: () => void | Promise<void>;
}

export interface EventBridge {
  connect: (
    executionAdapter: ExecutionStoreAdapter,
    boosterAdapter: BoosterStoreAdapter
  ) => void;
  disconnect: () => void;
  isConnected: () => boolean;
  
  // Operations
  stageOperation: (boosterId: string, operation: 'apply' | 'revert') => void;
  stageBatch: (operations: StagedOperations) => void;
  clearStaging: () => void;
  executeStagedBatch: () => Promise<{ success: boolean; batchId?: string; error?: string }>;
  executeBatch: (operations: StagedOperations) => Promise<{ success: boolean; batchId?: string; error?: string }>;
  
  // State queries
  getStagedOperations: () => StagedOperations;
  getStagedCount: () => number;
  getStats: () => any;
}

export function createEventBridge(
  executionOrchestrator: ExecutionOrchestrator,
  eventOrchestrator: EventOrchestrator,
  manageStagingUseCase: { 
    stageOperation: (params: { boosterId: string; operation: 'apply' | 'revert' }) => void;
    stageBatch: (params: { operations: StagedOperations }) => void;
    clearStaging: () => void;
    getStagedOperations: () => StagedOperations;
    getStagedCount: () => number;
  }
): EventBridge {
  let isConnected = false;
  let executionAdapter: ExecutionStoreAdapter | null = null;
  let boosterAdapter: BoosterStoreAdapter | null = null;
  let currentBatch: BatchEntity | null = null;

  const handleExecutionStatusChanged = (
    boosterId: string, 
    status: string, 
    progress?: number, 
    error?: string
  ) => {
    if (!executionAdapter) return;

    try {
      executionAdapter.updateExecution(boosterId, status, progress, error);
      
      // Update booster status for completed operations
      if (status === 'completed' && boosterAdapter) {
        const stagedOps = manageStagingUseCase.getStagedOperations();
        const operation = stagedOps[boosterId];
        if (operation) {
          const isApplied = operation === 'apply';
          boosterAdapter.updateBoosterStatus(boosterId, isApplied);
        }
      }
    } catch (error) {
      console.error('Error handling execution status change:', error);
    }
  };

  const handleBatchStarted = (batchId: string) => {
    if (!executionAdapter) return;

    try {
      // Create batch entity
      const batch: BatchEntity = {
        id: batchId,
        name: `Batch ${batchId}`,
        executionIds: Object.keys(manageStagingUseCase.getStagedOperations()),
        status: 'processing',
        totalProgress: 0,
        createdAt: new Date(),
        startedAt: new Date(),
      };

      currentBatch = batch;
      executionAdapter.setBatch(batch);

      // Add executions to store
      const stagedOps = manageStagingUseCase.getStagedOperations();
      executionAdapter.addExecutions(stagedOps);
    } catch (error) {
      console.error('Error handling batch start:', error);
    }
  };

  const handleBatchCompleted = (batchId: string) => {
    if (!executionAdapter || !currentBatch) return;

    try {
      const completedBatch: BatchEntity = {
        ...currentBatch,
        status: 'completed',
        totalProgress: 100,
        completedAt: new Date(),
      };

      executionAdapter.setBatch(completedBatch);
      
      // Clear batch after a delay to show completion
      setTimeout(() => {
        if (executionAdapter) {
          executionAdapter.clearBatch();
        }
        currentBatch = null;
      }, 2000);
    } catch (error) {
      console.error('Error handling batch completion:', error);
    }
  };

  const handleBatchError = (batchId: string, error: string) => {
    if (!executionAdapter || !currentBatch) return;

    try {
      const errorBatch: BatchEntity = {
        ...currentBatch,
        status: 'error',
        completedAt: new Date(),
      };

      executionAdapter.setBatch(errorBatch);
    } catch (error) {
      console.error('Error handling batch error:', error);
    }
  };

  const handleSyncRequired = async () => {
    try {
      // This could trigger a sync with backend or refresh booster data
      console.log('Sync required - refreshing state');
    } catch (error) {
      console.error('Error handling sync requirement:', error);
    }
  };

  const connect = (
    execAdapter: ExecutionStoreAdapter,
    boostAdapter: BoosterStoreAdapter
  ) => {
    if (isConnected) return;

    executionAdapter = execAdapter;
    boosterAdapter = boostAdapter;

    // Set up orchestrator callbacks
    executionOrchestrator.setCallbacks({
      onExecutionStatusChanged: handleExecutionStatusChanged,
      onBatchStarted: handleBatchStarted,
      onBatchCompleted: handleBatchCompleted,
      onBatchError: handleBatchError,
    });

    eventOrchestrator.setCallbacks({
      onExecutionStatusChanged: handleExecutionStatusChanged,
      onSyncRequired: handleSyncRequired,
    });

    // Start event orchestrator
    eventOrchestrator.start();
    isConnected = true;
  };

  const disconnect = () => {
    if (!isConnected) return;

    eventOrchestrator.stop();
    executionAdapter = null;
    boosterAdapter = null;
    currentBatch = null;
    isConnected = false;
  };

  const stageOperation = (boosterId: string, operation: 'apply' | 'revert') => {
    manageStagingUseCase.stageOperation({ boosterId, operation });
  };

  const stageBatch = (operations: StagedOperations) => {
    manageStagingUseCase.stageBatch({ operations });
  };

  const clearStaging = () => {
    manageStagingUseCase.clearStaging();
  };

  const executeStagedBatch = async () => {
    try {
      return await executionOrchestrator.executeStagedBatch();
    } catch (error) {
      console.error('Error executing staged batch:', error);
      return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
    }
  };

  const executeBatch = async (operations: StagedOperations) => {
    try {
      return await executionOrchestrator.executeBatch(operations);
    } catch (error) {
      console.error('Error executing batch:', error);
      return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
    }
  };

  const getStagedOperations = () => {
    return manageStagingUseCase.getStagedOperations();
  };

  const getStagedCount = () => {
    return manageStagingUseCase.getStagedCount();
  };

  const getStats = () => {
    return {
      event: eventOrchestrator.getStats(),
      staged: {
        count: manageStagingUseCase.getStagedCount(),
        operations: manageStagingUseCase.getStagedOperations(),
      },
      batch: currentBatch,
    };
  };

  return {
    connect,
    disconnect,
    isConnected: () => isConnected,
    stageOperation,
    stageBatch,
    clearStaging,
    executeStagedBatch,
    executeBatch,
    getStagedOperations,
    getStagedCount,
    getStats,
  };
}