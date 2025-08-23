// src/core/services/global-booster-service.ts

import { BoosterExecutionService, ExecutionServiceCallbacks } from "./booster-execution.service";

type ServiceSubscriber = {
  id: string;
  callbacks: ExecutionServiceCallbacks;
};


class GlobalBoosterService {
  private static instance: GlobalBoosterService | null = null;
  private service: BoosterExecutionService | null = null;
  private subscribers: Map<string, ServiceSubscriber> = new Map();
  private isServiceInitialized = false;

  private constructor() {}

  static getInstance(): GlobalBoosterService {
    if (!GlobalBoosterService.instance) {
      GlobalBoosterService.instance = new GlobalBoosterService();
    }
    return GlobalBoosterService.instance;
  }

  /**
   * Registra um subscriber (página/componente) para receber callbacks
   */
  subscribe(subscriberId: string, callbacks: ExecutionServiceCallbacks): void {
    console.log(`[GlobalBoosterService] Subscribing: ${subscriberId}`);
    
    this.subscribers.set(subscriberId, {
      id: subscriberId,
      callbacks
    });

    // Inicializa o service se for o primeiro subscriber
    if (!this.isServiceInitialized) {
      this.initializeService();
    }
  }

  /**
   * Remove um subscriber
   */
  unsubscribe(subscriberId: string): void {
    console.log(`[GlobalBoosterService] Unsubscribing: ${subscriberId}`);
    
    this.subscribers.delete(subscriberId);

    // Se não há mais subscribers, mantém o service ativo
    // (BoosterStatusProvider sempre existe)
    console.log(`[GlobalBoosterService] Active subscribers: ${this.subscribers.size}`);
  }

  /**
   * Executa um batch de operações
   */
  async executeBatch(operations: Record<string, any>): Promise<string> {
    if (!this.service || !this.isServiceInitialized) {
      throw new Error('Global booster service not initialized');
    }

    return this.service.executeBatch(operations);
  }

  /**
   * Sincroniza com o backend
   */
  async syncWithBackend(): Promise<void> {
    if (!this.service || !this.isServiceInitialized) {
      console.warn('[GlobalBoosterService] Service not available for sync');
      return;
    }

    return this.service.syncWithBackend();
  }

  /**
   * Verifica se o service está rodando
   */
  isRunning(): boolean {
    return this.isServiceInitialized && this.service?.isRunning() === true;
  }

  /**
   * Inicializa o service com callbacks agregados
   */
  private initializeService(): void {
    if (this.isServiceInitialized) {
      console.warn('[GlobalBoosterService] Service already initialized');
      return;
    }

    console.log('[GlobalBoosterService] Initializing service...');

    // Cria callbacks que distribuem para todos os subscribers
    const aggregatedCallbacks: ExecutionServiceCallbacks = {
      onExecutionStatusChanged: (boosterId, status, progress, error) => {
        this.subscribers.forEach(subscriber => {
          try {
            subscriber.callbacks.onExecutionStatusChanged(boosterId, status, progress, error);
          } catch (error) {
            console.error(`[GlobalBoosterService] Error in subscriber ${subscriber.id}:`, error);
          }
        });
      },

      onSyncRequired: async () => {
        // Executa sequencialmente para evitar conflitos
        for (const subscriber of this.subscribers.values()) {
          try {
            await subscriber.callbacks.onSyncRequired();
          } catch (error) {
            console.error(`[GlobalBoosterService] Sync error in subscriber ${subscriber.id}:`, error);
          }
        }
      },

      onBatchStarted: (batchId) => {
        this.subscribers.forEach(subscriber => {
          try {
            subscriber.callbacks.onBatchStarted(batchId);
          } catch (error) {
            console.error(`[GlobalBoosterService] Batch start error in subscriber ${subscriber.id}:`, error);
          }
        });
      },

      onBatchCompleted: (batchId) => {
        this.subscribers.forEach(subscriber => {
          try {
            subscriber.callbacks.onBatchCompleted(batchId);
          } catch (error) {
            console.error(`[GlobalBoosterService] Batch complete error in subscriber ${subscriber.id}:`, error);
          }
        });
      },

      onBatchError: (batchId, error) => {
        this.subscribers.forEach(subscriber => {
          try {
            subscriber.callbacks.onBatchError(batchId, error);
          } catch (error) {
            console.error(`[GlobalBoosterService] Batch error in subscriber ${subscriber.id}:`, error);
          }
        });
      },
    };

    this.service = new BoosterExecutionService(aggregatedCallbacks);
    this.service.initialize();
    this.isServiceInitialized = true;

    console.log('[GlobalBoosterService] Service initialized successfully');
  }

  /**
   * Cleanup completo - apenas para testes
   */
  destroy(): void {
    console.log('[GlobalBoosterService] Destroying service...');
    
    if (this.service && this.isServiceInitialized) {
      this.service.destroy();
    }
    
    this.service = null;
    this.subscribers.clear();
    this.isServiceInitialized = false;
    GlobalBoosterService.instance = null;
  }
}

export const globalBoosterService = GlobalBoosterService.getInstance();