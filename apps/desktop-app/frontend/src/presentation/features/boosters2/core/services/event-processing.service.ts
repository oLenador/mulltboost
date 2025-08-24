import { EventData, EventType } from '../types/events.types';
import { ExecutionEntity } from '../entities/execution.entity';

export interface EventCallback {
  (boosterId: string, status: string, progress?: number, error?: string): void;
}

export interface SyncCallback {
  (): void | Promise<void>;
}

export interface EventProcessingService {
  processEvent(eventData: EventData): ProcessedEvent | null;
  shouldProcessEvent(eventData: EventData, lastSequence: number, processedIds: Set<string>): boolean;
  convertEventToExecution(eventData: EventData): Partial<ExecutionEntity> | null;
}

export interface ProcessedEvent {
  readonly eventData: EventData;
  readonly executionUpdate: Partial<ExecutionEntity>;
  readonly shouldTriggerSync: boolean;
}

export function createEventProcessingService(): EventProcessingService {
  return {
    processEvent(eventData: EventData): ProcessedEvent | null {
      const executionUpdate = convertEventToExecution(eventData);
      if (!executionUpdate) return null;

      return {
        eventData,
        executionUpdate,
        shouldTriggerSync: isBatchEvent(eventData),
      };
    },

    shouldProcessEvent(eventData: EventData, lastSequence: number, processedIds: Set<string>): boolean {
      // Check idempotency
      if (processedIds.has(eventData.idempotencyId)) {
        return false;
      }

      // Check sequence
      if (eventData.sequence > 0 && eventData.sequence <= lastSequence) {
        return false;
      }

      return true;
    },

    convertEventToExecution(eventData: EventData): Partial<ExecutionEntity> | null {
      return convertEventToExecution(eventData);
    },
  };
}

function convertEventToExecution(eventData: EventData): Partial<ExecutionEntity> | null {
  const baseUpdate = {
    boosterId: eventData.boosterId,
  };

  switch (eventData.eventType) {
    case EventType.Queued:
      return { ...baseUpdate, status: 'queued' as const };
    
    case EventType.Processing:
      return { ...baseUpdate, status: 'processing' as const, startedAt: eventData.timestamp };
    
    case EventType.Success:
      return { 
        ...baseUpdate, 
        status: 'completed' as const, 
        progress: 100, 
        completedAt: eventData.endAt || eventData.timestamp 
      };
    
    case EventType.Error:
    case EventType.Failed:
      return { 
        ...baseUpdate, 
        status: 'error' as const, 
        error: eventData.error,
        completedAt: eventData.endAt || eventData.timestamp 
      };
    
    case EventType.Cancelled:
      return { 
        ...baseUpdate, 
        status: 'cancelled' as const,
        completedAt: eventData.endAt || eventData.timestamp 
      };
    
    default:
      return null;
  }
}

function isBatchEvent(eventData: EventData): boolean {
  return eventData.eventType === EventType.BatchQueued;
}