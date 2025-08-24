import { atom } from 'jotai';
import { ExecutionEntity } from '../../core/entities/execution.entity';
import { ExecutionStats, createExecutionStats } from '../../core/types/execution.types';
import { BatchEntity } from '../../core/entities/batch.entity';

export const executionsAtom = atom<Map<string, ExecutionEntity>>(new Map());
export const currentBatchAtom = atom<BatchEntity | undefined>(undefined);

// Computed atoms
export const executionsByStatusAtom = atom((get) => {
  const executions = get(executionsAtom);
  const byStatus: Record<string, ExecutionEntity[]> = {};
  
  Array.from(executions.values()).forEach(execution => {
    if (!byStatus[execution.status]) {
      byStatus[execution.status] = [];
    }
    byStatus[execution.status].push(execution);
  });
  
  return byStatus;
});

export const executionStatsAtom = atom((get) => {
  const executions = get(executionsAtom);
  const allExecutions = Array.from(executions.values());
  
  const total = allExecutions.length;
  const queued = allExecutions.filter(e => e.status === 'queued').length;
  const processing = allExecutions.filter(e => e.status === 'processing').length;
  const completed = allExecutions.filter(e => e.status === 'completed').length;
  const errors = allExecutions.filter(e => e.status === 'error').length;
  const cancelled = allExecutions.filter(e => e.status === 'cancelled').length;

  return createExecutionStats(total, queued, processing, completed, errors, cancelled);
});

export const isExecutingAtom = atom((get) => {
  const stats = get(executionStatsAtom);
  return stats.queued > 0 || stats.processing > 0;
});

// Actions
export const setExecutionAtom = atom(
  null,
  (get, set, execution: ExecutionEntity) => {
    const executions = get(executionsAtom);
    const newExecutions = new Map(executions);
    newExecutions.set(execution.boosterId, execution);
    set(executionsAtom, newExecutions);
  }
);

export const updateExecutionAtom = atom(
  null,
  (get, set, { boosterId, updates }: { boosterId: string; updates: Partial<ExecutionEntity> }) => {
    const executions = get(executionsAtom);
    const existing = executions.get(boosterId);
    
    if (existing) {
      const updated = { ...existing, ...updates };
      const newExecutions = new Map(executions);
      newExecutions.set(boosterId, updated);
      set(executionsAtom, newExecutions);
    }
  }
);

export const removeExecutionAtom = atom(
  null,
  (get, set, boosterId: string) => {
    const executions = get(executionsAtom);
    const newExecutions = new Map(executions);
    newExecutions.delete(boosterId);
    set(executionsAtom, newExecutions);
  }
);

export const removeExecutionsAtom = atom(
  null,
  (get, set, boosterIds: string[]) => {
    const executions = get(executionsAtom);
    const newExecutions = new Map(executions);
    boosterIds.forEach(id => newExecutions.delete(id));
    set(executionsAtom, newExecutions);
  }
);

export const clearExecutionsAtom = atom(
  null,
  (get, set) => {
    set(executionsAtom, new Map());
    set(currentBatchAtom, undefined);
  }
);
