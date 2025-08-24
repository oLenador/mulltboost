export type ExecutionStatus = 'idle' | 'queued' | 'processing' | 'completed' | 'error' | 'cancelled';
export type OperationType = 'apply' | 'revert';

export interface ExecutionEntity {
  readonly boosterId: string;
  readonly operation: OperationType;
  readonly status: ExecutionStatus;
  readonly progress: number;
  readonly error?: string;
  readonly startedAt?: Date;
  readonly completedAt?: Date;
  readonly canCancel: boolean;
}

export function createExecutionEntity(
  boosterId: string,
  operation: OperationType,
  overrides: Partial<ExecutionEntity> = {}
): ExecutionEntity {
  return {
    boosterId,
    operation,
    status: overrides.status || 'idle',
    progress: overrides.progress || 0,
    error: overrides.error,
    startedAt: overrides.startedAt,
    completedAt: overrides.completedAt,
    canCancel: overrides.canCancel ?? true,
  };
}

export function updateExecutionEntity(
  execution: ExecutionEntity,
  updates: Partial<Pick<ExecutionEntity, 'status' | 'progress' | 'error' | 'startedAt' | 'completedAt'>>
): ExecutionEntity {
  return {
    ...execution,
    ...updates,
    canCancel: updates.status ? updates.status === 'queued' || updates.status === 'processing' : execution.canCancel,
  };
}