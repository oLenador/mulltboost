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
import { globalBoosterService } from '../domain/booster-global.service';

let hookInstanceCounter = 0;

export const useBoosterExecution = () => {
  const [executions] = useAtom(boosterExecutionsAtom);
  const [currentBatch, setCurrentBatch] = useAtom(currentBatchAtom);
  const [stats] = useAtom(executionStatsAtom);
  const [executionsByStatus] = useAtom(executionsByStatusAtom);
  const [isExecuting] = useAtom(isExecutingAtom);
  const [canExecute] = useAtom(canExecuteBatchAtom);
  const [stagedOperations] = useAtom(stagedOperationsAtom);

  const [, updateExecution] = useAtom(updateExecutionAtom);
  const [, addExecutions] = useAtom(addExecutionsAtom);
  const [, removeExecutions] = useAtom(removeExecutionsAtom);
  const [, clearExecutions] = useAtom(clearExecutionsAtom);
  const [, stageOperation] = useAtom(stageOperationAtom);
  const [, stageBatchOperations] = useAtom(stageBatchOperationsAtom);
  const [, clearStaging] = useAtom(clearStagingAtom);

  const subscriberIdRef = useRef(`hook-${++hookInstanceCounter}-${Date.now()}`);
  const syncDebounceRef = useRef<NodeJS.Timeout | null>(null);
  const isSubscribedRef = useRef(false);

  const onExecutionStatusChanged = useCallback((
    boosterId: string, 
    status: ExecutionStatus, 
    progress?: number, 
    error?: string
  ) => {
    updateExecution({ boosterId, status, progress, error });
  }, [updateExecution]);

  const onSyncRequired = useCallback(async () => {
    if (syncDebounceRef.current) {
      clearTimeout(syncDebounceRef.current);
    }
    
    syncDebounceRef.current = setTimeout(async () => {
      try {
        await globalBoosterService.syncWithBackend();
      } catch (error) {
        console.error(`Sync failed:`, error);
      }
    }, 300);
  }, []);

  const onBatchStarted = useCallback((batchId: string) => {
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
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'completed',
      totalProgress: 100,
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  const onBatchError = useCallback((batchId: string, error: string) => {
    setCurrentBatch(prev => prev ? {
      ...prev,
      status: 'error',
      completedAt: new Date(),
    } : null);
  }, [setCurrentBatch]);

  useEffect(() => {
    if (isSubscribedRef.current) return;

    const subscriberId = subscriberIdRef.current;
    const callbacks: ExecutionServiceCallbacks = {
      onExecutionStatusChanged,
      onSyncRequired,
      onBatchStarted,
      onBatchCompleted,
      onBatchError,
    };

    globalBoosterService.subscribe(subscriberId, callbacks);
    isSubscribedRef.current = true;

    return () => {
      if (syncDebounceRef.current) {
        clearTimeout(syncDebounceRef.current);
        syncDebounceRef.current = null;
      }
      
      if (isSubscribedRef.current) {
        globalBoosterService.unsubscribe(subscriberId);
        isSubscribedRef.current = false;
      }
    };
  }, []);

  const executeBatch = useCallback(async (operations?: StagedOperations): Promise<string> => {
    if (!globalBoosterService.isRunning()) {
      throw new Error('Global booster service not running');
    }

    const opsToExecute = operations || stagedOperations;
    
    if (Object.keys(opsToExecute).length === 0) {
      throw new Error('No operations to execute');
    }

    try {
      addExecutions(opsToExecute);
      const batchId = await globalBoosterService.executeBatch(opsToExecute);
      return batchId;
    } catch (error) {
      throw error;
    } finally {
      clearStaging();
    }
  }, [stagedOperations, addExecutions, clearStaging]);

  const executeStagedBatch = useCallback(async (): Promise<string | null> => {
    return Object.keys(stagedOperations).length === 0 ? null : executeBatch();
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
      .filter(e => ['completed', 'error', 'cancelled'].includes(e.status))
      .map(e => e.boosterId);
      
    if (completedIds.length > 0) {
      removeExecutions(completedIds);
    }
  }, [executions, removeExecutions]);

  const syncWithBackend = useCallback(async (): Promise<void> => {
    if (!globalBoosterService.isRunning()) return;
    
    try {
      await globalBoosterService.syncWithBackend();
    } catch (error) {
      console.error('Sync with backend failed:', error);
    }
  }, []);

  const resetExecution = useCallback((): void => {
    clearExecutions();
    clearStaging();
  }, [clearExecutions, clearStaging]);

  return {
    executions: Array.from(executions.values()),
    executionsByStatus,
    currentBatch,
    stats,
    isExecuting,
    canExecute,
    stagedOperations,
    stagedCount: Object.keys(stagedOperations).length,
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