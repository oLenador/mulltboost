// src/presentation/hooks/use-real-time-updates.hook.ts
import { useContext, useEffect, useRef, useMemo } from 'react';
import { useAtomValue, useSetAtom } from 'jotai';
import { ExecutionContext } from '../providers/execution.provider';
import { EventBridge } from '../bridges/event-bridge';
import { createExecutionStoreAdapter, ExecutionStoreAtoms } from '../stores/adapters/execution-store.adapter';
import { createBoosterStoreAdapter, BoosterStoreAtoms } from '../stores/adapters/booster-store.adapter';
import {
  setExecutionAtom,
  updateExecutionAtom,
  removeExecutionAtom,
  clearExecutionsAtom,
  setCurrentBatchAtom,
} from '../stores/execution-state.store';
import {
  setBoostersAtom,
  updateBoosterAtom,
  clearBoostersAtom,
} from '../stores/booster-data.store';

export interface RealTimeUpdatesState {
  isConnected: boolean;
  stats: {
    totalProcessed: number;
    totalPending: number;
    stagedCount: number;
    batchActive: boolean;
  };
}

export const useRealTimeUpdates = () => {
  const executionContext = useContext(ExecutionContext);
  const isInitialized = useRef(false);
  const eventBridge = useRef<EventBridge | null>(null);

  // Atom setters
  const setExecution = useSetAtom(setExecutionAtom);
  const updateExecution = useSetAtom(updateExecutionAtom);
  const removeExecution = useSetAtom(removeExecutionAtom);
  const clearExecutions = useSetAtom(clearExecutionsAtom);
  const setCurrentBatch = useSetAtom(setCurrentBatchAtom);
  
  const setBoosters = useSetAtom(setBoostersAtom);
  const updateBooster = useSetAtom(updateBoosterAtom);
  const clearBoosters = useSetAtom(clearBoostersAtom);

  // Create store adapters
  const executionAdapter = useMemo(() => {
    const atoms: ExecutionStoreAtoms = {
      setExecutionAtom,
      updateExecutionAtom,
      removeExecutionAtom,
      clearExecutionsAtom,
      setCurrentBatchAtom,
    };

    return createExecutionStoreAdapter(
      atoms,
      // getAtomValue function - using closures to capture current setters
      (atom) => {
        if (atom === setExecutionAtom) return setExecution;
        if (atom === updateExecutionAtom) return updateExecution;
        if (atom === removeExecutionAtom) return removeExecution;
        if (atom === clearExecutionsAtom) return clearExecutions;
        if (atom === setCurrentBatchAtom) return setCurrentBatch;
        throw new Error('Unknown atom');
      },
      // setAtomValue function - not needed for this pattern
      () => {}
    );
  }, [setExecution, updateExecution, removeExecution, clearExecutions, setCurrentBatch]);

  const boosterAdapter = useMemo(() => {
    const atoms: BoosterStoreAtoms = {
      setBoostersAtom,
      updateBoosterAtom,
      clearBoostersAtom,
    };

    return createBoosterStoreAdapter(
      atoms,
      // getAtomValue function
      (atom) => {
        if (atom === setBoostersAtom) return setBoosters;
        if (atom === updateBoosterAtom) return updateBooster;
        if (atom === clearBoostersAtom) return clearBoosters;
        throw new Error('Unknown atom');
      },
      // setAtomValue function - not needed for this pattern
      () => {}
    );
  }, [setBoosters, updateBooster, clearBoosters]);

  // Initialize real-time connection
  useEffect(() => {
    if (!executionContext || isInitialized.current || !eventBridge.current) {
      return;
    }

    try {
      // Connect adapters to event bridge
      eventBridge.current.connect(executionAdapter, boosterAdapter);
      isInitialized.current = true;
      
      console.log('Real-time updates initialized');
    } catch (error) {
      console.error('Failed to initialize real-time updates:', error);
    }

    // Cleanup on unmount
    return () => {
      if (eventBridge.current && isInitialized.current) {
        eventBridge.current.disconnect();
        isInitialized.current = false;
        console.log('Real-time updates disconnected');
      }
    };
  }, [executionContext, executionAdapter, boosterAdapter]);

  // Set event bridge from context
  useEffect(() => {
    if (executionContext?.eventBridge) {
      eventBridge.current = executionContext.eventBridge;
    }
  }, [executionContext?.eventBridge]);

  // Get current state
  const state: RealTimeUpdatesState = useMemo(() => {
    if (!eventBridge.current || !isInitialized.current) {
      return {
        isConnected: false,
        stats: {
          totalProcessed: 0,
          totalPending: 0,
          stagedCount: 0,
          batchActive: false,
        },
      };
    }

    const stats = eventBridge.current.getStats();
    
    return {
      isConnected: eventBridge.current.isConnected(),
      stats: {
        totalProcessed: stats.event?.totalProcessed || 0,
        totalPending: stats.event?.totalPending || 0,
        stagedCount: stats.staged?.count || 0,
        batchActive: !!stats.batch,
      },
    };
  }, []);

  // Operations
  const operations = useMemo(() => {
    if (!eventBridge.current) {
      return {
        stageOperation: () => {},
        stageBatch: () => {},
        clearStaging: () => {},
        executeStagedBatch: async () => ({ success: false, error: 'Not connected' }),
        executeBatch: async () => ({ success: false, error: 'Not connected' }),
        getStagedOperations: () => ({}),
        getStagedCount: () => 0,
      };
    }

    return {
      stageOperation: eventBridge.current.stageOperation,
      stageBatch: eventBridge.current.stageBatch,
      clearStaging: eventBridge.current.clearStaging,
      executeStagedBatch: eventBridge.current.executeStagedBatch,
      executeBatch: eventBridge.current.executeBatch,
      getStagedOperations: eventBridge.current.getStagedOperations,
      getStagedCount: eventBridge.current.getStagedCount,
    };