import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';
import { ExecutionEntity } from '../../core/entities/execution.entity';
import { ExecutionStats, StagedOperations, createExecutionStats } from '../../core/types/execution.types';
import { BatchEntity } from '../../core/entities/batch.entity';

export function createInMemoryExecutionRepository(): ExecutionStateRepository {
  const executions = new Map<string, ExecutionEntity>();
  let currentBatch: BatchEntity | undefined;
  let stagedOperations: StagedOperations = {};

  return {
    // Executions
    getExecution(boosterId: string): ExecutionEntity | undefined {
      return executions.get(boosterId);
    },

    getAllExecutions(): ExecutionEntity[] {
      return Array.from(executions.values());
    },

    setExecution(execution: ExecutionEntity): void {
      executions.set(execution.boosterId, execution);
    },

    updateExecution(boosterId: string, updates: Partial<ExecutionEntity>): void {
      const existing = executions.get(boosterId);
      if (existing) {
        executions.set(boosterId, { ...existing, ...updates });
      }
    },

    removeExecution(boosterId: string): void {
      executions.delete(boosterId);
    },

    removeExecutions(boosterIds: string[]): void {
      boosterIds.forEach(id => executions.delete(id));
    },

    clearExecutions(): void {
      executions.clear();
    },

    // Batches
    getCurrentBatch(): BatchEntity | undefined {
      return currentBatch;
    },

    setCurrentBatch(batch: BatchEntity | undefined): void {
      currentBatch = batch;
    },

    // Stats
    getExecutionStats(): ExecutionStats {
      const allExecutions = Array.from(executions.values());
      const total = allExecutions.length;
      const queued = allExecutions.filter(e => e.status === 'queued').length;
      const processing = allExecutions.filter(e => e.status === 'processing').length;
      const completed = allExecutions.filter(e => e.status === 'completed').length;
      const errors = allExecutions.filter(e => e.status === 'error').length;
      const cancelled = allExecutions.filter(e => e.status === 'cancelled').length;

      return createExecutionStats(total, queued, processing, completed, errors, cancelled);
    },

    // Staging
    getStagedOperations(): StagedOperations {
      return { ...stagedOperations };
    },

    setStagedOperations(operations: StagedOperations): void {
      stagedOperations = { ...operations };
    },

    addStagedOperation(boosterId: string, operation: 'apply' | 'revert'): void {
      stagedOperations = { ...stagedOperations, [boosterId]: operation };
    },

    removeStagedOperation(boosterId: string): void {
      const { [boosterId]: removed, ...rest } = stagedOperations;
      stagedOperations = rest;
    },

    clearStagedOperations(): void {
      stagedOperations = {};
    },
  };
}