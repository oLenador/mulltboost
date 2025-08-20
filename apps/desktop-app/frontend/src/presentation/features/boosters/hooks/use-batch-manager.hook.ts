// src/presentation/features/boosters/hooks/use-batch-manager.hook.ts

import { useEffect, useRef, useCallback } from 'react';
import { useAtom } from 'jotai';
import { BoosterBatchManager } from '../domain/booster-batch.manager';
import {
    itemsByStatusAtom,
    updateItemStatusAtom,
    batchItemsAtom,
    batchStatsAtom,
    selectedItemsAtom,
    canStartBatchAtom,
    isProcessingAtom,
    addBatchItemsAtom,
    removeBatchItemsAtom,
    toggleItemSelectionAtom,
    selectAllItemsAtom,
    clearSelectionAtom,
} from '@/core/store/batch.store';
import { BoosterItem } from '../types/booster.types';
import { useBoosterBatchManagerContext } from '../providers/batch-managers.provider';


export const useBatchManager = () => {
    const { manager } = useBoosterBatchManagerContext();
    const [items] = useAtom(batchItemsAtom);
    const [stats] = useAtom(batchStatsAtom);
    const [selectedItems] = useAtom(selectedItemsAtom);
    const [canStart] = useAtom(canStartBatchAtom);
    const [isProcessing] = useAtom(isProcessingAtom);
    const [itemsByStatus] = useAtom(itemsByStatusAtom);
    
    // Actions
    const [, updateItemStatus] = useAtom(updateItemStatusAtom);
    const [, addItems] = useAtom(addBatchItemsAtom);
    const [, removeItems] = useAtom(removeBatchItemsAtom);
    const [, toggleSelection] = useAtom(toggleItemSelectionAtom);
    const [, selectAll] = useAtom(selectAllItemsAtom);
    const [, clearSelection] = useAtom(clearSelectionAtom);

    // Initialize manager
    useEffect(() => {
        if (!manager) return;
    
        // Subscribe to manager updates
        const unsubscribe = manager.subscribe((managedItems) => {
          // Sync manager state with Jotai store without causing infinite loops
          const currentItemIds = items.map(item => item.item.id);
          const newItems = managedItems.filter(item => !currentItemIds.includes(item.item.id));
          
          if (newItems.length > 0) {
            addItems(newItems.map(item => item.item));
          }
        });
    
        // Setup event listeners
        manager.on('onItemStatusChanged', (item) => {
          updateItemStatus({
            id: item.item.id,
            status: item.status,
            progress: item.progress,
            error: item.error
          });
        });
    
        manager.on('onBatchStarted', (operation) => {
          console.log('Batch started:', operation);
        });
    
        manager.on('onBatchCompleted', (operation) => {
          console.log('Batch completed:', operation);
        });
    
        manager.on('onBatchError', (operation, error) => {
          console.error('Batch error:', error);
        });
    
        return unsubscribe;
      }, [manager, addItems, updateItemStatus, items]);
    
      const addBoosterItems = useCallback((boosters: BoosterItem[]) => {
        if (!manager) return;
        manager.addItems(boosters);
      }, [manager]);
    
      const removeBoosterItems = useCallback((ids: string[]) => {
        if (!manager) return;
        manager.removeItems(ids);
        removeItems(ids);
      }, [manager, removeItems]);
    
      const startBatch = useCallback(async (itemIds?: string[]) => {
        if (!manager) return;
        
        try {
          const operationId = await manager.startBatch(itemIds);
          return operationId;
        } catch (error) {
          console.error('Failed to start batch:', error);
          throw error;
        }
      }, [manager]);
    
      const cancelItems = useCallback((ids: string[]) => {
        if (!manager) return;
        manager.cancelItems(ids);
      }, [manager]);
    
      const startSelectedBatch = useCallback(async () => {
        if (selectedItems.length > 0) {
          return startBatch(selectedItems);
        }
      }, [selectedItems, startBatch]);
    
      const cancelSelected = useCallback(() => {
        if (selectedItems.length > 0) {
          cancelItems(selectedItems);
        }
      }, [selectedItems, cancelItems]);
    
      const syncWithBackend = useCallback(async () => {
        if (!manager) return;
        
        try {
          await manager.syncWithBackend();
        } catch (error) {
          console.error('Failed to sync with backend:', error);
        }
      }, [manager]);
    
    
      return {
        // State
        items,
        stats,
        selectedItems,
        canStart,
        isProcessing,
        itemsByStatus,
        
        addBoosterItems,
        removeBoosterItems,
        startBatch,
        startSelectedBatch,
        cancelItems,
        cancelSelected,
        syncWithBackend,
        
        // Selection
        toggleSelection,
        selectAll,
        clearSelection,
        
        manager
      };
    };