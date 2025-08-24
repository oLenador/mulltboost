export interface IdempotencyState {
    processedIds: Set<string>;
    lastCleanup: number;
}

export interface IdempotencyManager {
    hasProcessed(idempotencyId: string, state: IdempotencyState): boolean;
    markAsProcessed(idempotencyId: string, state: IdempotencyState): void;
    cleanup(state: IdempotencyState, retentionTime: number): void;
}

export function createIdempotencyManager(): IdempotencyManager {
    return {
        hasProcessed(idempotencyId: string, state: IdempotencyState): boolean {
            return state.processedIds.has(idempotencyId);
        },

        markAsProcessed(idempotencyId: string, state: IdempotencyState): void {
            state.processedIds.add(idempotencyId);
        },

        cleanup(state: IdempotencyState, retentionTime: number): void {
            const now = Date.now();
            if (now - state.lastCleanup < retentionTime / 10) return; // Cleanup every 10% of retention time

            // For simplicity, we clear all. In production, you'd track timestamps per ID
            if (now - state.lastCleanup > retentionTime) {
                state.processedIds.clear();
                state.lastCleanup = now;
            }
        },
    };
}

export function createIdempotencyState(): IdempotencyState {
    return {
        processedIds: new Set(),
        lastCleanup: Date.now(),
    };
}
