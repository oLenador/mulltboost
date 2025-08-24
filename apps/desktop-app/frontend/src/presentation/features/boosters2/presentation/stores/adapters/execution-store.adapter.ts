// src/presentation/stores/adapters/execution-store.adapter.ts
import { PrimitiveAtom } from 'jotai';
import { ExecutionEntity } from '../../../core/entities/execution.entity';
import { BatchEntity } from '../../../core/entities/batch.entity';

export interface ExecutionStoreAdapter {
  updateExecution: (boosterId: string, status: string, progress?: number, error?: string) => void;
  addExecutions: (operations: Record<string, 'apply' | 'revert'>) => void;
  setBatch: (batch: BatchEntity) => void;
  clearBatch: () => void;
  removeExecution: (boosterId: string) => void;
  clearAllExecutions: () => void;
}

export interface ExecutionStoreAtoms {
  setExecutionAtom: PrimitiveAtom<(execution: Partial<ExecutionEntity> & { boosterId: string }) => void>;
  updateExecutionAtom: PrimitiveAtom<(boosterId: string, updates: Partial<ExecutionEntity>) => void>;
  removeExecutionAtom: PrimitiveAtom<(boosterId: string) => void>;
  clearExecutionsAtom: PrimitiveAtom<() => void>;
  setCurrentBatchAtom: PrimitiveAtom<(batch: BatchEntity | undefined) => void>;
}

export function createExecutionStoreAdapter(
  atoms: ExecutionStoreAtoms,
  getAtomValue: <T>(atom: PrimitiveAtom<T>) => T,
  setAtomValue: <T>(atom: PrimitiveAtom<T>, value: T) => void
): ExecutionStoreAdapter {
  
  const updateExecution = (boosterId: string, status: string, progress?: number, error?: string) => {
    const updateFn = getAtomValue(atoms.updateExecutionAtom);
    
    const updates: Partial<ExecutionEntity> = {
      status: status as any,
      progress: progress || 0,
      error,
      startedAt: status === 'processing' ? new Date() : undefined,
      completedAt: ['completed', 'error', 'cancelled'].includes(status) ? new Date() : undefined,
    };

    updateFn(boosterId, updates);
  };

  const addExecutions = (operations: Record<string, 'apply' | 'revert'>) => {
    const setFn = getAtomValue(atoms.setExecutionAtom);
    
    Object.entries(operations).forEach(([boosterId, operation]) => {
      setFn({
        boosterId,
        operation,
        status: 'queued',
        progress: 0,
        canCancel: true,
      });
    });
  };

  const setBatch = (batch: BatchEntity) => {
    const setBatchFn = getAtomValue(atoms.setCurrentBatchAtom);
    setBatchFn(batch);
  };

  const clearBatch = () => {
    const setBatchFn = getAtomValue(atoms.setCurrentBatchAtom);
    setBatchFn(undefined);
  };

  const removeExecution = (boosterId: string) => {
    const removeFn = getAtomValue(atoms.removeExecutionAtom);
    removeFn(boosterId);
  };

  const clearAllExecutions = () => {
    const clearFn = getAtomValue(atoms.clearExecutionsAtom);
    clearFn();
  };

  return {
    updateExecution,
    addExecutions,
    setBatch,
    clearBatch,
    removeExecution,
    clearAllExecutions,
  };
}