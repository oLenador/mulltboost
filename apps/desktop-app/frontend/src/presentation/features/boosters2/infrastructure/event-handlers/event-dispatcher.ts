import { EventData } from '../../core/types/events.types';
import { ExecutionEntity } from '../../core/entities/execution.entity';

export interface EventCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: string, progress?: number, error?: string) => void;
  onSyncRequired: () => void | Promise<void>;
}

export interface EventDispatcher {
  dispatch(eventData: EventData, executionUpdate: Partial<ExecutionEntity>): void;
  setCallbacks(callbacks: EventCallbacks): void;
}

export function createEventDispatcher(): EventDispatcher {
  let callbacks: EventCallbacks | null = null;

  return {
    dispatch(eventData: EventData, executionUpdate: Partial<ExecutionEntity>): void {
      if (!callbacks) return;

      try {
        // Dispatch execution status change
        if (executionUpdate.status) {
          callbacks.onExecutionStatusChanged(
            eventData.boosterId,
            executionUpdate.status,
            executionUpdate.progress,
            executionUpdate.error
          );
        }

        // Dispatch sync requirement for batch events
        if (eventData.eventType.includes('batch')) {
          callbacks.onSyncRequired();
        }
      } catch (error) {
        console.error('Error dispatching event:', error);
      }
    },

    setCallbacks(newCallbacks: EventCallbacks): void {
      callbacks = newCallbacks;
    },
  };
}
