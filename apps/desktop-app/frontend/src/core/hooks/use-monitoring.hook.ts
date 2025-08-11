import { useEffect } from 'react';
import { useMonitoringStore } from '../store/monitoringStore';

export const useMonitoring = (autoStart: boolean = false) => {
  const store = useMonitoringStore();

  useEffect(() => {
    if (autoStart && !store.isMonitoring) {
      store.startMonitoring();
    }
    
    return () => {
      if (store.isMonitoring) {
        store.stopMonitoring();
      }
    };
  }, [autoStart]);

  return {
    currentMetrics: store.currentMetrics,
    metricsHistory: store.metricsHistory,
    isMonitoring: store.isMonitoring,
    isLoading: store.isLoading,
    error: store.error,
    updateInterval: store.updateInterval,
    
    // Actions
    startMonitoring: store.startMonitoring,
    stopMonitoring: store.stopMonitoring,
    getMetrics: store.getMetrics,
    clearHistory: store.clearHistory,
    setUpdateInterval: store.setUpdateInterval,
    clearError: store.clearError,
  };
};