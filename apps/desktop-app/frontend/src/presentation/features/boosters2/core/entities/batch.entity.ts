import { ExecutionStatus } from "./execution.entity";

  export interface BatchEntity {
    readonly id: string;
    readonly name: string;
    readonly executionIds: string[];
    readonly status: ExecutionStatus;
    readonly totalProgress: number;
    readonly createdAt: Date;
    readonly startedAt?: Date;
    readonly completedAt?: Date;
  }
  
  export function createBatchEntity(
    id: string,
    name: string,
    executionIds: string[],
    overrides: Partial<BatchEntity> = {}
  ): BatchEntity {
    return {
      id,
      name,
      executionIds,
      status: overrides.status || 'idle',
      totalProgress: overrides.totalProgress || 0,
      createdAt: overrides.createdAt || new Date(),
      startedAt: overrides.startedAt,
      completedAt: overrides.completedAt,
    };
  }
  
  export function updateBatchEntity(
    batch: BatchEntity,
    updates: Partial<Pick<BatchEntity, 'status' | 'totalProgress' | 'startedAt' | 'completedAt'>>
  ): BatchEntity {
    return {
      ...batch,
      ...updates,
    };
  }