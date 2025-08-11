import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { apiClient } from '../api/client';
import type { SystemMetrics } from '../api/types';

interface MonitoringStore {
  // State
  currentMetrics: SystemMetrics | null;
  metricsHistory: SystemMetrics[];
  isMonitoring: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Settings
  updateInterval: number;
  maxHistorySize: number;
  
  // Actions
  startMonitoring: (interval?: number) => Promise<void>;
  stopMonitoring: () => Promise<void>;
  getMetrics: () => Promise<void>;
  clearHistory: () => void;
  setUpdateInterval: (interval: number) => void;
  clearError: () => void;
}

export const useMonitoringStore = create<MonitoringStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      currentMetrics: null,
      metricsHistory: [],
      isMonitoring: false,
      isLoading: false,
      error: null,
      updateInterval: 1000, // 1 second
      maxHistorySize: 100,

      // Actions
      startMonitoring: async (interval?: number) => {
        const { updateInterval: currentInterval } = get();
        const intervalToUse = interval ?? currentInterval;
        
        set({ isLoading: true, error: null });
        try {
          await apiClient.startRealTimeMonitoring(Math.floor(intervalToUse / 1000));
          
          set({ 
            isMonitoring: true, 
            isLoading: false,
            updateInterval: intervalToUse
          });
          
          // Start local polling for UI updates
          get().startLocalPolling();
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to start monitoring',
            isLoading: false 
          });
        }
      },

      stopMonitoring: async () => {
        set({ isLoading: true });
        try {
          await apiClient.stopRealTimeMonitoring();
          set({ 
            isMonitoring: false, 
            isLoading: false 
          });
          
          get().stopLocalPolling();
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to stop monitoring',
            isLoading: false 
          });
        }
      },

      getMetrics: async () => {
        try {
          const metrics = await apiClient.getSystemMetrics();
          const { metricsHistory, maxHistorySize } = get();
          
          // Add to history
          const newHistory = [...metricsHistory, metrics].slice(-maxHistorySize);
          
          set({ 
            currentMetrics: metrics,
            metricsHistory: newHistory,
            error: null
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to get metrics'
          });
        }
      },

      clearHistory: () => {
        set({ metricsHistory: [] });
      },

      setUpdateInterval: (interval: number) => {
        set({ updateInterval: interval });
      },

      clearError: () => set({ error: null }),

      // Private methods (not exposed in interface)
      pollingInterval: null as NodeJS.Timeout | null,
      
      startLocalPolling: () => {
        const state = get() as any;
        if (state.pollingInterval) {
          clearInterval(state.pollingInterval);
        }
        
        const interval = setInterval(() => {
          get().getMetrics();
        }, get().updateInterval);
        
        // Store interval reference
        (get() as any).pollingInterval = interval;
      },

      stopLocalPolling: () => {
        const state = get() as any;
        if (state.pollingInterval) {
          clearInterval(state.pollingInterval);
          state.pollingInterval = null;
        }
      }
    }),
    {
      name: 'monitoring-store'
    }
  )
);