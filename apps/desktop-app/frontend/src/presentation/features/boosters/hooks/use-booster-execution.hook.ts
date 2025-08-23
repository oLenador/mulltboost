// src/presentation/features/boosters/hooks/use-booster-execution.hook.ts
import { useCallback, useEffect, useRef } from 'react';
import { useAtom } from 'jotai';
import { 
  boosterExecutionsAtom, 
  currentBatchAtom, 
  executionStatsAtom, 
  executionsByStatusAtom, 
  isExecutingAtom, 
  canExecuteBatchAtom, 
  stagedOperationsAtom,
  updateExecutionAtom, 
  addExecutionsAtom, 
  removeExecutionsAtom, 
  clearExecutionsAtom, 
  stageOperationAtom, 
  stageBatchOperationsAtom, 
  clearStagingAtom, 
  StagedOperations 
} from '@/core/store/booster-execution.store';
import { BoosterExecutionService, ExecutionServiceCallbacks } from '../domain/booster-execution.service';
import { ExecutionBatch, ExecutionStatus } from '../domain/booster-queue.types';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

export const useBoosterExecution = () => {
  // Atoms
  const [executions] = useAtom(boosterExecutionsAtom);
  const [currentBatch, setCurrentBatch] = useAtom(currentBatchAtom);
  const [stats] = useAtom(executionStatsAtom);
  const [executionsByStatus] = useAtom(executionsByStatusAtom);
  const [isExecuting] = useAtom(isExecutingAtom);
  const [canExecute] = useAtom(canExecuteBatchAtom);
  const [stagedOperations] = useAtom(stagedOperationsAtom);

  // Actions
  const [, updateExecution] = useAtom(updateExecutionAtom);
  const [, addExecutions] = useAtom(addExecutionsAtom);
  const [, removeExecutions] = useAtom(removeExecutionsAtom);
  const [, clearExecutions] = useAtom(clearExecutionsAtom);
  const [, stageOperation] = useAtom(stageOperationAtom);
  const [, stageBatchOperations] = useAtom(stageBatchOperationsAtom);
  const [, clearStaging] = useAtom(clearStagingAtom);

  const serviceRef = useRef<BoosterExecutionService | null>(null);
  const isInitializedRef = useRef(false);
  const syncDebounceRef = useRef<NodeJS.Timeout | null>(null);

  const onExecutionStatusChanged = useCallback((
    boosterId: string, 
    status: ExecutionStatus, 
    progress?: number, 
    error?: string
  ) => {
    console.log(`[useBoosterExecution] Execution state changed:`, { boosterId, status, progress, error });
    updateExecution({ boosterId, status, progress, error });
  }, [updateExecution]);

  const onSyncRequired = useCallback(async () => {
    if (syncDebounceRef.current) {
      clearTimeout(syncDebounceRef.current);
    }
    
    syncDebounceRef.current = setTimeout(async () => {
      try {
        await serviceRef.current?.syncWithBackend();
      } catch (error) {
        console.error('[useBoosterExecution] Sync failed:', error);
      }
    }, 300);
  }, []);

  const onBatchStarted = useCallback((batchId: string) => {
    console.log(`[useBoosterExecution] Batch started: ${batchId}`);
    const batch: ExecutionBatch = {
      id: batchId,
      name: `Batch ${new Date().toLocaleTimeString()}`,
      executions: Array.from(executions.values()),
      status: 'processing',
      totalProgress: 0,
      createdAt: new Date(),
      startedAt: new Date(),
    };
    setCurrentBatch(batch);
  }, [executions, setCurrentBatch]);

  const onBatchCompleted = useCallback((batchId: string) => {
    console.log(`[useBoosterExecution] Batch completed: ${batchId}`);
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'completed',
      totalProgress: 100,
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  const onBatchError = useCallback((batchId: string, error: string) => {
    console.log(`[useBoosterExecution] Batch error: ${batchId} - ${error}`);
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'error',
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  useEffect(() => {
    if (isInitializedRef.current) {
      return;
    }

    console.log('[useBoosterExecution] Initializing service...');

    const callbacks: ExecutionServiceCallbacks = {
      onExecutionStatusChanged,
      onSyncRequired,
      onBatchStarted,
      onBatchCompleted,
      onBatchError,
    };

    serviceRef.current = BoosterExecutionService.getInstance(callbacks);
    serviceRef.current.initialize();
    isInitializedRef.current = true;

    return () => {
      console.log('[useBoosterExecution] Cleaning up service...');
      
      if (syncDebounceRef.current) {
        clearTimeout(syncDebounceRef.current);
        syncDebounceRef.current = null;
      }
      
      if (serviceRef.current && isInitializedRef.current) {
        serviceRef.current.destroy();
        serviceRef.current = null;
        isInitializedRef.current = false;
      }
    };
  }, []); 

  // Actions
  const executeBatch = useCallback(async (operations?: StagedOperations): Promise<string> => {
    if (!serviceRef.current || !isInitializedRef.current) {
      throw new Error('Execution service not initialized');
    }

    const opsToExecute = operations || stagedOperations;
    
    if (Object.keys(opsToExecute).length === 0) {
      throw new Error('No operations to execute');
    }

    try {
      console.log('[useBoosterExecution] Executing batch with operations:', opsToExecute);
      
      // Add executions to store
      addExecutions(opsToExecute);
      
      // Execute on service
      const batchId = await serviceRef.current.executeBatch(opsToExecute);
      return batchId;
    } catch (error) {
      console.error('[useBoosterExecution] Failed to execute batch:', error);
      throw error;
    } finally {
      clearStaging();
    }
  }, [stagedOperations, addExecutions, clearStaging]);

  const executeStagedBatch = useCallback(async (): Promise<string | null> => {
    if (Object.keys(stagedOperations).length === 0) {
      return null;
    }
    return executeBatch();
  }, [stagedOperations, executeBatch]);

  const cancelExecution = useCallback((boosterIds: string[]): void => {
    boosterIds.forEach(boosterId => {
      const execution = executions.get(boosterId);
      if (execution?.canCancel) {
        updateExecution({ boosterId, status: 'cancelled' });
      }
    });
  }, [executions, updateExecution]);

  const cancelStagedExecutions = useCallback((): void => {
    const boosterIds = Object.keys(stagedOperations);
    if (boosterIds.length > 0) {
      cancelExecution(boosterIds);
    }
  }, [stagedOperations, cancelExecution]);

  const removeCompletedExecutions = useCallback((): void => {
    const completedIds = Array.from(executions.values())
      .filter(e => e.status === 'completed' || e.status === 'error' || e.status === 'cancelled')
      .map(e => e.boosterId);
      
    if (completedIds.length > 0) {
      removeExecutions(completedIds);
    }
  }, [executions, removeExecutions]);

  const syncWithBackend = useCallback(async (): Promise<void> => {
    if (!serviceRef.current || !isInitializedRef.current) {
      console.warn('[useBoosterExecution] Service not available for sync');
      return;
    }
    
    try {
      await serviceRef.current.syncWithBackend();
    } catch (error) {
      console.error('[useBoosterExecution] Sync with backend failed:', error);
    }
  }, []);

  const resetExecution = useCallback((): void => {
    clearExecutions();
    clearStaging();
  }, [clearExecutions, clearStaging]);

  return {
    // State
    executions: Array.from(executions.values()),
    executionsByStatus,
    currentBatch,
    stats,
    isExecuting,
    canExecute,
    stagedOperations,
    stagedCount: Object.keys(stagedOperations).length,

    // Actions
    stageOperation,
    stageBatchOperations,
    clearStaging,
    executeBatch,
    executeStagedBatch,
    cancelExecution,
    cancelStagedExecutions,
    removeCompletedExecutions,
    syncWithBackend,
    resetExecution,
  };
};