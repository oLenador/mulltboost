// src/presentation/features/boosters/domain/booster-events.handler.ts
import { Events } from '@wailsio/runtime';
import { BoosterEvent, BoosterEventData, ExecutionStatus } from './booster-queue.types';

export interface EventHandlerCallbacks {
  onExecutionStatusChanged: (boosterId: string, status: ExecutionStatus, progress?: number, error?: string) => void;
  onSyncRequired: () => void;
}

interface BoosterEventState {
  lastSequence: number;
  lastTimestamp: number;
  processedIdempotencyIds: Set<string>;
  pendingEvents: Map<number, BoosterEventData>; // sequência -> evento
  currentStatus: ExecutionStatus
}

export class BoosterEventsHandler {
  private unsubscribeFunctions: Map<string, () => void> = new Map();
  private isStarted = false;
  
  // Estado por booster para controle de sequência
  private boosterStates: Map<string, BoosterEventState> = new Map();
  
  // Configurações
  private readonly CONFIG = {
    IDEMPOTENCY_CLEANUP_INTERVAL: 30000, // 30 segundos
    IDEMPOTENCY_RETENTION_TIME: 300000,  // 5 minutos
    PENDING_EVENTS_TIMEOUT: 10000,       // 10 segundos
    MAX_PENDING_EVENTS: 50,              // Máximo de eventos pendentes por booster
  };

  private cleanupInterval?: NodeJS.Timeout;

  constructor(private callbacks: EventHandlerCallbacks) {}

  start(): void {
    if (this.isStarted) {
      console.warn('[BoosterEventsHandler] Already started, skipping...');
      return;
    }

    console.log('[BoosterEventsHandler] Starting event listeners...');
    this.setupEventListeners();
    this.startCleanupRoutine();
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
    this.boosterStates.clear();
    
    if (this.cleanupInterval) {
      clearInterval(this.cleanupInterval);
      this.cleanupInterval = undefined;
    }
    
    this.isStarted = false;
  }

  private setupEventListeners(): void {
    const events = [
      BoosterEvent.Processing,
      BoosterEvent.Success,
      BoosterEvent.Error,
      BoosterEvent.Failed,
      BoosterEvent.Queued,
      BoosterEvent.BatchQueued,
      BoosterEvent.Cancelled,
    ];

    events.forEach(eventType => {
      console.log(eventType)
      const unsubscribe = Events.On(eventType, (data: any) => {
        console.log("EventReceived: ", data)
        this.processIncomingEvent(eventType, data?.data);
      });
      this.unsubscribeFunctions.set(eventType, unsubscribe);
    });
  }

  private processIncomingEvent(eventType: string, rawData: any): void {
    try {
      // Validar e normalizar dados do evento
      const eventData = rawData as BoosterEventData;
 

      console.log(`[BoosterEventsHandler] Processing event: ${eventType}`, {
        boosterId: eventData.BoosterID,
        sequence: eventData.Sequency,
        idempotencyId: eventData.IdempotencyID
      });

      // Verificar se deve processar este evento
      if (!this.shouldProcessEvent(eventData)) {
        return;
      }

      // Processar evento imediatamente ou armazenar para processamento futuro
      this.handleEventWithSequence(eventData);
      
      // Tentar processar eventos pendentes
      this.processPendingEvents(eventData.BoosterID);

    } catch (error) {
      console.error(`[BoosterEventsHandler] Error processing event ${eventType}:`, error);
    }
  }

  private shouldProcessEvent(eventData: BoosterEventData): boolean {
    const state = this.getOrCreateBoosterState(eventData.BoosterID);

    // 1. Verificar idempotência
    if (state.processedIdempotencyIds.has(eventData.IdempotencyID)) {
      console.log(`[BoosterEventsHandler] Event already processed (idempotency): ${eventData.IdempotencyID}`);
      return false;
    }

    // 2. Verificar se a sequência é válida (não é menor que a última processada)
    if (eventData.Sequency > 0 && eventData.Sequency <= state.lastSequence) {
      console.log(`[BoosterEventsHandler] Event sequence too old: ${eventData.Sequency} <= ${state.lastSequence}`);
      return false;
    }

    return true;
  }

