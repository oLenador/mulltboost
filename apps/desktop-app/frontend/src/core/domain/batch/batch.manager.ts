import { ProcessableItem, ProcessableItemWithStatus, ItemStatus, BatchOperation, BatchManagerEvents } from './batch.types';

export abstract class BaseBatchManager<T extends ProcessableItem> {
  protected items = new Map<string, ProcessableItemWithStatus<T>>();
  protected operations = new Map<string, BatchOperation>();
  protected observers: ((items: ProcessableItemWithStatus<T>[]) => void)[] = [];
  protected events: Partial<BatchManagerEvents<T>> = {};

  constructor(protected maxConcurrent: number = 5) {}

  // Observer Pattern
  subscribe(observer: (items: ProcessableItemWithStatus<T>[]) => void): () => void {
    this.observers.push(observer);
    return () => {
      const index = this.observers.indexOf(observer);
      if (index > -1) {
        this.observers.splice(index, 1);
      }
    };
  }

  protected notify(): void {
    const items = Array.from(this.items.values());
    this.observers.forEach(observer => observer(items));
  }

  // Events
  on<K extends keyof BatchManagerEvents<T>>(event: K, handler: BatchManagerEvents<T>[K]): void {
    this.events[event] = handler;
  }

  protected emit<K extends keyof BatchManagerEvents<T>>(
    event: K, 
    ...args: Parameters<NonNullable<BatchManagerEvents<T>[K]>>
  ): void {
    const handler = this.events[event];
    if (handler) {
      (handler as any)(...args);
    }
  }

  // Core Methods
  addItems(items: T[]): void {
    items.forEach(item => {
      this.items.set(item.id, {
        item,
        status: 'idle',
        progress: 0,
        canCancel: false,
      });
    });
    this.notify();
  }

  removeItems(ids: string[]): void {
    ids.forEach(id => {
      const item = this.items.get(id);
      if (item && item.status === 'idle') {
        this.items.delete(id);
      }
    });
    this.notify();
  }

  getItems(): ProcessableItemWithStatus<T>[] {
    return Array.from(this.items.values());
  }

  getItemsByStatus(status: ItemStatus): ProcessableItemWithStatus<T>[] {
    return this.getItems().filter(item => item.status === status);
  }

  getItem(id: string): ProcessableItemWithStatus<T> | undefined {
    return this.items.get(id);
  }

  async startBatch(ids?: string[]): Promise<string> {
    const itemsToProcess = ids 
      ? ids.map(id => this.items.get(id)).filter(Boolean) as ProcessableItemWithStatus<T>[]
      : this.getItemsByStatus('idle');

    if (itemsToProcess.length === 0) {
      throw new Error('No items to process');
    }

    const operationId = this.generateId();
    const operation: BatchOperation = {
      id: operationId,
      name: `Batch Operation ${new Date().toLocaleTimeString()}`,
      itemIds: itemsToProcess.map(item => item.item.id),
      status: 'queued',
      progress: 0,
      createdAt: new Date()
    };

    this.operations.set(operationId, operation);
    
    // Queue items
    itemsToProcess.forEach(item => {
      this.updateItemStatus(item.item.id, 'queued');
    });

    this.emit('onBatchStarted', operation);
    
    // Start processing
    this.processBatch(operationId);
    
    return operationId;
  }

  cancelItems(ids: string[]): void {
    ids.forEach(id => {
      const item = this.items.get(id);
      if (item && item.canCancel) {
        this.updateItemStatus(id, 'cancelled');
      }
    });
  }

  protected updateItemStatus(
    id: string, 
    status: ItemStatus, 
    progress?: number, 
    error?: string
  ): void {
    const item = this.items.get(id);
    if (!item) return;

    const now = new Date();
    const updatedItem: ProcessableItemWithStatus<T> = {
      ...item,
      status,
      progress: progress ?? item.progress,
      error,
      canCancel: status === 'processing' || status === 'queued',
      ...(status === 'processing' && !item.startedAt ? { startedAt: now } : {}),
      ...(status === 'completed' || status === 'error' || status === 'cancelled' 
        ? { completedAt: now, canCancel: false } 
        : {})
    };

    this.items.set(id, updatedItem);
    this.emit('onItemStatusChanged', updatedItem);
    this.notify();
  }

  private async processBatch(operationId: string): Promise<void> {
    const operation = this.operations.get(operationId);
    if (!operation) return;

    try {
      operation.status = 'processing';
      operation.startedAt = new Date();
      
      const items = operation.itemIds.map(id => this.items.get(id)!).filter(Boolean);
      const chunks = this.chunkArray(items, this.maxConcurrent);
      
      let processedCount = 0;
      
      for (const chunk of chunks) {
        await Promise.all(
          chunk.map(async (item) => {
            if (item.status === 'cancelled') return;
            
            try {
              this.updateItemStatus(item.item.id, 'processing');
              await this.processItem(item.item);
              this.updateItemStatus(item.item.id, 'completed', 100);
            } catch (error) {
              const errorMessage = error instanceof Error ? error.message : 'Unknown error';
              this.updateItemStatus(item.item.id, 'error', item.progress, errorMessage);
            }
            
            processedCount++;
            operation.progress = (processedCount / items.length) * 100;
          })
        );
      }

      operation.status = 'completed';
      operation.completedAt = new Date();
      operation.progress = 100;
      
      this.emit('onBatchCompleted', operation);
      
    } catch (error) {
      operation.status = 'error';
      operation.completedAt = new Date();
      
      const errorMessage = error instanceof Error ? error.message : 'Batch processing failed';
      this.emit('onBatchError', operation, errorMessage);
    }
  }

  private chunkArray<U>(array: U[], size: number): U[][] {
    const chunks: U[][] = [];
    for (let i = 0; i < array.length; i += size) {
      chunks.push(array.slice(i, i + size));
    }
    return chunks;
  }

  private generateId(): string {
    return `batch_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}