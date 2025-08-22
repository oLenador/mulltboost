import { useEffect, useState, useMemo, useCallback } from 'react';
import { useAtom } from 'jotai';
import { BoosterItem, BoosterPageConfig } from '../types/booster.types';
import { userDataAtom } from '@/core/store/user-data.store';

import { GetBoostersByCategory } from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler';
import { Language } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/services/i18n';
import { GetBoosterDto } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/dto';
import { 
  stagedItemsAtom, 
  itemStageAtom, 
  batchStageItemsAtom, 
  clearStagingAtom, 
  StagedItemsType,
  listingBoostersAtom
} from '@/core/store/batch.store';
import { BoosterOperationType } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';



type BoosterState = 'original' | 'staged-apply' | 'staged-revert';
type EffectiveBoosterItem = BoosterItem & {
  effectiveState: BoosterState;
  hasChanges: boolean;
};


function mapRiskLevelToImpact(riskLevel: string): 'low' | 'medium' | 'high' {
  const normalizedLevel = (riskLevel || '').toLowerCase();
  return ['low', 'medium', 'high'].includes(normalizedLevel) 
    ? normalizedLevel as 'low' | 'medium' | 'high'
    : 'low';
}

function mapBoosterToItem(dto: GetBoosterDto): BoosterItem {
  return {
    id: dto.ID,
    name: dto.Name,
    description: dto.Description,
    category: dto.Category,
    level: dto.Level,
    platform: dto.Platform.join(', '), 
    dependencies: dto.Dependencies || [],
    conflicts: dto.Conflicts || [],
    reversible: dto.Reversible,
    riskLevel: mapRiskLevelToImpact(dto.RiskLevel),
    version: dto.Version,
    isApplied: dto.IsApplied,
    appliedAt: dto.AppliedAt ? new Date(dto.AppliedAt.toString()) : new Date(0),
    revertedAt: dto.RevertedAt ? new Date(dto.RevertedAt.toString()) : new Date(0),
    tags: dto.Tags || [],
  };
}

function calculateEffectiveState(
  originalBooster: BoosterItem, 
  stagedAction?: BoosterOperationType
): { isApplied: boolean; effectiveState: BoosterState; hasChanges: boolean } {
  const wasOriginallyApplied = originalBooster.isApplied;
  
  if (!stagedAction) {
    return {
      isApplied: wasOriginallyApplied,
      effectiveState: 'original',
      hasChanges: false
    };
  }

  const effectiveIsApplied = stagedAction === 'apply';
  const hasChanges = wasOriginallyApplied !== effectiveIsApplied;

  return {
    isApplied: effectiveIsApplied,
    effectiveState: hasChanges ? `staged-${stagedAction}` : 'original',
    hasChanges
  };
}

function createEffectiveBooster(
  originalBooster: BoosterItem,
  stagedAction?: BoosterOperationType
): EffectiveBoosterItem {
  const { isApplied, effectiveState, hasChanges } = calculateEffectiveState(
    originalBooster, 
    stagedAction
  );

  return {
    ...originalBooster,
    isApplied,
    effectiveState,
    hasChanges
  };
}

