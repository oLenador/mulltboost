import { useEffect } from 'react';
import { useSystemStore } from '../store/systemStore';

export const useSystemInfo = (loadOnMount: boolean = true) => {
  const store = useSystemStore();

  useEffect(() => {
    if (loadOnMount && !store.systemInfo) {
      Promise.all([
        store.loadSystemInfo(),
        store.loadHardwareInfo()
      ]);
    }
  }, [loadOnMount]);

  return {
    systemInfo: store.systemInfo,
    hardwareInfo: store.hardwareInfo,
    isLoading: store.isLoading,
    error: store.error,
    lastUpdated: store.lastUpdated,
    
    // Actions
    loadSystemInfo: store.loadSystemInfo,
    loadHardwareInfo: store.loadHardwareInfo,
    refreshSystemInfo: store.refreshSystemInfo,
    clearError: store.clearError,
  };
};