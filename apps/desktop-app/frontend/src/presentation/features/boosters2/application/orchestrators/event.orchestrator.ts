import { EventListener } from '../../infrastructure/event-handlers/event-listener';
import { SequenceManager, SequenceState, createSequenceState } from '../../infrastructure/event-handlers/sequence-manager';
import { IdempotencyManager, IdempotencyState, createIdempotencyState } from '../../infrastructure/event-handlers/idempotency-manager';
import { EventDispatcher, EventCallbacks } from '../../infrastructure/event-handlers/event-dispatcher';
import { EventProcessingService } from '../../core/services/event-processing.service';
import { EventData } from '../../core/types/events.types';

export interface EventOrchestrator {
  start(): void;
  stop(): void;
  setCallbacks(callbacks: EventCallbacks): void;
  isRunning(): boolean;
  getStats(): {
    totalProcessed: number;
    totalPending: number;
    boostersTracked: number;
  };
}

export function createEventOrchestrator(
  eventListener: EventListener,
  sequenceManager: SequenceManager,
  idempotencyManager: IdempotencyManager,
  eventDispatcher: EventDispatcher,
  eventProcessingService: EventProcessingService
): EventOrchestrator {
  let isStarted = false;
  const boosterStates = new Map<string, { sequence: SequenceState; idempotency: IdempotencyState }>();
  let cleanupInterval: NodeJS.Timeout | undefined;

  const getOrCreateBoosterState = (boosterId: string) => {
    if (!boosterStates.has(boosterId)) {
      boosterStates.set(boosterId, {
        sequence: createSequenceState(),
        idempotency: createIdempotencyState(),
      });
    }
    return boosterStates.get(boosterId)!;
  };

  const processEvent = (eventData: EventData) => {
    const state = getOrCreateBoosterState(eventData.boosterId);

    // Check idempotency first
    if (idempotencyManager.hasProcessed(eventData.idempotencyId, state.idempotency)) {
      return;
    }

    // Check if should process based on sequence
    if (!eventProcessingService.shouldProcessEvent(eventData, state.sequence.lastSequence, state.idempotency.processedIds)) {
      return;
    }

    // Handle sequence logic
    if (sequenceManager.shouldProcessImmediately(eventData, state.sequence)) {
      processEventImmediately(eventData, state);
      processSequentialPendingEvents(eventData.boosterId, state);
    } else {
      sequenceManager.addPendingEvent(eventData, state.sequence);
    }
  };

  const processEventImmediately = (eventData: EventData, state: { sequence: SequenceState; idempotency: IdempotencyState }) => {
    // Process the event
    const processedEvent = eventProcessingService.processEvent(eventData);
    if (!processedEvent) return;

    // Mark as processed
    idempotencyManager.markAsProcessed(eventData.idempotencyId, state.idempotency);
    sequenceManager.updateSequence(eventData, state.sequence);

    // Dispatch to UI
    eventDispatcher.dispatch(eventData, processedEvent.executionUpdate);
  };

  const processSequentialPendingEvents = (boosterId: string, state: { sequence: SequenceState; idempotency: IdempotencyState }) => {
    const nextEvents = sequenceManager.getNextSequentialEvents(state.sequence);
    nextEvents.forEach(nextEvent => {
      processEventImmediately(nextEvent, state);
    });
  };

  const startCleanupRoutine = () => {
    cleanupInterval = setInterval(() => {
      boosterStates.forEach(state => {
        idempotencyManager.cleanup(state.idempotency, 300000); // 5 minutes
        sequenceManager.cleanupOldPendingEvents(state.sequence, 300000);
      });
    }, 30000); // Every 30 seconds
  };

  return {
    start() {
      if (isStarted) return;

      eventListener.onEvent(processEvent);
      eventListener.start();
      startCleanupRoutine();
      isStarted = true;
    },

    stop() {
      if (!isStarted) return;

      eventListener.stop();
      if (cleanupInterval) {
        clearInterval(cleanupInterval);
        cleanupInterval = undefined;
      }
      boosterStates.clear();
      isStarted = false;
    },

    setCallbacks(callbacks: EventCallbacks) {
      eventDispatcher.setCallbacks(callbacks);
    },

    isRunning() {
      return isStarted;
    },

    getStats() {
      let totalProcessed = 0;
      let totalPending = 0;

      boosterStates.forEach(state => {
        totalProcessed += state.idempotency.processedIds.size;
        totalPending += state.sequence.pendingEvents.size;
      });

      return {
        totalProcessed,
        totalPending,
        boostersTracked: boosterStates.size,
      };
    },
  };
}
