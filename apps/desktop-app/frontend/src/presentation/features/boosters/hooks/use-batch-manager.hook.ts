// src/presentation/features/boosters/hooks/use-batch-manager.hook.ts

import { useEffect, useRef, useCallback } from 'react';
import { useAtom } from 'jotai';
import { BoosterBatchManager } from '../domain/booster-batch.manager';
import {
    itemsByStatusAtom,
    updateItemStatusAtom,
    batchItemsAtom,
    batchStatsAtom,
    stagedItemsAtom,
    canStartBatchAtom,
    isProcessingAtom,
    addBatchItemsAtom,
    removeBatchItemsAtom,
    batchStageItemsAtom,
    itemStageAtom,
    clearStagingAtom,
    StagedItemsType,
} from '@/core/store/batch.store';
import { BoosterItem } from '../types/booster.types';
import { useBoosterBatchManagerContext } from '../providers/batch-managers.provider';

export const useBatchManager = () => {
    const { manager } = useBoosterBatchManagerContext();
    const [items] = useAtom(batchItemsAtom);
    const [stats] = useAtom(batchStatsAtom);
    const [stagedItems, setStagedItems] = useAtom(stagedItemsAtom);
    const [canStart] = useAtom(canStartBatchAtom);
    const [isProcessing] = useAtom(isProcessingAtom);
    const [itemsByStatus] = useAtom(itemsByStatusAtom);

    // Actions
    const [, updateItemStatus] = useAtom(updateItemStatusAtom);
    const [, addItems] = useAtom(addBatchItemsAtom);
    const [, removeItems] = useAtom(removeBatchItemsAtom);
    const [, itemStage] = useAtom(itemStageAtom);
    const [, batchStageItems] = useAtom(batchStageItemsAtom);
    const [, clearStaging] = useAtom(clearStagingAtom);

    useEffect(() => {
        if (!manager) return;

        const unsubscribe = manager.subscribe((managedItems) => {
            const currentItemIds = items.map(item => item.item.id);
            const newItems = managedItems.filter(item => !currentItemIds.includes(item.item.id));
            if (newItems.length > 0) {
                addItems(newItems.map(item => item.item));
            }
        });

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
    }, [manager, updateItemStatus, items]);

    const removeBoosterItems = useCallback((ids: string[]) => {
        if (!manager) return;
        manager.removeItems(ids);
        removeItems(ids);
    }, [manager, removeItems]);

    const startBatch = useCallback(async (items: StagedItemsType) => {
        if (!manager) return;
        try {
            const operationId = await manager.startBatch(items);
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

    const startStagedBatch = useCallback(async () => {
        if (Object.keys(stagedItems).length > 0) {
            return startBatch();
        }
    }, [stagedItems, startBatch]);
    

    const cancelStaged = useCallback(() => {
        if (Object.keys(stagedItems).length > 0) {
            cancelItems(Object.keys(stagedItems));
        }
    }, [stagedItems, cancelItems]);

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
        stagedItems,
        canStart,
        isProcessing,
        itemsByStatus,

        removeBoosterItems,
        startBatch,
        startStagedBatch,
        cancelItems,
        cancelStaged,
        syncWithBackend,

        itemStage,
        batchStageItems,
        clearStaging,

        manager
    };
};
