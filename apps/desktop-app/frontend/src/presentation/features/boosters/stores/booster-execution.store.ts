// src/core/store/booster-execution.store.ts

import { atom } from 'jotai';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { BoosterExecution, ExecutionStatus, ExecutionBatch, ExecutionStats } from '../domain/booster-queue.types';
import { BoosterItem } from '../types/booster.types';

// Base atoms
export const boosterExecutionsAtom = atom<Map<string, BoosterExecution>>(new Map());
export const currentBatchAtom = atom<ExecutionBatch | null>(null);
export const listingBoostersAtom = atom<BoosterItem[]>([]);

// Staged operations (user selections before execution)
export type StagedOperations = Record<string, BoosterOperationType>;
export const stagedOperationsAtom = atom<StagedOperations>({});

// Derived atoms
export const executionStatsAtom = atom((get) => {
  const executions = Array.from(get(boosterExecutionsAtom).values());
  
  const stats: ExecutionStats = {
    total: executions.length,
    queued: executions.filter(e => e.status === 'queued').length,
    processing: executions.filter(e => e.status === 'processing').length,
    completed: executions.filter(e => !!e.completedAt).length,
    errors: executions.filter(e => e.status === 'error').length,
    cancelled: executions.filter(e => e.status === 'cancelled').length,
    totalProgress: executions.length > 0 
      ? executions.reduce((sum, e) => sum + e.progress, 0) / executions.length 
      : 0
  };
  
  return stats;
});

export const executionsByStatusAtom = atom((get) => {
  const executions = Array.from(get(boosterExecutionsAtom).values());
  
  return {
    idle: executions.filter(e => e.status === 'idle'),
    queued: executions.filter(e => e.status === 'queued'),
    processing: executions.filter(e => e.status === 'processing'),
    success: executions.filter(e => !!e.completedAt),
    error: executions.filter(e => e.status === 'error'),
    cancelled: executions.filter(e => e.status === 'cancelled'),
    completed: executions.filter(e => e.status === 'completed'),

  };
});

export const isExecutingAtom = atom((get) => {
  const stats = get(executionStatsAtom);
  return stats.queued > 0 || stats.processing > 0;
});

export const canExecuteBatchAtom = atom((get) => {
  const staged = get(stagedOperationsAtom);
  const isExecuting = get(isExecutingAtom);
  return Object.keys(staged).length > 0 && !isExecuting;
});

// Actions
export const updateExecutionAtom = atom(
  null,
  (
    get,
    set,
    update: {
      boosterId: string;
      status?: ExecutionStatus;
      progress?: number;
      error?: string;
    }
  ) => {
    console.log("UpdateExectutionAtom: ", update)
    const executions = new Map(get(boosterExecutionsAtom));
    const existing = executions.get(update.boosterId);

    if (!existing) return;

    const now = new Date();

    const updated: BoosterExecution = {
      ...existing,

      // Update status if provided
      ...(update.status && { status: update.status }),

      // Update progress if provided
      ...(update.progress !== undefined && { progress: update.progress }),

      // Update error if provided
      ...(update.error && { error: update.error }),

      // Determine if can cancel
      canCancel: update.status === 'processing' || update.status === 'queued',

      // Set startedAt if processing and not already started
      ...(update.status === 'processing' && !existing.startedAt
        ? { startedAt: now }
        : {}),

      // Set completedAt and disable cancel if status is final
      ...((['completed', 'error', 'cancelled'] as ExecutionStatus[]).includes(
        update.status!
      )
        ? { completedAt: now, canCancel: false }
        : {}),
    };
    console.log("UpdateExectutionAtom, updated item: ", updated)

    executions.set(update.boosterId, updated);
    console.log("UpdateExectutionAtom, executions: ", executions)

    set(boosterExecutionsAtom, executions);
  }
);

export const addExecutionsAtom = atom(
  null,
  (get, set, operations: StagedOperations) => {
    const executions = new Map(get(boosterExecutionsAtom));
    
    Object.entries(operations).forEach(([boosterId, operation]) => {
      const execution: BoosterExecution = {
        boosterId,
        operation,
        status: 'idle',
        progress: 0,
        canCancel: false,
      };
      executions.set(boosterId, execution);
    });
    
    set(boosterExecutionsAtom, executions);
  }
);

export const removeExecutionsAtom = atom(
  null,
  (get, set, boosterIds: string[]) => {
    const executions = new Map(get(boosterExecutionsAtom));
    
    boosterIds.forEach(id => {
      const execution = executions.get(id);
      if (execution && execution.status === 'idle') {
        executions.delete(id);
      }
    });
    
    set(boosterExecutionsAtom, executions);
  }
);

export const clearExecutionsAtom = atom(
  null,
  (get, set) => {
    set(boosterExecutionsAtom, new Map());
    set(currentBatchAtom, null);
  }
);

export const stageOperationAtom = atom(
  null,
  (get, set, { boosterId, operation }: { boosterId: string; operation: BoosterOperationType }) => {
    const staged = get(stagedOperationsAtom);
    
    if (staged[boosterId] === operation) {
      // Remove if same operation
      const { [boosterId]: _, ...rest } = staged;
      set(stagedOperationsAtom, rest);
    } else {
      // Add or update operation
      set(stagedOperationsAtom, { ...staged, [boosterId]: operation });
    }
  }
);

export const stageBatchOperationsAtom = atom(
  null,
  (get, set, operations: StagedOperations) => {
    const staged = get(stagedOperationsAtom);
    set(stagedOperationsAtom, { ...staged, ...operations });
  }
);

export const clearStagingAtom = atom(
  null,
  (get, set) => {
    set(stagedOperationsAtom, {});
  }
);