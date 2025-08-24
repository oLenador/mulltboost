import { useMemo, useCallback } from 'react';

// Ajuste os imports abaixo conforme sua estrutura real de arquivos.
// Os hooks de baixo nível foram usados conforme as suas snippets.
import { useBoosterData } from '../../hooks/use-booster-data.hook';
import { useUiFilters } from '../../hooks/use-ui-filters.hook';
import { useStaging } from '../../hooks/use-staging.hook';
import { useExecutionState } from '../../hooks/use-execution-state.hook';
import { BoosterPageConfig } from '../../../types/booster.types';

export function useBoosters(config: BoosterPageConfig) {
  const { boosters, boosterById, setBoosters, updateBooster, clearBoosters } = useBoosterData();
  const { filterState, filterActions, applyFilters } = useUiFilters();
  const {
    stagedOperations,
    stagedCount,
    hasChanges,
    canExecuteBatch,
    isStaged,
    getStagedOperation,
    stageOperation,
    unstageOperation,
    toggleStaging,
    stageBatchOperations,
    clearStaging,
  } = useStaging();

  const {
    executions, // array
    executionsByStatus,
    stats,
    isExecuting,
    getExecution,
    addExecutions,
    removeExecutions,
    clearExecutions,
  } = useExecutionState();

  // Filters application
  const filteredResult = useMemo(() => applyFilters(boosters), [applyFilters, boosters]);
  const filteredBoosters = filteredResult.filteredItems;
  const totalCount = boosters.length;

  // Localized accessors for filter state/actions
  const { searchTerm, impactFilter, statusFilter, showAdvanced } = filterState;
  const { setSearchTerm, setImpactFilter, setStatusFilter, setShowAdvanced, resetFilters } = filterActions;

  // Basic computed values
  const appliedCount = useMemo(() => boosters.filter(b => b.isApplied).length, [boosters]);
  const changesCount = stagedCount; // number of staged changes (from staging store)

  // Toggle a single booster (stage apply/revert)
  const toggleBooster = useCallback((boosterIdOrObj: string | { id: string; isApplied?: boolean }) => {
    // accept either booster id or booster object
    let id: string;
    let isApplied: boolean | undefined;

    if (typeof boosterIdOrObj === 'string') {
      id = boosterIdOrObj;
      const b = boosterById.get(id);
      isApplied = b?.isApplied;
    } else {
      id = boosterIdOrObj.id;
      isApplied = boosterIdOrObj.isApplied;
    }

    const operation = isApplied ? 'revert' : 'apply';
    toggleStaging(id, operation);
  }, [boosterById, toggleStaging]);

  // Reset only staged changes (UI-level reset). resetExecution deve ser chamado separadamente.
  const resetChanges = useCallback(() => {
    clearStaging();
  }, [clearStaging]);

  // Loading heuristic: substitua se tiver um atom / flag real de loading
  const loading = boosters.length === 0;

  return {
    // data
    filteredBoosters,
    boosters,
    // filters / UI
    searchTerm,
    setSearchTerm,
    impactFilter,
    setImpactFilter,
    statusFilter,
    setStatusFilter,
    showAdvanced,
    setShowAdvanced,
    resetFilters,
    // staging / changes
    toggleBooster,
    resetChanges,
    appliedCount,
    changesCount,
    hasChanges,
    stagedCount,
    canExecute: canExecuteBatch,
    // execution / misc
    executions,
    executionsByStatus,
    stats,
    isExecuting,
    loading,
    // low-level actions exposed (se estiverem necessários)
    stageOperation,
    unstageOperation,
    getStagedOperation,
  };
}
