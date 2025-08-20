import { useState, useEffect, useCallback } from 'react';
import { Events } from '@wailsio/runtime';
import { GetSystemMetrics } from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/metricshandler';
import { SystemMetrics } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

export const useMonitoring = (autoStart: boolean = false) => {
  const [currentMetrics, setCurrentMetrics] = useState<any | null>(null);
  const [metricsHistory, setMetricsHistory] = useState<SystemMetrics[]>([]);
  const [isMonitoring, setIsMonitoring] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [updateInterval, setUpdateInterval] = useState(1000);
  const [maxHistorySize] = useState(100);

  const clearHistory = useCallback(() => {
    setMetricsHistory([]);
  }, []);
  
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  useEffect(() => {
    let cancelEvent: (() => void) | null = null;

    const fetchInitialMetrics = async () => {
      setIsLoading(true);
      try {
        const initialMetrics = await GetSystemMetrics()
        setCurrentMetrics(initialMetrics);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setIsLoading(false);
      }
    };

    // Busca o estado atual na inicialização
    fetchInitialMetrics();

    // Escuta eventos em tempo real
    cancelEvent = Events.On('metrics-change', (metrics: any) => {
      setCurrentMetrics(metrics);
    });

    return () => {
      if (cancelEvent) cancelEvent();
    };
  }, [autoStart, maxHistorySize]);

  return {
    currentMetrics,
    metricsHistory,
    isMonitoring,
    isLoading,
    error,
    updateInterval,
    clearHistory,
    setUpdateInterval,
    clearError,
  };
};
