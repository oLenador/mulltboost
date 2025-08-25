import React, { ReactElement, useMemo, useCallback } from 'react';
import { PageType } from '@/presentation/pages/dashboard/dashboard';
import { useBoosterExecution } from '../../hooks/use-booster-execution.hook';
import { BoosterStatus } from './booster-status.component';

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
    const currentProcessing = executionsByStatus.processing[0];
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
  
    // Mapear status dos itens baseado nas execuções
    const itemStatuses = executions.map(execution => {
      switch (execution.status) {
        case 'completed':
          return 'completed' as const;
        case 'processing':
          return 'loading' as const;
        case 'error':
          return 'error' as const;
        default:
          return 'pending' as const;
      }
    });
  
    return {
      items: totalItems,
      completed: completed,
      currentProgress: currentProgress,
      itemStatuses: itemStatuses 
    };
  }, [progressMetrics, executions]);

  const handleApply = useCallback(async () => {
    try {
      console.log('[BoosterStatusProvider-GLOBAL] Executing staged batch...');
      const result = await executeStagedBatch();
      
      if (result) {
        console.log(`[BoosterStatusProvider-GLOBAL] Batch executed successfully: ${result}`);
      } else {
        console.log('[BoosterStatusProvider-GLOBAL] No staged operations to execute');
      }
    } catch (error) {
      console.error('[BoosterStatusProvider-GLOBAL] Failed to execute staged batch:', error);
    }
  }, [executeStagedBatch]);

  React.useEffect(() => {
    console.log('[BoosterStatusProvider-GLOBAL] Re-rendered with:', {
      executionsCount: executions.length,
      isExecuting,
      stagedCount: Object.keys(stagedOperations).length,
      path
    } ,    executions,
    executionsByStatus,
    stats,
    stagedOperations,
    isExecuting,);
  }, [executions.length, isExecuting, stagedOperations, path]);

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