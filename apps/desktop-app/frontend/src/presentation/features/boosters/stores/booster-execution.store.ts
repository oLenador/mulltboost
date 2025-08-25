import { atom } from 'jotai';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { BoosterExecution, ExecutionStatus, ExecutionBatch, ExecutionStats } from '../domain/booster-queue.types';
import { BoosterItem } from '../types/booster.types';

export const boosterExecutionsAtom = atom<Map<string, BoosterExecution>>(new Map());
export const currentBatchAtom = atom<ExecutionBatch | null>(null);
export const listingBoostersAtom = atom<BoosterItem[]>([]);

export type StagedOperations = Record<string, BoosterOperationType>;
export const stagedOperationsAtom = atom<StagedOperations>({});

const filterBoostersBySearchTerm = (boosters: BoosterItem[], searchTerm: string): BoosterItem[] => {
  if (!searchTerm) return boosters;
  
  const searchLower = searchTerm.toLowerCase();
  return boosters.filter(booster => {
    const matchesName = booster.name.toLowerCase().includes(searchLower);
    const matchesDescription = booster.description.toLowerCase().includes(searchLower);
    const matchesTags = booster.tags.some(tag => tag.toLowerCase().includes(searchLower));
    return matchesName || matchesDescription || matchesTags;
  });
};

const filterBoostersByImpact = (boosters: BoosterItem[], impactFilter: string): BoosterItem[] => {
  return impactFilter === 'all' ? boosters : boosters.filter(b => b.riskLevel === impactFilter);
};

const filterBoostersByStatus = (boosters: BoosterItem[], statusFilter: string): BoosterItem[] => {
  if (statusFilter === 'enabled') return boosters.filter(b => b.isApplied);
  if (statusFilter === 'disabled') return boosters.filter(b => !b.isApplied);
  return boosters;
};

const getNextOperation = (booster: BoosterItem, currentStaged?: BoosterOperationType): BoosterOperationType | null => {
  if (!currentStaged) {
    return booster.isApplied ? BoosterOperationType.RevertOperationType : BoosterOperationType.ApplyOperationType;
  }
  
  const nextAction = currentStaged === BoosterOperationType.ApplyOperationType ?  BoosterOperationType.RevertOperationType  : BoosterOperationType.ApplyOperationType;
  const wouldReturnToOriginal = (nextAction === BoosterOperationType.ApplyOperationType && booster.isApplied) || 
                                (nextAction ===  BoosterOperationType.RevertOperationType  && !booster.isApplied);
  
  return wouldReturnToOriginal ? null : nextAction;
};

const createStagingChangesForBatch = (boosters: BoosterItem[], apply: boolean): StagedOperations => {
  return boosters.reduce((acc, booster) => {
    if (booster.isApplied !== apply) {
      acc[booster.id] = apply ? BoosterOperationType.ApplyOperationType : BoosterOperationType.RevertOperationType
    }
    return acc;
  }, {} as StagedOperations);
};

export const filteredBoostersAtom = atom(
  (get) => get(listingBoostersAtom),
  (get, set, filters: { searchTerm: string; impactFilter: string; statusFilter: string }) => {
    const boosters = get(listingBoostersAtom);
    let filtered = filterBoostersBySearchTerm(boosters, filters.searchTerm);
    filtered = filterBoostersByImpact(filtered, filters.impactFilter);
    filtered = filterBoostersByStatus(filtered, filters.statusFilter);
    return filtered;
  }
);

export const executionStatsAtom = atom((get) => {
  const executions = Array.from(get(boosterExecutionsAtom).values());
  
  const stats: ExecutionStats = {
    total: executions.length,
    queued: executions.filter(e => e.status === 'queued').length,
    processing: executions.filter(e => e.status === 'processing').length,
    completed: executions.filter(e => !!e.completedAt).length,
    errors: executions.filter(e => e.status === 'error').length,
    cancelled: executions.filter(e => e.status === 'cancelled').length,
    totalProgress: executions.length > 0 
      ? executions.reduce((sum, e) => sum + e.progress, 0) / executions.length 
      : 0
  };
  
  return stats;
});

export const executionsByStatusAtom = atom((get) => {
  const executions = Array.from(get(boosterExecutionsAtom).values());
  
  return {
    idle: executions.filter(e => e.status === 'idle'),
    queued: executions.filter(e => e.status === 'queued'),
    processing: executions.filter(e => e.status === 'processing'),
    success: executions.filter(e => !!e.completedAt),
    error: executions.filter(e => e.status === 'error'),
    cancelled: executions.filter(e => e.status === 'cancelled'),
    completed: executions.filter(e => e.status === 'completed'),
  };
});

export const isExecutingAtom = atom((get) => {
  const stats = get(executionStatsAtom);
  return stats.queued > 0 || stats.processing > 0;
});

