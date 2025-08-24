import { useState, useMemo, useCallback } from 'react';
import { BoosterEntity } from '../../core/entities/booster.entity';

export interface FilterState {
  searchTerm: string;
  impactFilter: string;
  statusFilter: string;
  showAdvanced: boolean;
}

export interface FilterActions {
  setSearchTerm: (term: string) => void;
  setImpactFilter: (impact: string) => void;
  setStatusFilter: (status: string) => void;
  setShowAdvanced: (show: boolean) => void;
  resetFilters: () => void;
}

export interface FilteredResults<T> {
  filteredItems: T[];
  totalCount: number;
  filteredCount: number;
  hasActiveFilters: boolean;
}

export interface UseUiFiltersHook {
  filterState: FilterState;
  filterActions: FilterActions;
  applyFilters: <T extends BoosterEntity>(items: T[]) => FilteredResults<T>;
}

const initialFilterState: FilterState = {
  searchTerm: '',
  impactFilter: 'all',
  statusFilter: 'all',
  showAdvanced: false,
};

export function useUiFilters(): UseUiFiltersHook {
  const [filterState, setFilterState] = useState<FilterState>(initialFilterState);

  const filterActions: FilterActions = useMemo(() => ({
    setSearchTerm: (term: string) => 
      setFilterState(prev => ({ ...prev, searchTerm: term })),
    
    setImpactFilter: (impact: string) => 
      setFilterState(prev => ({ ...prev, impactFilter: impact })),
    
    setStatusFilter: (status: string) => 
      setFilterState(prev => ({ ...prev, statusFilter: status })),
    
    setShowAdvanced: (show: boolean) => 
      setFilterState(prev => ({ ...prev, showAdvanced: show })),
    
    resetFilters: () => setFilterState(initialFilterState),
  }), []);

  const applyFilters = useCallback(<T extends BoosterEntity>(items: T[]): FilteredResults<T> => {
    const { searchTerm, impactFilter, statusFilter } = filterState;
    
    const filteredItems = items.filter(item => {
      // Search filter
      if (searchTerm) {
        const searchLower = searchTerm.toLowerCase();
        const matchesName = item.name.toLowerCase().includes(searchLower);
        const matchesDescription = item.description.toLowerCase().includes(searchLower);
        const matchesTags = item.tags.some(tag => 
          tag.toLowerCase().includes(searchLower)
        );
        
        if (!matchesName && !matchesDescription && !matchesTags) {
          return false;
        }
      }

      // Impact filter
      if (impactFilter !== 'all' && item.riskLevel !== impactFilter) {
        return false;
      }

      // Status filter
      if (statusFilter === 'enabled' && !item.isApplied) return false;
      if (statusFilter === 'disabled' && item.isApplied) return false;

      return true;
    });

    const hasActiveFilters = searchTerm !== '' || 
                           impactFilter !== 'all' || 
                           statusFilter !== 'all';

    return {
      filteredItems,
      totalCount: items.length,
      filteredCount: filteredItems.length,
      hasActiveFilters,
    };
  }, [filterState]);

  return {
    filterState,
    filterActions,
    applyFilters,
  };
}