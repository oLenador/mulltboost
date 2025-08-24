import { ExecutionEntity, createExecutionEntity } from '../entities/execution.entity';
import { BatchEntity, createBatchEntity } from '../entities/batch.entity';
import { StagedOperations } from '../types/execution.types';

export interface ExecutionService {
  createExecutionsFromOperations(operations: StagedOperations): ExecutionEntity[];
  createBatchFromExecutions(executions: ExecutionEntity[]): BatchEntity;
  canExecuteOperations(operations: StagedOperations): { canExecute: boolean; reason?: string };
  generateBatchId(): string;
  calculateBatchProgress(executions: ExecutionEntity[]): number;
}

export function createExecutionService(): ExecutionService {
  return {
    createExecutionsFromOperations(operations: StagedOperations): ExecutionEntity[] {
      return Object.entries(operations).map(([boosterId, operation]) =>
        createExecutionEntity(boosterId, operation, { status: 'queued' })
      );
    },

    createBatchFromExecutions(executions: ExecutionEntity[]): BatchEntity {
      const batchId = generateBatchId();
      const executionIds = executions.map(e => e.boosterId);
      
      return createBatchEntity(
        batchId,
        `Batch ${new Date().toLocaleTimeString()}`,
        executionIds,
        { status: 'queued' }
      );
    },

    canExecuteOperations(operations: StagedOperations): { canExecute: boolean; reason?: string } {
      const operationCount = Object.keys(operations).length;
      
      if (operationCount === 0) {
        return { canExecute: false, reason: 'No operations to execute' };
      }
      
      if (operationCount > 100) {
        return { canExecute: false, reason: 'Too many operations (max 100)' };
      }
      
      return { canExecute: true };
    },

    generateBatchId(): string {
      return generateBatchId();
    },

    calculateBatchProgress(executions: ExecutionEntity[]): number {
      if (executions.length === 0) return 0;
      
      const totalProgress = executions.reduce((sum, exec) => sum + exec.progress, 0);
      return Math.round(totalProgress / executions.length);
    },
  };
}

function generateBatchId(): string {
  return `batch_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}
