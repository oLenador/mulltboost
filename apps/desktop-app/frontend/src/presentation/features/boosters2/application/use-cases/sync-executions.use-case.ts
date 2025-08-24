import { SyncService } from '../../core/services/sync.service';
import { BoosterApiService } from '../../core/services/booster-api.service';
import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';

export interface SyncExecutionsResponse {
  success: boolean;
  updatedCount: number;
  error?: string;
}

export interface SyncCallbacks {
  onExecutionUpdated: (boosterId: string, status: string, progress?: number, error?: string) => void;
  onSyncCompleted: (updatedCount: number) => void;
}

export interface SyncExecutionsUseCase {
  execute(): Promise<SyncExecutionsResponse>;
  setCallbacks(callbacks: SyncCallbacks): void;
}

export function createSyncExecutionsUseCase(
  syncService: SyncService,
  apiService: BoosterApiService,
  repository: ExecutionStateRepository
): SyncExecutionsUseCase {
  let callbacks: SyncCallbacks | null = null;

  return {
    async execute(): Promise<SyncExecutionsResponse> {
      try {
        // Get queue state from backend
        const queueState = await apiService.getQueueState();

        // Process queue state into execution updates
        const executionUpdates = syncService.processQueueState(queueState);

        // Update repository
        let updatedCount = 0;
        executionUpdates.forEach(update => {
          repository.updateExecution(update.boosterId, update);
          
          // Notify callbacks
          if (callbacks) {
            callbacks.onExecutionUpdated(
              update.boosterId,
              update.status,
              update.progress,
              update.error
            );
          }
          
          updatedCount++;
        });

        // Notify sync completed
        callbacks?.onSyncCompleted(updatedCount);

        return { success: true, updatedCount };
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Sync failed';
        return { success: false, updatedCount: 0, error: errorMessage };
      }
    },

    setCallbacks(newCallbacks: SyncCallbacks): void {
      callbacks = newCallbacks;
    },
  };
}