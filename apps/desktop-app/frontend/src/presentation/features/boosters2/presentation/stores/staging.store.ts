import { atom } from 'jotai';
import { StagedOperations, createStagedOperations } from '../../core/types/execution.types';

export const stagedOperationsAtom = atom<StagedOperations>(createStagedOperations());

// Computed atoms
export const stagedCountAtom = atom((get) => {
  const staged = get(stagedOperationsAtom);
  return Object.keys(staged).length;
});

export const hasChangesAtom = atom((get) => {
  const count = get(stagedCountAtom);
  return count > 0;
});

export const canExecuteBatchAtom = atom((get) => {
  const staged = get(stagedOperationsAtom);
  return Object.keys(staged).length > 0;
});

// Actions
export const stageOperationAtom = atom(
  null,
  (get, set, { boosterId, operation }: { boosterId: string; operation: 'apply' | 'revert' }) => {
    const current = get(stagedOperationsAtom);
    const updated = { ...current, [boosterId]: operation };
    set(stagedOperationsAtom, updated);
  }
);

export const unstageOperationAtom = atom(
  null,
  (get, set, boosterId: string) => {
    const current = get(stagedOperationsAtom);
    const { [boosterId]: removed, ...rest } = current;
    set(stagedOperationsAtom, rest);
  }
);

export const stageBatchOperationsAtom = atom(
  null,
  (get, set, operations: StagedOperations) => {
    set(stagedOperationsAtom, operations);
  }
);

export const clearStagingAtom = atom(
  null,
  (get, set) => {
    set(stagedOperationsAtom, createStagedOperations());
  }
);