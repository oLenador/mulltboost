import { useEffect } from 'react';
import { useOptimizationStore } from '../store/optimizationStore';
import type { OptimizationCategory } from '../api/types';

export const useOptimizations = (category?: OptimizationCategory) => {
  const store = useOptimizationStore();

  useEffect(() => {
    if (store.optimizations.length === 0) {
      store.loadOptimizations();
    }
  }, []);

  const optimizations = category 
    ? store.getOptimizationsByCategory(category)
    : store.optimizations;

  return {
    optimizations,
    states: store.optimizationStates,
    isLoading: store.isLoading,
    error: store.error,
    
    // Actions
    applyOptimization: store.applyOptimization,
    revertOptimization: store.revertOptimization,
    applyBatch: store.applyBatch,
    loadOptimizations: store.loadOptimizations,
    clearError: store.clearError,
    
    // Computed
    appliedOptimizations: store.getAppliedOptimizations(),
    
    // Utils
    isApplied: (id: string) => store.optimizationStates[id]?.applied ?? false,
    canApply: (id: string) => !store.optimizationStates[id]?.applied,
    canRevert: (id: string) => store.optimizationStates[id]?.applied ?? false,
  };
};

export const usePrecisionOptimizations = () => useOptimizations('precision');
export const usePerformanceOptimizations = () => useOptimizations('performance');
export const useNetworkOptimizations = () => useOptimizations('network');
export const useSystemOptimizations = () => useOptimizations('system');