  private handleEventWithSequence(eventData: BoosterEventData): void {
    const state = this.getOrCreateBoosterState(eventData.BoosterID);

    // Se não tem sequência ou é a próxima esperada, processar imediatamente
    if (eventData.Sequency === 0 || eventData.Sequency === state.lastSequence + 1) {
      this.processEvent(eventData);
      return;
    }

    // Se a sequência é maior que a esperada, armazenar para processar depois
    if (eventData.Sequency > state.lastSequence + 1) {
      console.log(`[BoosterEventsHandler] Event sequence out of order, storing: ${eventData.Sequency} (expected: ${state.lastSequence + 1})`);
      
      // Limitar número de eventos pendentes para evitar memory leak
      if (state.pendingEvents.size >= this.CONFIG.MAX_PENDING_EVENTS) {
        console.warn(`[BoosterEventsHandler] Too many pending events for booster ${eventData.BoosterID}, processing oldest`);
        const oldestSequence = Math.min(...state.pendingEvents.keys());
        const oldestEvent = state.pendingEvents.get(oldestSequence);
        if (oldestEvent) {
          state.pendingEvents.delete(oldestSequence);
          this.processEvent(oldestEvent);
        }
      }
      
      state.pendingEvents.set(eventData.Sequency, eventData);
      
      // Configurar timeout para processar evento mesmo sem a sequência correta
      setTimeout(() => {
        if (state.pendingEvents.has(eventData.Sequency)) {
          console.warn(`[BoosterEventsHandler] Processing event with gap in sequence: ${eventData.Sequency}`);
          const event = state.pendingEvents.get(eventData.Sequency);
          if (event) {
            state.pendingEvents.delete(eventData.Sequency);
            this.processEvent(event);
          }
        }
      }, this.CONFIG.PENDING_EVENTS_TIMEOUT);
    }
  }

  private processPendingEvents(boosterId: string): void {
    const state = this.getOrCreateBoosterState(boosterId);
    
    // Processar eventos pendentes em ordem de sequência
    const sortedSequences = Array.from(state.pendingEvents.keys()).sort((a, b) => a - b);
    
    for (const sequence of sortedSequences) {
      // Só processar se for a próxima sequência esperada
      if (sequence === state.lastSequence + 1) {
        const event = state.pendingEvents.get(sequence);
        if (event) {
          state.pendingEvents.delete(sequence);
          this.processEvent(event);
        }
      } else {
        break; // Parar no primeiro gap
      }
    }
  }

  private processEvent(eventData: BoosterEventData): void {
    const state = this.getOrCreateBoosterState(eventData.BoosterID);

    // Marcar como processado
    state.processedIdempotencyIds.add(eventData.IdempotencyID);
    
    // Atualizar sequência se aplicável
    if (eventData.Sequency > 0) {
      state.lastSequence = eventData.Sequency;
    }
    
    state.lastTimestamp = new Date(eventData.Timestamp).getTime();

    console.log(`[BoosterEventsHandler] Processing event: ${eventData.EventType} for ${eventData.BoosterID}`);

    // Processar baseado no tipo de evento
    if (this.isSingleBoosterEvent(eventData)) {
      this.handleSingleBoosterEvent(eventData);
    } else {
      this.callbacks.onSyncRequired();
    }
  }

