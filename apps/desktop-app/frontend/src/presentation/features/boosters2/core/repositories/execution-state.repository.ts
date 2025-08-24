import { BatchEntity } from '../entities/batch.entity';
import { ExecutionEntity} from '../entities/execution.entity';
import { ExecutionStats, StagedOperations } from '../types/execution.types';

export interface ExecutionStateRepository {
  // Executions
  getExecution(boosterId: string): ExecutionEntity | undefined;
  getAllExecutions(): ExecutionEntity[];
  setExecution(execution: ExecutionEntity): void;
  updateExecution(boosterId: string, updates: Partial<ExecutionEntity>): void;
  removeExecution(boosterId: string): void;
  removeExecutions(boosterIds: string[]): void;
  clearExecutions(): void;

  // Batches
  getCurrentBatch(): BatchEntity | undefined;
  setCurrentBatch(batch: BatchEntity | undefined): void;

  // Stats
  getExecutionStats(): ExecutionStats;

  // Staging
  getStagedOperations(): StagedOperations;
  setStagedOperations(operations: StagedOperations): void;
  addStagedOperation(boosterId: string, operation: 'apply' | 'revert'): void;
  removeStagedOperation(boosterId: string): void;
  clearStagedOperations(): void;
}