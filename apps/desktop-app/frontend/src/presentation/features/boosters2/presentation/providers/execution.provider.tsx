import React, { createContext, useContext, useEffect, useRef, ReactNode } from 'react';
import { ExecutionOrchestrator } from '../../application/orchestrators/execution.orchestrator';
import { EventOrchestrator } from '../../application/orchestrators/event.orchestrator';

export interface ExecutionContextValue {
  executionOrchestrator: ExecutionOrchestrator;
  eventOrchestrator: EventOrchestrator;
  isInitialized: boolean;
}

const ExecutionContext = createContext<ExecutionContextValue | null>(null);

export interface ExecutionProviderProps {
  children: ReactNode;
  executionOrchestrator: ExecutionOrchestrator;
  eventOrchestrator: EventOrchestrator;
}

export function ExecutionProvider({ 
  children, 
  executionOrchestrator, 
  eventOrchestrator 
}: ExecutionProviderProps) {
  const isInitializedRef = useRef(false);

  useEffect(() => {
    if (!isInitializedRef.current) {
      try {
        // Initialize event orchestrator first
        eventOrchestrator.start();
        
        // Set up connection between event and execution orchestrators
        eventOrchestrator.setCallbacks({
          onExecutionStatusChanged: (boosterId, status, progress, error) => {
            // This would typically update a store or trigger callbacks
            console.log('Execution status changed:', { boosterId, status, progress, error });
          },
          onSyncRequired: async () => {
            await executionOrchestrator.syncWithBackend();
          },
        });

        isInitializedRef.current = true;
        console.log('[ExecutionProvider] Initialized successfully');
      } catch (error) {
        console.error('[ExecutionProvider] Failed to initialize:', error);
      }
    }

    return () => {
      if (isInitializedRef.current) {
        eventOrchestrator.stop();
        isInitializedRef.current = false;
        console.log('[ExecutionProvider] Cleaned up');
      }
    };
  }, [executionOrchestrator, eventOrchestrator]);

  const contextValue: ExecutionContextValue = {
    executionOrchestrator,
    eventOrchestrator,
    isInitialized: isInitializedRef.current,
  };

  return (
    <ExecutionContext.Provider value={contextValue}>
      {children}
    </ExecutionContext.Provider>
  );
}

export function useExecutionContext(): ExecutionContextValue {
  const context = useContext(ExecutionContext);
  if (!context) {
    throw new Error('useExecutionContext must be used within ExecutionProvider');
  }
  return context;
}