  private handleSingleBoosterEvent(data: BoosterEventData): void {
    const { BoosterID, Error } = data;
    const state = this.getOrCreateBoosterState(BoosterID);

    let newStatus: ExecutionStatus;

    switch (data.EventType) {
      case BoosterEvent.Queued:
        newStatus = 'queued';
        break;
      case BoosterEvent.Processing:
        newStatus = 'processing';
        break;
      case BoosterEvent.Success:
        newStatus = 'completed';
        this.callbacks.onExecutionStatusChanged(BoosterID, newStatus, 100);
        return;
      case BoosterEvent.Error:
      case BoosterEvent.Failed:
        newStatus = 'error';
        this.callbacks.onExecutionStatusChanged(BoosterID, newStatus, undefined, Error);
        return;
      case BoosterEvent.Cancelled:
        newStatus = 'cancelled';
        break;
      default:
        console.warn(`[BoosterEventsHandler] Unknown event type: ${data.EventType}`);
        return;
    }

    // Atualizar estado e notificar apenas se mudou
    if (state.currentStatus !== newStatus) {
      state.currentStatus = newStatus;
      this.callbacks.onExecutionStatusChanged(BoosterID, newStatus);
    }
  }

  private getOrCreateBoosterState(boosterId: string): BoosterEventState {
    if (!this.boosterStates.has(boosterId)) {
      this.boosterStates.set(boosterId, {
        lastSequence: 0,
        lastTimestamp: 0,
        processedIdempotencyIds: new Set(),
        pendingEvents: new Map(),
        currentStatus: 'queued',
      });
    }
    return this.boosterStates.get(boosterId)!;
  }

  private isSingleBoosterEvent(data: any): data is BoosterEventData {
    return typeof data === 'object' && data !== null && 'BoosterID' in data;
  }

  private startCleanupRoutine(): void {
    this.cleanupInterval = setInterval(() => {
      this.cleanupOldData();
    }, this.CONFIG.IDEMPOTENCY_CLEANUP_INTERVAL);
  }

  private cleanupOldData(): void {
    const now = Date.now();
    const cutoffTime = now - this.CONFIG.IDEMPOTENCY_RETENTION_TIME;

    console.log('[BoosterEventsHandler] Running cleanup routine...');

    this.boosterStates.forEach((state, boosterId) => {
      // Limpar IDs de idempotência antigos
      const sizeBefore = state.processedIdempotencyIds.size;
      
      // Como Set não tem método de filtro por data, recreamos o set
      // (Esta é uma limitação - idealmente armazenaríamos timestamp com cada ID)
      if (state.lastTimestamp < cutoffTime) {
        state.processedIdempotencyIds.clear();
      }

      // Limpar eventos pendentes muito antigos
      const pendingToRemove: number[] = [];
      state.pendingEvents.forEach((event, sequence) => {
        const eventTime = new Date(event.Timestamp).getTime();
        if (eventTime < cutoffTime) {
          pendingToRemove.push(sequence);
        }
      });
      
      pendingToRemove.forEach(sequence => {
        console.warn(`[BoosterEventsHandler] Removing stale pending event: ${boosterId}-${sequence}`);
        state.pendingEvents.delete(sequence);
      });

      const sizeAfter = state.processedIdempotencyIds.size;
      if (sizeBefore !== sizeAfter || pendingToRemove.length > 0) {
        console.log(`[BoosterEventsHandler] Cleaned up booster ${boosterId}: ${sizeBefore - sizeAfter} idempotency IDs, ${pendingToRemove.length} pending events`);
      }
    });
  }

  // Métodos públicos para debug e monitoramento
  isRunning(): boolean {
    return this.isStarted;
  }

  getBoosterState(boosterId: string): Readonly<BoosterEventState> | undefined {
    return this.boosterStates.get(boosterId);
  }

  getStats(): {
    totalBoosters: number;
    totalProcessedEvents: number;
    totalPendingEvents: number;
  } {
    let totalProcessedEvents = 0;
    let totalPendingEvents = 0;

    this.boosterStates.forEach(state => {
      totalProcessedEvents += state.processedIdempotencyIds.size;
      totalPendingEvents += state.pendingEvents.size;
    });

    return {
      totalBoosters: this.boosterStates.size,
      totalProcessedEvents,
      totalPendingEvents,
    };
  }

  // Método para forçar limpeza manual
  forceCleanup(): void {
    this.cleanupOldData();
  }
}