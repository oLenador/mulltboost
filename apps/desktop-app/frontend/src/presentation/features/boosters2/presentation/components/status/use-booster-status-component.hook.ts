import { useCallback } from 'react';

import { useExecutionState } from '../../hooks/use-execution-state.hook';
import { useStaging } from '../../hooks/use-staging.hook';
import { useBoosterContext } from '../../providers/booster-context.provider';
import { useExecutionContext } from '../../providers/execution.provider';

export function useBoosterExecution() {

    const { executionOrchestrator } = useExecutionContext();
    const { manageStagingUseCase } = useBoosterContext();

  const {
    executions,
    executionsByStatus,
    stats,
    isExecuting,
    getExecution,
    addExecutions,
    updateExecution,
    removeExecution,
    removeExecutions,
    clearExecutions,
    setCurrentBatch,
  } = useExecutionState();

  const {
    stagedOperations,
    stagedCount,
    hasChanges,
    canExecuteBatch,
    clearStaging,
  } = useStaging();

  // Execute staged batch: adiciona execuções no estado e limpa staging.
  // NOTA: aqui apenas cria filas/entradas de execução. A rotina que processa
  // as execuções (progress, calls API, etc.) não foi fornecida nas snippets,
  // então essa função não irá "processar" as operações, apenas as enfileira.
  const executeStagedBatch = useCallback(async () => {
    const ops = stagedOperations;
    const keys = Object.keys(ops);
    if (keys.length === 0) {
      return false;
    }

    // transformar stagedOperations em formato aceito por addExecutions
    // addExecutions aceita Record<string,'apply'|'revert'>
    try {
      addExecutions(ops);
      // opcional: criar um batch id
      setCurrentBatch && setCurrentBatch(Date.now().toString());
      clearStaging();
      return true;
    } catch (err) {
      console.error('[useBoosterExecution] executeStagedBatch error', err);
      return false;
    }
  }, [stagedOperations, addExecutions, clearStaging, setCurrentBatch]);

  // Cancela execuções que estejam relacionadas com o staged (ou pelo param)
  const cancelStagedExecutions = useCallback(async () => {
    const keys = Object.keys(stagedOperations);
    if (keys.length === 0) return false;

    try {
      // remove execs se existirem
      removeExecutions && removeExecutions(keys);
      clearStaging();
      return true;
    } catch (err) {
      console.error('[useBoosterExecution] cancelStagedExecutions error', err);
      return false;
    }
  }, [stagedOperations, removeExecutions, clearStaging]);

  const resetExecution = useCallback(() => {
    // zera execuções no estado
    clearExecutions && clearExecutions();
    setCurrentBatch && setCurrentBatch(null);
  }, [clearExecutions, setCurrentBatch]);

  // expose useful aggregations
  const stagedCount = stagedCount;
  const canExecute = canExecuteBatch;

  return {
    executions,
    executionsByStatus,
    stats,
    isExecuting,
    getExecution,
    stagedCount,
    executeStagedBatch,
    cancelStagedExecutions,
    resetExecution,
    canExecute,
    hasChanges,
  };
}