export const canExecuteBatchAtom = atom((get) => {
  const staged = get(stagedOperationsAtom);
  const isExecuting = get(isExecutingAtom);
  return Object.keys(staged).length > 0 && !isExecuting;
});

export const appliedCountAtom = atom((get) => {
  return get(listingBoostersAtom).filter(b => b.isApplied).length;
});

export const changesCountAtom = atom((get) => {
  return Object.keys(get(stagedOperationsAtom)).length;
});

export const hasChangesAtom = atom((get) => {
  return get(changesCountAtom) > 0;
});

export const updateExecutionAtom = atom(
  null,
  (get, set, update: { boosterId: string; status?: ExecutionStatus; progress?: number; error?: string }) => {
    const executions = new Map(get(boosterExecutionsAtom));
    const existing = executions.get(update.boosterId);

    if (!existing) return;

    const now = new Date();
    const isFinalStatus = ['completed', 'error', 'cancelled'].includes(update.status!);

    const updated: BoosterExecution = {
      ...existing,
      ...(update.status && { status: update.status }),
      ...(update.progress !== undefined && { progress: update.progress }),
      ...(update.error && { error: update.error }),
      canCancel: update.status === 'processing' || update.status === 'queued',
      ...(update.status === 'processing' && !existing.startedAt ? { startedAt: now } : {}),
      ...(isFinalStatus ? { completedAt: now, canCancel: false } : {}),
    };

    executions.set(update.boosterId, updated);
    set(boosterExecutionsAtom, executions);
  }
);

export const addExecutionsAtom = atom(
  null,
  (get, set, operations: StagedOperations) => {
    const executions = new Map(get(boosterExecutionsAtom));
    
    Object.entries(operations).forEach(([boosterId, operation]) => {
      const execution: BoosterExecution = {
        boosterId,
        operation,
        status: 'idle',
        progress: 0,
        canCancel: false,
      };
      executions.set(boosterId, execution);
    });
    
    set(boosterExecutionsAtom, executions);
  }
);

export const removeExecutionsAtom = atom(
  null,
  (get, set, boosterIds: string[]) => {
    const executions = new Map(get(boosterExecutionsAtom));
    
    boosterIds.forEach(id => {
      const execution = executions.get(id);
      if (execution && execution.status === 'idle') {
        executions.delete(id);
      }
    });
    
    set(boosterExecutionsAtom, executions);
  }
);

export const clearExecutionsAtom = atom(
  null,
  (get, set) => {
    set(boosterExecutionsAtom, new Map());
    set(currentBatchAtom, null);
  }
);

export const stageOperationAtom = atom(
  null,
  (get, set, { boosterId }: { boosterId: string }) => {
    const boosters = get(listingBoostersAtom);
    const staged = get(stagedOperationsAtom);
    const booster = boosters.find(b => b.id === boosterId);
    
    if (!booster) return;

    const nextOperation = getNextOperation(booster, staged[boosterId]);
    
    if (nextOperation === null) {
      const { [boosterId]: _, ...rest } = staged;
      set(stagedOperationsAtom, rest);
    } else {
      set(stagedOperationsAtom, { ...staged, [boosterId]: nextOperation });
    }
  }
);

export const stageBatchOperationsAtom = atom(
  null,
  (get, set, operations: StagedOperations) => {
    const staged = get(stagedOperationsAtom);
    set(stagedOperationsAtom, { ...staged, ...operations });
  }
);

export const clearStagingAtom = atom(
  null,
  (get, set) => {
    set(stagedOperationsAtom, {});
  }
);

export const initializeStagingAtom = atom(
  null,
  (get, set) => {
    const boosters = get(listingBoostersAtom);
    if (boosters.length === 0) return;

    const currentlyAppliedOperations = boosters
      .filter(booster => booster.isApplied)
      .reduce((acc, booster) => {
        acc[booster.id] = BoosterOperationType.ApplyOperationType;
        return acc;
      }, {} as StagedOperations);

    if (Object.keys(currentlyAppliedOperations).length > 0) {
      set(stagedOperationsAtom, currentlyAppliedOperations);
    }
  }
);

export const toggleAllBoostersAtom = atom(
  null,
  (get, set, apply: boolean) => {
    const boosters = get(listingBoostersAtom);
    const stagingChanges = createStagingChangesForBatch(boosters, apply);
    
    set(stagedOperationsAtom, {});
    if (Object.keys(stagingChanges).length > 0) {
      set(stagedOperationsAtom, stagingChanges);
    }
  }
);

export const getFilteredBoostersAtom = atom(
  null,
  (get, set, filters: { searchTerm: string; impactFilter: string; statusFilter: string }) => {
    const boosters = get(listingBoostersAtom);
    let filtered = filterBoostersBySearchTerm(boosters, filters.searchTerm);
    filtered = filterBoostersByImpact(filtered, filters.impactFilter);
    filtered = filterBoostersByStatus(filtered, filters.statusFilter);
    return filtered;
  }
);