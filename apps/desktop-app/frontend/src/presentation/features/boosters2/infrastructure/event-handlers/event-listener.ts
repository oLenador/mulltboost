import { Events } from '@wailsio/runtime';
import { EventType, EventData, createEventData } from '../../core/types/events.types';

export interface EventListener {
  start(): void;
  stop(): void;
  onEvent(callback: (eventData: EventData) => void): void;
}

export function createEventListener(): EventListener {
  const unsubscribeFunctions = new Map<string, () => void>();
  let eventCallback: ((eventData: EventData) => void) | null = null;
  let isStarted = false;

  return {
    start(): void {
      if (isStarted) return;

      const events = [
        EventType.Processing,
        EventType.Success,
        EventType.Error,
        EventType.Failed,
        EventType.Queued,
        EventType.BatchQueued,
        EventType.Cancelled,
      ];

      events.forEach(eventType => {
        const unsubscribe = Events.On(eventType, (data: any) => {
          if (eventCallback) {
            try {
              const eventData = createEventData(data?.data);
              eventCallback(eventData);
            } catch (error) {
              console.error(`Error processing event ${eventType}:`, error);
            }
          }
        });
        unsubscribeFunctions.set(eventType, unsubscribe);
      });

      isStarted = true;
    },

    stop(): void {
      if (!isStarted) return;

      unsubscribeFunctions.forEach(unsubscribe => unsubscribe());
      unsubscribeFunctions.clear();
      eventCallback = null;
      isStarted = false;
    },

    onEvent(callback: (eventData: EventData) => void): void {
      eventCallback = callback;
    },
  };
}