// src/presentation/features/boosters/domain/booster-events.handler.ts
import { Events } from '@wailsio/runtime';
import { BoosterEvent, BoosterEventData, ExecutionStatus } from './booster-queue.types';

export interface EventHandlerCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: ExecutionStatus, progress?: number, error?: string) => void;
  onSyncRequired: () => void;
}

export class BoosterEventsHandler {
  private unsubscribeFunctions: Map<string, () => void> = new Map();
  private isStarted = false; 
  private processedEvents = new Set<string>(); 

  constructor(private callbacks: EventHandlerCallbacks) {}

  start(): void {
    if (this.isStarted) {
      console.warn('[BoosterEventsHandler] Already started, skipping...');
      return;
    }

    console.log('[BoosterEventsHandler] Starting event listeners...');
    
    const events = [
      BoosterEvent.BatchQueued,
      BoosterEvent.Cancelled,
      BoosterEvent.Error,
      BoosterEvent.Failed,
      BoosterEvent.Processing,
      BoosterEvent.Queued,
      BoosterEvent.Success,
    ];

    events.forEach(eventType => {
      const unsubscribe = Events.On(eventType, (data: any) => {
        const eventId = this.generateEventId(eventType, data?.data);
        
        if (this.processedEvents.has(eventId)) {
          console.warn(`[BoosterEventsHandler] Duplicate event ignored: ${eventId}`);
          return;
        }
        
        this.processedEvents.add(eventId);
        console.log(`[BoosterEventsHandler] EventReceived: ${eventType}`, data?.data);
        
        this.handleEvent(eventType, data?.data);
        
        setTimeout(() => {
          this.processedEvents.delete(eventId);
        }, 5000);
      });

      this.unsubscribeFunctions.set(eventType, unsubscribe);
    });

    this.isStarted = true;
  }

  stop(): void {
    if (!this.isStarted) {
      console.warn('[BoosterEventsHandler] Not started, skipping stop...');
      return;
    }

    console.log('[BoosterEventsHandler] Stopping event listeners...');
    
    this.unsubscribeFunctions.forEach(unsubscribe => unsubscribe());
    this.unsubscribeFunctions.clear();
    this.processedEvents.clear();
    this.isStarted = false;
  }

  private generateEventId(eventType: string, data: any): string {
    const operationId = data?.OperationID || '';
    const boosterId = data?.BoosterID || '';
    const timestamp = data?.Timestamp || Date.now();
    
    return `${eventType}-${boosterId}-${operationId}-${timestamp}`;
  }

  private handleEvent(eventType: string, data: any): void {
    try {
      if (this.isSingleBoosterEvent(data)) {
        this.handleSingleBoosterEvent(eventType, data);
      } else {
        this.callbacks.onSyncRequired();
      }
    } catch (error) {
      console.error(`[BoosterEventsHandler] Error handling event ${eventType}:`, error);
    }
  }

  private handleSingleBoosterEvent(eventType: string, data: BoosterEventData): void {
    const { BoosterID, Error } = data;

    switch (eventType) {
      case BoosterEvent.Queued:
        this.callbacks.onExecutionStatusChanged(BoosterID, 'queued');
        break;
      case BoosterEvent.Processing:
        this.callbacks.onExecutionStatusChanged(BoosterID, 'processing');
        break;
      case BoosterEvent.Success:
        this.callbacks.onExecutionStatusChanged(BoosterID, 'completed', 100);
        break;
      case BoosterEvent.Error:
      case BoosterEvent.Failed:
        this.callbacks.onExecutionStatusChanged(BoosterID, 'error', undefined, Error);
        break;
      case BoosterEvent.Cancelled:
        this.callbacks.onExecutionStatusChanged(BoosterID, 'cancelled');
        break;
      default:
        console.warn(`[BoosterEventsHandler] Unknown event type: ${eventType}`);
    }
  }

  private isSingleBoosterEvent(data: any): data is BoosterEventData {
    return typeof data === 'object' && data !== null && 'BoosterID' in data;
  }

  isRunning(): boolean {
    return this.isStarted;
  }
}