export function useBoosters(config: BoosterPageConfig) {
  // UI State
  const [searchTerm, setSearchTerm] = useState('');
  const [impactFilter, setImpactFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [showAdvanced, setShowAdvanced] = useState(false);
  
  // Data State - Keep original state immutable
  const [originalBoosters, setOriginalBoosters] = useAtom(listingBoostersAtom);
  const [loading, setLoading] = useState(false);

  // Jotai atoms
  const [userData] = useAtom(userDataAtom);
  const [stagedItems] = useAtom(stagedItemsAtom);
  const [, stageItem] = useAtom(itemStageAtom);
  const [, batchStageItems] = useAtom(batchStageItemsAtom);
  const [, clearStaging] = useAtom(clearStagingAtom);

  // Load boosters data
  useEffect(() => {
    const loadBoosters = async () => {
      setLoading(true);
      try {
        const response = await GetBoostersByCategory(
          config.category, 
          userData.language as Language
        );
        const mappedBoosters = response.map(mapBoosterToItem);
        setOriginalBoosters(mappedBoosters);
      } catch (error) {
        console.error('Error loading boosters:', error);
        setOriginalBoosters([]);
      } finally {
        setLoading(false);
      }
    };

    loadBoosters();

    // Cleanup staged items when component unmounts or category changes
    return () => clearStaging();
  }, [config.category, userData.language, clearStaging]);

  // Initialize staging with currently applied boosters
  useEffect(() => {
    if (originalBoosters.length === 0) return;

    const currentlyAppliedIds = originalBoosters
      .filter(booster => booster.isApplied)
      .reduce((acc, booster) => {
        acc[booster.id] = 'apply' as BoosterOperationType;
        return acc;
      }, {} as StagedItemsType);

    if (Object.keys(currentlyAppliedIds).length > 0) {
      batchStageItems(currentlyAppliedIds);
    }
  }, [originalBoosters, batchStageItems]);

  // Calculate effective boosters (original + staged changes)
  const effectiveBoosters = useMemo<EffectiveBoosterItem[]>(() => {
    return originalBoosters.map(booster => 
      createEffectiveBooster(booster, stagedItems[booster.id])
    );
  }, [originalBoosters, stagedItems]);

  // Filtered boosters for display
  const filteredBoosters = useMemo(() => {
    return effectiveBoosters.filter(booster => {
      // Search filter
      if (searchTerm) {
        const searchLower = searchTerm.toLowerCase();
        const matchesName = booster.name.toLowerCase().includes(searchLower);
        const matchesDescription = booster.description.toLowerCase().includes(searchLower);
        const matchesTags = booster.tags.some(tag => 
          tag.toLowerCase().includes(searchLower)
        );
        
        if (!matchesName && !matchesDescription && !matchesTags) {
          return false;
        }
      }

      // Impact filter
      if (impactFilter !== 'all' && booster.riskLevel !== impactFilter) {
        return false;
      }

      // Status filter
      if (statusFilter === 'enabled' && !booster.isApplied) return false;
      if (statusFilter === 'disabled' && booster.isApplied) return false;

      // Advanced filter (if implemented)
      // Could filter by effectiveState, hasChanges, etc.

      return true;
    });
  }, [effectiveBoosters, searchTerm, impactFilter, statusFilter, showAdvanced]);

  // Action handlers
  const toggleBooster = useCallback((id: string) => {
    const originalBooster = originalBoosters.find(b => b.id === id);
    if (!originalBooster) return;

    const currentStaged = stagedItems[id];
    const wasOriginallyApplied = originalBooster.isApplied;

    // Determine the action to stage
    let actionToStage: BoosterOperationType;

    if (!currentStaged) {
      // No staging exists, toggle from original state
      actionToStage = wasOriginallyApplied ? BoosterOperationType.RevertOperationType : BoosterOperationType.ApplyOperationType;
    } else {
      // Already staged, toggle the staged action
      actionToStage = currentStaged === 'apply' ? BoosterOperationType.RevertOperationType : BoosterOperationType.ApplyOperationType;
      
      // If toggling back to original state, remove from staging
      if ((actionToStage === 'apply' && wasOriginallyApplied) ||
          (actionToStage === 'revert' && !wasOriginallyApplied)) {
        stageItem({ itemId: id, action: actionToStage }); // This will remove it
        return;
      }
    }

    stageItem({ itemId: id, action: actionToStage as BoosterOperationType });
  }, [originalBoosters, stagedItems, stageItem]);


  const toggleAllBoosters = useCallback((apply: boolean) => {
    const stagingChanges: StagedItemsType = {};
  
    for (const booster of originalBoosters) {
      const shouldStage = booster.isApplied !== apply;
      if (shouldStage) {
        stagingChanges[booster.id] = apply ? BoosterOperationType.ApplyOperationType : BoosterOperationType.RevertOperationType 
      }
    }
  
    clearStaging();
    if (Object.keys(stagingChanges).length > 0) {
      batchStageItems(stagingChanges);
    }
  }, [originalBoosters, clearStaging, batchStageItems]);

  const resetChanges = useCallback(() => {
    clearStaging();
  }, [clearStaging]);

  // Computed values
  const appliedCount = useMemo(() => 
    effectiveBoosters.filter(b => b.isApplied).length, 
    [effectiveBoosters]
  );

  const changesCount = useMemo(() =>
    effectiveBoosters.filter(b => b.hasChanges).length,
    [effectiveBoosters]
  );

  const hasChanges = changesCount > 0;

  return {
    // Data
    boosters: effectiveBoosters,
    filteredBoosters,
    loading,
    
    // UI State
    searchTerm,
    setSearchTerm,
    impactFilter,
    setImpactFilter,
    statusFilter,
    setStatusFilter,
    showAdvanced,
    setShowAdvanced,
    
    // Actions
    toggleBooster,
    toggleAllBoosters,
    resetChanges,
    
    // Computed values
    appliedCount,
    changesCount,
    hasChanges,
    
    originalBoosters,
    stagedItems,
  };
}