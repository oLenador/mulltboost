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
} from '@/presentation/features/boosters/stores/booster-execution.store';
import { ExecutionServiceCallbacks } from '../domain/booster-execution.service';
import { ExecutionBatch, ExecutionStatus } from '../domain/booster-queue.types';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { globalBoosterService } from '../domain/booster-global.service';

let hookInstanceCounter = 0;

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

  // Unique subscriber ID para o service global
  const subscriberIdRef = useRef(`hook-${++hookInstanceCounter}-${Date.now()}`);
  const syncDebounceRef = useRef<NodeJS.Timeout | null>(null);
  const isSubscribedRef = useRef(false);

  // Callbacks estáveis usando useCallback
  const onExecutionStatusChanged = useCallback((
    boosterId: string, 
    status: ExecutionStatus, 
    progress?: number, 
    error?: string
  ) => {
    console.log(`[useBoosterExecution-${subscriberIdRef.current}] Execution state changed:`, { boosterId, status, progress, error });
    updateExecution({ boosterId, status, progress, error });
  }, [updateExecution]);

  const onSyncRequired = useCallback(async () => {
    // Debounce para evitar múltiplas sincronizações
    if (syncDebounceRef.current) {
      clearTimeout(syncDebounceRef.current);
    }
    
    syncDebounceRef.current = setTimeout(async () => {
      try {
        await globalBoosterService.syncWithBackend();
      } catch (error) {
        console.error(`[useBoosterExecution-${subscriberIdRef.current}] Sync failed:`, error);
      }
    }, 300);
  }, []);

  const onBatchStarted = useCallback((batchId: string) => {
    console.log(`[useBoosterExecution-${subscriberIdRef.current}] Batch started: ${batchId}`);
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
    console.log(`[useBoosterExecution-${subscriberIdRef.current}] Batch completed: ${batchId}`);
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'completed',
      totalProgress: 100,
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  const onBatchError = useCallback((batchId: string, error: string) => {
    console.log(`[useBoosterExecution-${subscriberIdRef.current}] Batch error: ${batchId} - ${error}`);
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'error',
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  // Subscribe to global service
  useEffect(() => {
    if (isSubscribedRef.current) {
      return;
    }

    const subscriberId = subscriberIdRef.current;
    console.log(`[useBoosterExecution] Subscribing to global service: ${subscriberId}`);

    const callbacks: ExecutionServiceCallbacks = {
      onExecutionStatusChanged,
      onSyncRequired,
      onBatchStarted,
      onBatchCompleted,
      onBatchError,
    };

    globalBoosterService.subscribe(subscriberId, callbacks);
    isSubscribedRef.current = true;

    // Cleanup
    return () => {
      console.log(`[useBoosterExecution] Unsubscribing from global service: ${subscriberId}`);
      
      if (syncDebounceRef.current) {
        clearTimeout(syncDebounceRef.current);
        syncDebounceRef.current = null;
      }
      
      if (isSubscribedRef.current) {
        globalBoosterService.unsubscribe(subscriberId);
        isSubscribedRef.current = false;
      }
    };
  }, []); // Empty dependency array - só subscribe uma vez

  // Actions
  const executeBatch = useCallback(async (operations?: StagedOperations): Promise<string> => {
    if (!globalBoosterService.isRunning()) {
      throw new Error('Global booster service not running');
    }

    const opsToExecute = operations || stagedOperations;
    
    if (Object.keys(opsToExecute).length === 0) {
      throw new Error('No operations to execute');
    }

    try {
      console.log(`[useBoosterExecution-${subscriberIdRef.current}] Executing batch with operations:`, opsToExecute);
      
      // Add executions to store
      addExecutions(opsToExecute);
      
      // Execute via global service
      const batchId = await globalBoosterService.executeBatch(opsToExecute);
      return batchId;
    } catch (error) {
      console.error(`[useBoosterExecution-${subscriberIdRef.current}] Failed to execute batch:`, error);
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
    if (!globalBoosterService.isRunning()) {
      console.warn(`[useBoosterExecution-${subscriberIdRef.current}] Global service not available for sync`);
      return;
    }
    
    try {
      await globalBoosterService.syncWithBackend();
    } catch (error) {
      console.error(`[useBoosterExecution-${subscriberIdRef.current}] Sync with backend failed:`, error);
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