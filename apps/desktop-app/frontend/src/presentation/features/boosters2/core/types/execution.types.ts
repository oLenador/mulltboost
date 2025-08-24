// src/presentation/features/boosters/core/types/events.types.ts
export enum EventType {
    Processing = "booster.processing",
    Success = "booster.success",
    Error = "booster.error",
    Failed = "booster.failed",
    Queued = "booster.queued",
    BatchQueued = "booster.batch_queued",
    Cancelled = "booster.cancelled",
  }
  
  export interface EventData {
    readonly eventType: EventType;
    readonly timestamp: Date;
    readonly operationType: string;
    readonly operationId: string;
    readonly boosterId: string;
    readonly sequence: number;
    readonly idempotencyId: string;
    readonly status: string;
    readonly endAt?: Date;
    readonly error?: string;
    readonly queueSize: number;
  }
  
  export function createEventData(rawData: any): EventData {
    return {
      eventType: rawData.EventType as EventType,
      timestamp: new Date(rawData.Timestamp),
      operationType: rawData.OperationType,
      operationId: rawData.OperationID,
      boosterId: rawData.BoosterID,
      sequence: rawData.Sequency || 0,
      idempotencyId: rawData.IdempotencyID,
      status: rawData.Status,
      endAt: rawData.EndAt ? new Date(rawData.EndAt) : undefined,
      error: rawData.Error,
      queueSize: rawData.QueueSize || 0,
    };
  }
  
  // src/presentation/features/boosters/core/types/execution.types.ts
  export interface ExecutionStats {
    readonly total: number;
    readonly queued: number;
    readonly processing: number;
    readonly completed: number;
    readonly errors: number;
    readonly cancelled: number;
    readonly totalProgress: number;
  }
  
  export function createExecutionStats(
    total: number = 0,
    queued: number = 0,
    processing: number = 0,
    completed: number = 0,
    errors: number = 0,
    cancelled: number = 0
  ): ExecutionStats {
    const totalProgress = total > 0 ? Math.round((completed / total) * 100) : 0;
    
    return {
      total,
      queued,
      processing,
      completed,
      errors,
      cancelled,
      totalProgress,
    };
  }
  
  export interface StagedOperations {
    readonly [boosterId: string]: 'apply' | 'revert';
  }
  
  export function createStagedOperations(operations: Record<string, 'apply' | 'revert'> = {}): StagedOperations {
    return { ...operations };
  }
  
  export function addStagedOperation(
    staged: StagedOperations,
    boosterId: string,
    operation: 'apply' | 'revert'
  ): StagedOperations {
    return {
      ...staged,
      [boosterId]: operation,
    };
  }
  
  export function removeStagedOperation(staged: StagedOperations, boosterId: string): StagedOperations {
    const { [boosterId]: removed, ...rest } = staged;
    return rest;
  }
  