import { useCallback } from 'react';
import { useAtom } from 'jotai';
import { StagedOperations } from '../../core/types/execution.types';
import {
  stagedOperationsAtom,
  stagedCountAtom,
  hasChangesAtom,
  canExecuteBatchAtom,
  stageOperationAtom,
  unstageOperationAtom,
  stageBatchOperationsAtom,
  clearStagingAtom,
} from '../stores/staging.store';

export const useStaging = () => {
  const [stagedOperations] = useAtom(stagedOperationsAtom);
  const [stagedCount] = useAtom(stagedCountAtom);
  const [hasChanges] = useAtom(hasChangesAtom);
  const [canExecuteBatch] = useAtom(canExecuteBatchAtom);

  // Actions
  const [, stageOperation] = useAtom(stageOperationAtom);
  const [, unstageOperation] = useAtom(unstageOperationAtom);
  const [, stageBatchOperations] = useAtom(stageBatchOperationsAtom);
  const [, clearStaging] = useAtom(clearStagingAtom);

  const isStaged = useCallback((boosterId: string): boolean => {
    return boosterId in stagedOperations;
  }, [stagedOperations]);

  const getStagedOperation = useCallback((boosterId: string): 'apply' | 'revert' | undefined => {
    return stagedOperations[boosterId];
  }, [stagedOperations]);

  const toggleStaging = useCallback((boosterId: string, operation: 'apply' | 'revert') => {
    if (isStaged(boosterId)) {
      unstageOperation(boosterId);
    } else {
      stageOperation({ boosterId, operation });
    }
  }, [isStaged, stageOperation, unstageOperation]);

  return {
    stagedOperations,
    stagedCount,
    hasChanges,
    canExecuteBatch,
    isStaged,
    getStagedOperation,
    stageOperation: ({ boosterId, operation }: { boosterId: string; operation: 'apply' | 'revert' }) => 
      stageOperation({ boosterId, operation }),
    unstageOperation,
    toggleStaging,
    stageBatchOperations,
    clearStaging,
  };
};