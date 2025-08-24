import { QueueStateData, QueueItem } from '../types/api.types';
import { ExecutionEntity } from '../entities/execution.entity';

export interface SyncService {
  processQueueState(queueState: QueueStateData): ExecutionEntity[];
  determineProcessingExecutions(items: QueueItem[], inProgressCount: number): Set<string>;
}

export function createSyncService(): SyncService {
  return {
    processQueueState(queueState: QueueStateData): ExecutionEntity[] {
      const { items, inProgress } = queueState;
      const processingSet = determineProcessingExecutions(items, inProgress);

      return items
        .filter(item => item.boosterId)
        .map(item => ({
          boosterId: item.boosterId,
          operation: 'apply' as const, // We don't know from queue state, assume apply
          status: processingSet.has(item.boosterId) ? 'processing' as const : 'queued' as const,
          progress: typeof item.progress === 'number' ? item.progress : 0,
          error: item.error,
          canCancel: true,
        }));
    },

    determineProcessingExecutions(items: QueueItem[], inProgressCount: number): Set<string> {
      const processingSet = new Set<string>();
      
      for (let i = 0; i < Math.min(inProgressCount, items.length); i++) {
        const boosterId = items[i]?.boosterId;
        if (boosterId) {
          processingSet.add(boosterId);
        }
      }
      
      return processingSet;
    },
  };
}

function determineProcessingExecutions(items: QueueItem[], inProgressCount: number): Set<string> {
  const processingSet = new Set<string>();
  
  for (let i = 0; i < Math.min(inProgressCount, items.length); i++) {
    const boosterId = items[i]?.boosterId;
    if (boosterId) {
      processingSet.add(boosterId);
    }
  }
  
  return processingSet;
}