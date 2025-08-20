import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { apiClient } from '../api/client';
import type { 
  Optimization, 
  OptimizationState, 
  OptimizationResult, 
  BatchResult, 
  OptimizationCategory 
} from '../api/types';

interface OptimizationStore {
  // State
  optimizations: Optimization[];
  optimizationStates: Record<string, OptimizationState>;
  isLoading: boolean;
  error: string | null;
  
  // Computed
  getOptimizationsByCategory: (category: OptimizationCategory) => Optimization[];
  getAppliedOptimizations: () => Optimization[];
  
  // Actions
  loadOptimizations: () => Promise<void>;
  loadOptimizationState: (id: string) => Promise<void>;
  applyOptimization: (id: string) => Promise<OptimizationResult>;
  revertOptimization: (id: string) => Promise<OptimizationResult>;
  applyBatch: (ids: string[]) => Promise<BatchResult>;
  clearError: () => void;
}

export const useOptimizationStore = create<OptimizationStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      optimizations: [],
      optimizationStates: {},
      isLoading: false,
      error: null,

      // Computed selectors
      getOptimizationsByCategory: (category: OptimizationCategory) => {
        return get().optimizations.filter(opt => opt.category === category);
      },

      getAppliedOptimizations: () => {
        const { optimizations, optimizationStates } = get();
        return optimizations.filter(opt => 
          optimizationStates[opt.id]?.applied === true
        );
      },

      // Actions
      loadOptimizations: async () => {
        set({ isLoading: true, error: null });
        try {
          const optimizations = await apiClient.getAvailableOptimizations();
          
          // Load states for all optimizations
          const states: Record<string, OptimizationState> = {};
          await Promise.allSettled(
            optimizations.map(async (opt) => {
              try {
                const state = await apiClient.getOptimizationState(opt.id);
                states[opt.id] = state;
              } catch (error) {
                // If state doesn't exist, create default
                states[opt.id] = {
                  id: opt.id,
                  applied: false,
                  version: opt.version,
                  backupData: {},
                  status: 'not_applied'
                };
              }
            })
          );

          set({ 
            optimizations, 
            optimizationStates: states,
            isLoading: false 
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Unknown error',
            isLoading: false 
          });
        }
      },

      loadOptimizationState: async (id: string) => {
        try {
          const state = await apiClient.getOptimizationState(id);
          set(state => ({
            optimizationStates: {
              ...state.optimizationStates,
              [id]: state
            }
          }));
        } catch (error) {
          console.warn(`Failed to load state for ${id}:`, error);
        }
      },

      applyOptimization: async (id: string) => {
        set({ isLoading: true, error: null });
        try {
          const result = await apiClient.applyOptimization(id);
          
          // Update local state
          if (result.success) {
            set(state => ({
              optimizationStates: {
                ...state.optimizationStates,
                [id]: {
                  ...state.optimizationStates[id],
                  applied: true,
                  appliedAt: new Date().toISOString(),
                  status: 'applied',
                  backupData: result.backupData || {}
                }
              }
            }));
          }
          
          set({ isLoading: false });
          return result;
        } catch (error) {
          const errorMsg = error instanceof Error ? error.message : 'Unknown error';
          set({ error: errorMsg, isLoading: false });
          throw error;
        }
      },

      revertOptimization: async (id: string) => {
        set({ isLoading: true, error: null });
        try {
          const result = await apiClient.revertOptimization(id);
          
          // Update local state
          if (result.success) {
            set(state => ({
              optimizationStates: {
                ...state.optimizationStates,
                [id]: {
                  ...state.optimizationStates[id],
                  applied: false,
                  revertedAt: new Date().toISOString(),
                  status: 'reverted'
                }
              }
            }));
          }
          
          set({ isLoading: false });
          return result;
        } catch (error) {
          const errorMsg = error instanceof Error ? error.message : 'Unknown error';
          set({ error: errorMsg, isLoading: false });
          throw error;
        }
      },

      applyBatch: async (ids: string[]) => {
        set({ isLoading: true, error: null });
        try {
          const result = await apiClient.applyOptimizationBatch(ids);
          
          // Update states based on results
          const updatedStates = { ...get().optimizationStates };
          Object.entries(result.results).forEach(([id, optResult]) => {
            if (optResult.success) {
              updatedStates[id] = {
                ...updatedStates[id],
                applied: true,
                appliedAt: new Date().toISOString(),
                status: 'applied',
                backupData: optResult.backupData || {}
              };
            }
          });
          
          set({ 
            optimizationStates: updatedStates,
            isLoading: false 
          });
          return result;
        } catch (error) {
          const errorMsg = error instanceof Error ? error.message : 'Unknown error';
          set({ error: errorMsg, isLoading: false });
          throw error;
        }
      },

      clearError: () => set({ error: null })
    }),
    {
      name: 'optimization-store'
    }
  )
);