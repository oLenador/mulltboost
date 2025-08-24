import { EventData } from '../../core/types/events.types';

export interface SequenceState {
  lastSequence: number;
  lastTimestamp: number;
  pendingEvents: Map<number, EventData>;
}

export interface SequenceManager {
  shouldProcessImmediately(eventData: EventData, state: SequenceState): boolean;
  addPendingEvent(eventData: EventData, state: SequenceState): void;
  getNextSequentialEvents(state: SequenceState): EventData[];
  updateSequence(eventData: EventData, state: SequenceState): void;
  cleanupOldPendingEvents(state: SequenceState, maxAge: number): void;
}

export function createSequenceManager(config: {
  maxPendingEvents: number;
  pendingTimeout: number;
}): SequenceManager {
  const { maxPendingEvents, pendingTimeout } = config;

  return {
    shouldProcessImmediately(eventData: EventData, state: SequenceState): boolean {
      return eventData.sequence === 0 || eventData.sequence === state.lastSequence + 1;
    },

    addPendingEvent(eventData: EventData, state: SequenceState): void {
      if (eventData.sequence <= state.lastSequence + 1) return;

      // Remove oldest if at capacity
      if (state.pendingEvents.size >= maxPendingEvents) {
        const oldestSequence = Math.min(...state.pendingEvents.keys());
        state.pendingEvents.delete(oldestSequence);
      }

      state.pendingEvents.set(eventData.sequence, eventData);

      // Auto-cleanup after timeout
      setTimeout(() => {
        if (state.pendingEvents.has(eventData.sequence)) {
          const event = state.pendingEvents.get(eventData.sequence);
          if (event) {
            state.pendingEvents.delete(eventData.sequence);
          }
        }
      }, pendingTimeout);
    },

    getNextSequentialEvents(state: SequenceState): EventData[] {
      const sequentialEvents: EventData[] = [];
      const sortedSequences = Array.from(state.pendingEvents.keys()).sort((a, b) => a - b);

      for (const sequence of sortedSequences) {
        if (sequence === state.lastSequence + 1) {
          const event = state.pendingEvents.get(sequence);
          if (event) {
            sequentialEvents.push(event);
            state.pendingEvents.delete(sequence);
            state.lastSequence = sequence;
          }
        } else {
          break;
        }
      }

      return sequentialEvents;
    },

    updateSequence(eventData: EventData, state: SequenceState): void {
      if (eventData.sequence > 0) {
        state.lastSequence = eventData.sequence;
      }
      state.lastTimestamp = eventData.timestamp.getTime();
    },

    cleanupOldPendingEvents(state: SequenceState, maxAge: number): void {
      const cutoffTime = Date.now() - maxAge;
      const toRemove: number[] = [];

      state.pendingEvents.forEach((event, sequence) => {
        if (event.timestamp.getTime() < cutoffTime) {
          toRemove.push(sequence);
        }
      });

      toRemove.forEach(sequence => state.pendingEvents.delete(sequence));
    },
  };
}

export function createSequenceState(): SequenceState {
  return {
    lastSequence: 0,
    lastTimestamp: 0,
    pendingEvents: new Map(),
  };
}
