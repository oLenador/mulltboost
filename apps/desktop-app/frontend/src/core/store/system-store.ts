import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { apiClient } from '../api/client';
import type { SystemInfo } from '../api/types';

interface SystemStore {
  // State
  systemInfo: SystemInfo | null;
  hardwareInfo: SystemInfo | null;
  isLoading: boolean;
  error: string | null;
  lastUpdated: Date | null;
  
  // Actions
  loadSystemInfo: () => Promise<void>;
  loadHardwareInfo: () => Promise<void>;
  refreshSystemInfo: () => Promise<void>;
  clearError: () => void;
}

export const useSystemStore = create<SystemStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      systemInfo: null,
      hardwareInfo: null,
      isLoading: false,
      error: null,
      lastUpdated: null,

      // Actions
      loadSystemInfo: async () => {
        set({ isLoading: true, error: null });
        try {
          const systemInfo = await apiClient.getSystemInfo();
          set({ 
            systemInfo, 
            isLoading: false,
            lastUpdated: new Date(),
            error: null
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load system info',
            isLoading: false 
          });
        }
      },

      loadHardwareInfo: async () => {
        set({ isLoading: true, error: null });
        try {
          const hardwareInfo = await apiClient.getHardwareInfo();
          set({ 
            hardwareInfo, 
            isLoading: false,
            lastUpdated: new Date(),
            error: null
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to load hardware info',
            isLoading: false 
          });
        }
      },

      refreshSystemInfo: async () => {
        set({ isLoading: true, error: null });
        try {
          await apiClient.refreshSystemInfo();
          // Reload both system and hardware info
          const [systemInfo, hardwareInfo] = await Promise.all([
            apiClient.getSystemInfo(),
            apiClient.getHardwareInfo()
          ]);
          
          set({ 
            systemInfo,
            hardwareInfo,
            isLoading: false,
            lastUpdated: new Date(),
            error: null
          });
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Failed to refresh system info',
            isLoading: false 
          });
        }
      },

      clearError: () => set({ error: null })
    }),
    {
      name: 'system-store'
    }
  )
);