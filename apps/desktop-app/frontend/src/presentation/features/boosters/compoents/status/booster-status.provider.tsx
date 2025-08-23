// BoosterStatusProvider.tsx - Versão Corrigida
import React, { ReactElement, useMemo, useCallback } from 'react';
import BoosterStatus from './booster-status.component';
import { PageType } from '@/presentation/pages/dashboard/dashboard';
import { useBoosterExecution } from '../../hooks/use-booster-execution.hook';

interface BoosterStatusProviderProps {
  children: ReactElement;
  path: PageType;
}

function BoosterStatusProvider({ children, path }: BoosterStatusProviderProps) {
  const {
    executions,
    executionsByStatus,
    stats,
    stagedOperations,
    isExecuting,
    executeStagedBatch
  } = useBoosterExecution();

  const progressMetrics = useMemo(() => {
    const totalItems = executions.length;
    const completed = executionsByStatus.completed.length;
    const currentProcessing = executionsByStatus.processing[0]; // Get first processing item
    const currentProgress = currentProcessing?.progress || 0;

    return {
      totalItems,
      completed,
      currentProgress,
      currentProcessing
    };
  }, [executions.length, executionsByStatus.completed.length, executionsByStatus.processing]);

  const circularLoaderProps = useMemo(() => {
    const { totalItems, completed, currentProgress } = progressMetrics;
    
    if (totalItems === 0) {
      return { 
        items: 4, 
        completed: 0, 
        currentProgress: 0 
      };
    }

    return {
      items: totalItems,
      completed: completed,
      currentProgress: currentProgress,
    };
  }, [progressMetrics]);

  const handleApply = useCallback(async () => {
    try {
      console.log('[BoosterStatusProvider] Executing staged batch...');
      await executeStagedBatch();
    } catch (error) {
      console.error('[BoosterStatusProvider] Failed to execute staged batch:', error);
    }
  }, [executeStagedBatch]);

  React.useEffect(() => {
    console.log('[BoosterStatusProvider] Re-rendered with:', {
      executionsCount: executions.length,
      isExecuting,
      stagedCount: Object.keys(stagedOperations).length
    });
  }, [executions.length, isExecuting, stagedOperations]);

  return (
    <>
      {children}
      <BoosterStatus
        path={path}
        executions={executions}
        stats={stats}
        stagedOperations={stagedOperations}
        isExecuting={isExecuting}
        handleApply={handleApply}
        circularLoaderProps={circularLoaderProps}
      />
    </>
  );
}

export default React.memo(BoosterStatusProvider);