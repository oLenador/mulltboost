import { BoosterOperationType } from "bindings/github.com/oLenador/mulltbost/internal/core/domain/entities";

export type ItemStatus = 'idle' | 'queued' | 'processing' | 'completed' | 'error' | 'cancelled';

export interface ProcessableItem {
  id: string;
  name: string;
}

export interface ProcessableItemWithStatus<T extends ProcessableItem = ProcessableItem> {
  item: T;
  status: ItemStatus;
  progress?: number;
  error?: string;
  startedAt?: Date;
  completedAt?: Date;
  canCancel: boolean;
  operation: BoosterOperationType
}

type BatchOperationItems = Record<string, BoosterOperationType>

export interface BatchOperation {
  id: string;
  name: string;
  description?: string;
  items: ProcessableItemWithStatus;
  status: ItemStatus;
  progress: number;
  createdAt: Date;
  startedAt?: Date;
  completedAt?: Date;
}

export interface TabConfig {
  id: string;
  label: string;
  filter: (item: ProcessableItemWithStatus) => boolean;
  sortBy?: (a: ProcessableItemWithStatus, b: ProcessableItemWithStatus) => number;
}

export interface BatchManagerEvents<T extends ProcessableItem> {
  onItemStatusChanged: (item: ProcessableItemWithStatus<T>) => void;
  onBatchStarted: (operation: BatchOperation) => void;
  onBatchCompleted: (operation: BatchOperation) => void;
  onBatchError: (operation: BatchOperation, error: string) => void;
}