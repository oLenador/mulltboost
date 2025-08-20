import { atom } from 'jotai';
import type { SystemMetrics } from '../api/types';

// ---------- State atoms ----------
export const currentMetricsAtom = atom<SystemMetrics | null>(null);
export const metricsHistoryAtom = atom<SystemMetrics[]>([]);
export const isMonitoringAtom = atom(false);
export const isLoadingAtom = atom(false);
export const errorAtom = atom<string | null>(null);

export const updateIntervalAtom = atom(1000); // 1s
export const maxHistorySizeAtom = atom(100);

const pollingIntervalAtom = atom<NodeJS.Timeout | null>(null);

// ---------- Actions ----------
export const startMonitoringAtom = atom(
  null,
  async (get, set, interval?: number) => {
    const currentInterval = get(updateIntervalAtom);
    const intervalToUse = interval ?? currentInterval;

    set(isLoadingAtom, true);
    set(errorAtom, null);

    try {
      await apiClient.startRealTimeMonitoring(Math.floor(intervalToUse / 1000));

      set(isMonitoringAtom, true);
      set(isLoadingAtom, false);
      set(updateIntervalAtom, intervalToUse);

      // Start local polling
      const existingInterval = get(pollingIntervalAtom);
      if (existingInterval) {
        clearInterval(existingInterval);
      }
      const newInterval = setInterval(() => {
        set(getMetricsAtom); // dispara a action abaixo
      }, get(updateIntervalAtom));

      set(pollingIntervalAtom, newInterval);
    } catch (err) {
      set(errorAtom, err instanceof Error ? err.message : 'Failed to start monitoring');
      set(isLoadingAtom, false);
    }
  }
);

export const stopMonitoringAtom = atom(
  null,
  async (get, set) => {
    set(isLoadingAtom, true);
    try {
      await apiClient.stopRealTimeMonitoring();
      set(isMonitoringAtom, false);
      set(isLoadingAtom, false);

      const intervalId = get(pollingIntervalAtom);
      if (intervalId) {
        clearInterval(intervalId);
        set(pollingIntervalAtom, null);
      }
    } catch (err) {
      set(errorAtom, err instanceof Error ? err.message : 'Failed to stop monitoring');
      set(isLoadingAtom, false);
    }
  }
);

export const getMetricsAtom = atom(
  null,
  async (get, set) => {
    try {
      const metrics = await apiClient.getSystemMetrics();
      const history = get(metricsHistoryAtom);
      const maxSize = get(maxHistorySizeAtom);

      const newHistory = [...history, metrics].slice(-maxSize);

      set(currentMetricsAtom, metrics);
      set(metricsHistoryAtom, newHistory);
      set(errorAtom, null);
    } catch (err) {
      set(errorAtom, err instanceof Error ? err.message : 'Failed to get metrics');
    }
  }
);

export const clearHistoryAtom = atom(null, (_get, set) => {
  set(metricsHistoryAtom, []);
});

export const setUpdateIntervalAtom = atom(
  null,
  (_get, set, interval: number) => {
    set(updateIntervalAtom, interval);
  }
);

export const clearErrorAtom = atom(null, (_get, set) => {
  set(errorAtom, null);
});
