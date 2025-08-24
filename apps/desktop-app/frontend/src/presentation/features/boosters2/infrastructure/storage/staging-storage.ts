import { StagedOperations } from '../../core/types/execution.types';

export interface StagingStorage {
  get(): StagedOperations;
  set(operations: StagedOperations): void;
  add(boosterId: string, operation: 'apply' | 'revert'): void;
  remove(boosterId: string): void;
  clear(): void;
  has(boosterId: string): boolean;
  count(): number;
}

export function createInMemoryStagingStorage(): StagingStorage {
  let operations: StagedOperations = {};

  return {
    get(): StagedOperations {
      return { ...operations };
    },

    set(newOperations: StagedOperations): void {
      operations = { ...newOperations };
    },

    add(boosterId: string, operation: 'apply' | 'revert'): void {
      operations = { ...operations, [boosterId]: operation };
    },

    remove(boosterId: string): void {
      const { [boosterId]: removed, ...rest } = operations;
      operations = rest;
    },

    clear(): void {
      operations = {};
    },

    has(boosterId: string): boolean {
      return boosterId in operations;
    },

    count(): number {
      return Object.keys(operations).length;
    },
  };
}
