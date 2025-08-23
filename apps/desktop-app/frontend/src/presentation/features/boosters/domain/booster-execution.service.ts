// src/presentation/features/boosters/domain/booster-execution.service.ts
import { GetExecutionQueueState, InitBoosterApply, InitRevertBooster } from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { QueueState, ExecutionStatus } from './booster-queue.types';
import { BoosterEventsHandler, EventHandlerCallbacks } from './booster-events.handler';

export interface ExecutionServiceCallbacks extends EventHandlerCallbacks {
  onBatchStarted: (batchId: string) => void;
  onBatchCompleted: (batchId: string) => void;
  onBatchError: (batchId: string, error: string) => void;
}

export class BoosterExecutionService {
  private eventsHandler: BoosterEventsHandler;
  private isInitialized = false;
  private isSyncing = false; 
  private static instance: BoosterExecutionService | null = null;

  constructor(private callbacks: ExecutionServiceCallbacks) {
    this.eventsHandler = new BoosterEventsHandler(callbacks);
  }

  static getInstance(callbacks: ExecutionServiceCallbacks): BoosterExecutionService {
    if (!BoosterExecutionService.instance) {
      BoosterExecutionService.instance = new BoosterExecutionService(callbacks);
    } else {
      // Atualiza callbacks se necessário
      BoosterExecutionService.instance.updateCallbacks(callbacks);
    }
    return BoosterExecutionService.instance;
  }

  private updateCallbacks(callbacks: ExecutionServiceCallbacks): void {
    this.callbacks = callbacks;

    if (this.eventsHandler) {
      this.eventsHandler = new BoosterEventsHandler(callbacks);
    }
  }

  initialize(): void {
    if (this.isInitialized) {
      console.warn('[BoosterExecutionService] Already initialized, skipping...');
      return;
    }

    console.log('[BoosterExecutionService] Initializing...');
    
    try {
      this.eventsHandler.start();
      this.isInitialized = true;
      console.log('[BoosterExecutionService] Initialized successfully');
    } catch (error) {
      console.error('[BoosterExecutionService] Failed to initialize:', error);
      this.isInitialized = false;
      throw error;
    }
  }

  destroy(): void {
    if (!this.isInitialized) {
      console.warn('[BoosterExecutionService] Not initialized, skipping destroy...');
      return;
    }

    console.log('[BoosterExecutionService] Destroying...');
    
    try {
      this.eventsHandler.stop();
      this.isInitialized = false;
      BoosterExecutionService.instance = null;
      console.log('[BoosterExecutionService] Destroyed successfully');
    } catch (error) {
      console.error('[BoosterExecutionService] Failed to destroy:', error);
    }
  }

  async executeBatch(operations: Record<string, BoosterOperationType>): Promise<string> {
    if (!this.isInitialized) {
      throw new Error('[BoosterExecutionService] Service not initialized');
    }

    if (Object.keys(operations).length === 0) {
      throw new Error('No operations to execute');
    }

    const batchId = this.generateBatchId();
    
    try {
      console.log(`[BoosterExecutionService] Starting batch ${batchId} with ${Object.keys(operations).length} operations`);
      this.callbacks.onBatchStarted(batchId);

      for (const [boosterId, operation] of Object.entries(operations)) {
        await this.executeBooster(boosterId, operation);
        // Little delay to prevent overload
        await new Promise(resolve => setTimeout(resolve, 100));
      }

      this.callbacks.onBatchCompleted(batchId);
      console.log(`[BoosterExecutionService] Batch ${batchId} completed successfully`);
      return batchId;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Batch execution failed';
      console.error(`[BoosterExecutionService] Batch ${batchId} failed:`, errorMessage);
      this.callbacks.onBatchError(batchId, errorMessage);
      throw error;
    }
  }

  private async executeBooster(boosterId: string, operation: BoosterOperationType): Promise<void> {
    try {
      console.log(`[BoosterExecutionService] Executing ${operation} on ${boosterId}`);
      
      const result = operation === 'apply' 
        ? await InitBoosterApply(boosterId) 
        : await InitRevertBooster(boosterId);

      if (!result.Success) {
        const error = result.Error || result.Message || `Failed to ${operation} booster`;
        throw new Error(error);
      }

      console.log(`[BoosterExecutionService] Successfully initiated ${operation} for ${boosterId}, Operation ID: ${result.OperationID}`);
    } catch (error) {
      console.error(`[BoosterExecutionService] Failed to ${operation} booster ${boosterId}:`, error);
      throw error;
    }
  }

  async syncWithBackend(): Promise<void> {
    if (this.isSyncing) {
      console.warn('[BoosterExecutionService] Sync already in progress, skipping...');
      return;
    }

    if (!this.isInitialized) {
      console.warn('[BoosterExecutionService] Service not initialized, skipping sync...');
      return;
    }

    this.isSyncing = true;
    
    try {
      console.log('[BoosterExecutionService] Syncing with backend...');
      const stateOrArray: any = await GetExecutionQueueState();
      
      if (!stateOrArray) {
        console.log('[BoosterExecutionService] No backend state to sync');
        return;
      }

      const queueState = this.normalizeQueueState(stateOrArray);
      this.processQueueState(queueState);
      console.log('[BoosterExecutionService] Backend sync completed');
    } catch (error) {
      console.error('[BoosterExecutionService] Failed to sync with backend:', error);
      throw error;
    } finally {
      this.isSyncing = false;
    }
  }

  private normalizeQueueState(data: any): QueueState {
    if (Array.isArray(data)) {
      return { items: data, InProgress: 0 };
    }
    return {
      items: data.Items || [],
      InProgress: data.InProgress || data.inProgress || 0
    };
  }

  private processQueueState(queueState: QueueState): void {
    const { items, InProgress } = queueState;

    const processingSet = new Set<string>();
    for (let i = 0; i < Math.min(InProgress, items.length); i++) {
      const id = items[i]?.BoosterID;
      if (id) processingSet.add(id);
    }

    // Update status for items in queue
    items.forEach(item => {
      if (!item?.BoosterID) return;

      const status: ExecutionStatus = processingSet.has(item.BoosterID) ? 'processing' : 'queued';
      const progress = typeof item.Progress === 'number' ? item.Progress : 0;

      this.callbacks.onExecutionStatusChanged(
        item.BoosterID,
        status,
        progress,
        item.Error
      );
    });

    // Trigger sync for UI updates
    this.callbacks.onSyncRequired();
  }

  private generateBatchId(): string {
    return `batch_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  isRunning(): boolean {
    return this.isInitialized && this.eventsHandler.isRunning();
  }

  static cleanup(): void {
    if (BoosterExecutionService.instance) {
      BoosterExecutionService.instance.destroy();
    }
  }
}