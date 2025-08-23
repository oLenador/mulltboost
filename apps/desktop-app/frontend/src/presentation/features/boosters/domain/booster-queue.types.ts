// src/presentation/features/boosters/domain/booster-queue.types.ts

import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { BoosterItem } from '../types/booster.types';

export type ExecutionStatus = 'idle' | 'queued' | 'processing' | 'completed' | 'error' | 'cancelled';

export interface BoosterExecution {
  boosterId: string;
  operation: BoosterOperationType;
  status: ExecutionStatus;
  progress: number;
  error?: string;
  startedAt?: Date;
  completedAt?: Date;
  canCancel: boolean;
}

export interface ExecutionBatch {
  id: string;
  name: string;
  executions: BoosterExecution[];
  status: ExecutionStatus;
  totalProgress: number;
  createdAt: Date;
  startedAt?: Date;
  completedAt?: Date;
}

export interface QueueState {
  items: Array<{
    BoosterID: string;
    OperationID?: string;
    Progress?: number;
    Error?: string;
  }>;
  InProgress: number;
}

export interface ExecutionStats {
  total: number;
  queued: number;
  processing: number;
  completed: number;
  errors: number;
  cancelled: number;
  totalProgress: number;
}

// Events from backend
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