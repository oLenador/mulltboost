import { useEffect, useState, useCallback } from 'react';
import { useAtom } from 'jotai';
import { BoosterItem, BoosterPageConfig } from '../types/booster.types';
import { userDataAtom } from '@/core/store/user-data.store';
import { 
  listingBoostersAtom, 
  stagedOperationsAtom, 
  stageOperationAtom, 
  clearStagingAtom, 
  initializeStagingAtom,
  toggleAllBoostersAtom,
  getFilteredBoostersAtom,
  appliedCountAtom,
  changesCountAtom,
  hasChangesAtom
} from '@/presentation/features/boosters/stores/booster-execution.store';
import { GetBoostersByCategory } from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler';
import { Language } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/services/i18n';
import { GetBoosterDto } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/dto';

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

export function useBoosters(config: BoosterPageConfig) {
  const [searchTerm, setSearchTerm] = useState('');
  const [impactFilter, setImpactFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [loading, setLoading] = useState(false);

  const [originalBoosters, setOriginalBoosters] = useAtom(listingBoostersAtom);
  const [stagedOperations] = useAtom(stagedOperationsAtom);
  const [userData] = useAtom(userDataAtom);
  const [appliedCount] = useAtom(appliedCountAtom);
  const [changesCount] = useAtom(changesCountAtom);
  const [hasChanges] = useAtom(hasChangesAtom);

  const [, stageOperation] = useAtom(stageOperationAtom);
  const [, clearStaging] = useAtom(clearStagingAtom);
  const [, initializeStaging] = useAtom(initializeStagingAtom);
  const [, toggleAllBoosters] = useAtom(toggleAllBoostersAtom);
  const [, getFilteredBoosters] = useAtom(getFilteredBoostersAtom);

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
    return () => clearStaging();
  }, [config.category, userData.language, clearStaging, setOriginalBoosters]);

  useEffect(() => {
    initializeStaging();
  }, [originalBoosters, initializeStaging]);

  const filteredBoosters = getFilteredBoosters({ searchTerm, impactFilter, statusFilter });

  const toggleBooster = useCallback((id: string) => {
    stageOperation({ boosterId: id });
  }, [stageOperation]);

  const toggleAll = useCallback((apply: boolean) => {
    toggleAllBoosters(apply);
  }, [toggleAllBoosters]);

  const resetChanges = useCallback(() => {
    clearStaging();
  }, [clearStaging]);

  return {
    boosters: originalBoosters,
    filteredBoosters,
    loading,
    searchTerm,
    setSearchTerm,
    impactFilter,
    setImpactFilter,
    statusFilter,
    setStatusFilter,
    showAdvanced,
    setShowAdvanced,
    toggleBooster,
    toggleAllBoosters: toggleAll,
    resetChanges,
    appliedCount,
    changesCount,
    hasChanges,
    originalBoosters,
    stagedOperations,
  };
}