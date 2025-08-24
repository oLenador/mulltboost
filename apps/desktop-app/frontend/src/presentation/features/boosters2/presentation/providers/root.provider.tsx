
// src/presentation/features/boosters/presentation/providers/root.provider.tsx
import React, { ReactNode, useMemo } from 'react';
import { ExecutionProvider } from './execution.provider';
import { BoosterContextProvider } from './booster-context.provider';

// Infrastructure
import { createInMemoryExecutionRepository } from '../../infrastructure/storage/execution-storage';
import { createBackendClient } from '../../infrastructure/external/backend-client';
import { createEventListener } from '../../infrastructure/event-handlers/event-listener';
import { createSequenceManager } from '../../infrastructure/event-handlers/sequence-manager';
import { createIdempotencyManager } from '../../infrastructure/event-handlers/idempotency-manager';
import { createEventDispatcher } from '../../infrastructure/event-handlers/event-dispatcher';

// Core Services
import { createEventProcessingService } from '../../core/services/event-processing.service';
import { createExecutionService } from '../../core/services/execution.service';
import { createSyncService } from '../../core/services/sync.service';

// Use Cases
import { createLoadBoostersUseCase } from '../../application/use-cases/load-boosters.use-case';
import { createManageStagingUseCase } from '../../application/use-cases/manage-staging.use-case';
import { createExecuteBatchUseCase } from '../../application/use-cases/execute-batch.use-case';
import { createSyncExecutionsUseCase } from '../../application/use-cases/sync-executions.use-case';

// Orchestrators
import { createExecutionOrchestrator } from '../../application/orchestrators/execution.orchestrator';
import { createEventOrchestrator } from '../../application/orchestrators/event.orchestrator';
import { createInMemoryBoosterRepository } from '../../infrastructure/storage/booster-storage';

export interface BoosterRootProviderProps {
  children: ReactNode;
}

export function BoosterRootProvider({ children }: BoosterRootProviderProps) {
  const dependencies = useMemo(() => {
    // Infrastructure layer
    const executionRepository = createInMemoryExecutionRepository();
    const boosterRepository = createInMemoryBoosterRepository();
    const backendClient = createBackendClient();
    
    const eventListener = createEventListener();
    const sequenceManager = createSequenceManager({
      maxPendingEvents: 50,
      pendingTimeout: 10000,
    });
    const idempotencyManager = createIdempotencyManager();
    const eventDispatcher = createEventDispatcher();

    // Core services
    const eventProcessingService = createEventProcessingService();
    const executionService = createExecutionService();
    const syncService = createSyncService();

    // Use cases
    const loadBoostersUseCase = createLoadBoostersUseCase(backendClient, boosterRepository);
    const manageStagingUseCase = createManageStagingUseCase(executionRepository, boosterRepository);
    const executeBatchUseCase = createExecuteBatchUseCase(executionService, backendClient, executionRepository);
    const syncExecutionsUseCase = createSyncExecutionsUseCase(syncService, backendClient, executionRepository);

    // Orchestrators
    const executionOrchestrator = createExecutionOrchestrator(
      executeBatchUseCase,
      syncExecutionsUseCase,
      manageStagingUseCase
    );

    const eventOrchestrator = createEventOrchestrator(
      eventListener,
      sequenceManager,
      idempotencyManager,
      eventDispatcher,
      eventProcessingService,
    );

    return {
      boosterRepository,
      executionRepository,
      loadBoostersUseCase,
      manageStagingUseCase,
      executionOrchestrator,
      eventOrchestrator,
    };
  }, []);

  return (
    <BoosterContextProvider
      boosterRepository={dependencies.boosterRepository}
      executionRepository={dependencies.executionRepository}
      loadBoostersUseCase={dependencies.loadBoostersUseCase}
      manageStagingUseCase={dependencies.manageStagingUseCase}
    >
      <ExecutionProvider
        executionOrchestrator={dependencies.executionOrchestrator}
        eventOrchestrator={dependencies.eventOrchestrator}
      >
        {children}
      </ExecutionProvider>
    </BoosterContextProvider>
  );
}
