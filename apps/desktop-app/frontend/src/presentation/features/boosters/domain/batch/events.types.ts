
export enum BoosterEvent {
  Processing = "booster.processing",
  Success = "booster.success", 
  Error = "booster.error",
  Failed = "booster.failed",
  Queued = "booster.queued",
  BatchQueued = "booster.batch_queued",
  Cancelled = "booster.cancelled",
}

export interface BoosterEventData {
  EventType: string;
  Timestamp: string;
  OperationType: string;
  OperationID: string;
  BoosterID: string;
  Status: string;
  EndAt?: string;
  Error?: string;
  QueueSize: number;
}

export interface BoosterBatchProgressEventData {
  EventType: string;
  Timestamp: string;
  BatchID: string;
  OperationType: string;
  TotalCount: number;
  QueuedCount: number;
  ValidationErrors?: Record<string, string>;
  QueueSize: number;
}

export interface BackendBoosterDto {
  id: string;
  name: string;
  description: string;
  impact: 'low' | 'medium' | 'high';
  category: string;
  enabled: boolean;
  advanced?: boolean;
  requiresRestart?: boolean;
}

export interface QueueItem {
  id: string;
  boosterId: string;
  status: string;
  progress: number;
  error?: string;
  createdAt: string;
  startedAt?: string;
  completedAt?: string;
}

export interface InitResult {
  operationId: string;
  success: boolean;
  message?: string;
  error?: string;
}