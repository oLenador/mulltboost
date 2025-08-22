import { atom } from 'jotai';
import { ProcessableItem, ProcessableItemWithStatus, BatchOperation, ItemStatus } from '../../presentation/features/boosters/domain/batch/batch.types';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { BoosterItem } from '@/presentation/features/boosters/types/booster.types';

// Base atoms
export const batchItemsAtom = atom<ProcessableItemWithStatus[]>([]);
export const batchOperationsAtom = atom<BatchOperation[]>([]);

// BOosters Array to batch manager access
export const listingBoostersAtom = atom<BoosterItem[]>([]);


export type StagedItemsType = Record<string, BoosterOperationType>
export const stagedItemsAtom = atom<StagedItemsType>({});

export const activeTabAtom = atom<string>('all');

// Derived atoms
export const itemsByStatusAtom = atom((get) => {
  const items = get(batchItemsAtom);
  return {
    idle: items.filter(item => item.status === 'idle'),
    queued: items.filter(item => item.status === 'queued'),
    processing: items.filter(item => item.status === 'processing'),
    completed: items.filter(item => item.status === 'completed'),
    error: items.filter(item => item.status === 'error'),
    cancelled: items.filter(item => item.status === 'cancelled'),
  };
});

export const batchStatsAtom = atom((get) => {
  const itemsByStatus = get(itemsByStatusAtom);
  const total = get(batchItemsAtom).length;
  return {
    total,
    idle: itemsByStatus.idle.length,
    queued: itemsByStatus.queued.length,
    processing: itemsByStatus.processing.length,
    completed: itemsByStatus.completed.length,
    error: itemsByStatus.error.length,
    cancelled: itemsByStatus.cancelled.length,
    progress: total > 0 ? (itemsByStatus.completed.length / total) * 100 : 0
  };
});

export const canStartBatchAtom = atom((get) => {
  const itemsByStatus = get(itemsByStatusAtom);
  return itemsByStatus.idle.length > 0;
});

export const isProcessingAtom = atom((get) => {
  const itemsByStatus = get(itemsByStatusAtom);
  return itemsByStatus.processing.length > 0 || itemsByStatus.queued.length > 0;
});

// Write atoms (ações)
export const updateItemStatusAtom = atom(
  null,
  (get, set, { id, status, progress, error }: { 
    id: string; 
    status: ItemStatus; 
    progress?: number; 
    error?: string; 
  }) => {
    const items = get(batchItemsAtom);
    const updatedItems = items.map(item => 
      item.item.id === id 
        ? {
            ...item,
            status,
            progress: progress ?? item.progress,
            error,
            canCancel: status === 'processing' || status === 'queued',
            ...(status === 'processing' && !item.startedAt ? { startedAt: new Date() } : {}),
            ...(status === 'completed' || status === 'error' || status === 'cancelled' 
              ? { completedAt: new Date(), canCancel: false } 
              : {})
          }
        : item
    );
    set(batchItemsAtom, updatedItems);
  }
);

export const addBatchItemsAtom = atom(
  null,
  (get, set, items: ProcessableItem[]) => {
    const currentItems = get(batchItemsAtom);
    const newItems: ProcessableItemWithStatus[] = items.map(item => ({
      item,
      status: 'idle' as ItemStatus,
      progress: 0,
      canCancel: false,
      permissions: []
    }));
    set(batchItemsAtom, [...currentItems, ...newItems]);
  }
);

export const removeBatchItemsAtom = atom(
  null,
  (get, set, ids: string[]) => {
    const items = get(batchItemsAtom);
    const filteredItems = items.filter(item => 
      !ids.includes(item.item.id) || item.status !== 'idle'
    );
    set(batchItemsAtom, filteredItems);
  }
);

// Atom de estágio de item individual
export const itemStageAtom = atom(
  null,
  (get, set, { itemId, action }: { itemId: string; action: BoosterOperationType }) => {
    const staged = get(stagedItemsAtom);
    if (staged[itemId]) {
      // Se já existe, remove
      const { [itemId]: _, ...rest } = staged;
      set(stagedItemsAtom, rest);
    } else {
      set(stagedItemsAtom, { ...staged, [itemId]: action });
    }
  }
);

// Atom de estágio em lote
export const batchStageItemsAtom = atom(
  null,
  (get, set, items: StagedItemsType) => {
    const staged = get(stagedItemsAtom);
    set(stagedItemsAtom, { ...staged, ...items });
  }
);

export const clearStagingAtom = atom(
  null,
  (get, set) => {
    set(stagedItemsAtom, {});
  }
);
