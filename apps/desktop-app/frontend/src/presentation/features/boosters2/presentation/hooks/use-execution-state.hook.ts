import { useCallback } from 'react';
import { useAtom } from 'jotai';
import { ExecutionEntity } from '../../core/entities/execution.entity';
import {
  executionsAtom,
  currentBatchAtom,
  executionsByStatusAtom,
  executionStatsAtom,
  isExecutingAtom,
  setExecutionAtom,
  updateExecutionAtom,
  removeExecutionAtom,
  removeExecutionsAtom,
  clearExecutionsAtom,
} from '../stores/execution-state.store';

export const useExecutionState = () => {
  const [executions] = useAtom(executionsAtom);
  const [currentBatch, setCurrentBatch] = useAtom(currentBatchAtom);
  const [executionsByStatus] = useAtom(executionsByStatusAtom);
  const [stats] = useAtom(executionStatsAtom);
  const [isExecuting] = useAtom(isExecutingAtom);

  // Actions
  const [, setExecution] = useAtom(setExecutionAtom);
  const [, updateExecution] = useAtom(updateExecutionAtom);
  const [, removeExecution] = useAtom(removeExecutionAtom);
  const [, removeExecutions] = useAtom(removeExecutionsAtom);
  const [, clearExecutions] = useAtom(clearExecutionsAtom);

  const getExecution = useCallback((boosterId: string): ExecutionEntity | undefined => {
    return executions.get(boosterId);
  }, [executions]);

  const getAllExecutions = useCallback((): ExecutionEntity[] => {
    return Array.from(executions.values());
  }, [executions]);

  const getExecutionsByStatus = useCallback((status: string): ExecutionEntity[] => {
    return executionsByStatus[status] || [];
  }, [executionsByStatus]);

  const addExecutions = useCallback((operations: Record<string, 'apply' | 'revert'>) => {
    Object.entries(operations).forEach(([boosterId, operation]) => {
      setExecution({
        boosterId,
        operation,
        status: 'queued',
        progress: 0,
        canCancel: true,
      });
    });
  }, [setExecution]);

  return {
    executions: getAllExecutions(),
    executionsByStatus,
    currentBatch,
    stats,
    isExecuting,
    getExecution,
    getAllExecutions,
    getExecutionsByStatus,
    addExecutions,
    setExecution,
    updateExecution,
    removeExecution,
    removeExecutions,
    clearExecutions,
    setCurrentBatch,
  };
};