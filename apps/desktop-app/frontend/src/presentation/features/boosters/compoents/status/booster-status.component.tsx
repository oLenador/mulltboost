// src/presentation/features/boosters/components/booster-status.component.tsx

import { FloatElement } from '@/presentation/components/floating-manager';
import React, { useEffect, useState } from 'react';
import { CircularLoader } from './circular-loader.component';
import { Check, Play } from 'lucide-react';
import { PageType } from '@/presentation/pages/dashboard/dashboard';
import { MdPlayArrow } from "react-icons/md";
import { BoosterExecution, ExecutionStats } from '../../domain/booster-queue.types';
import { StagedOperations } from '@/presentation/features/boosters/stores/booster-execution.store';

interface BoosterStatusProps {
  path: PageType;
  executions: BoosterExecution[];
  stats: ExecutionStats;
  stagedOperations: StagedOperations;
  isExecuting: boolean;
  handleApply: () => void;
  circularLoaderProps: {
    items: number;
    completed: number;
    currentProgress: number;
  };
}

type StatusState = "apply" | 'showing_progress' | 'completed_animation' | 'idle';

interface ShowingProgressProps {
  completed: number;
  total: number;
  processing: number;
}

function ShowingProgress({ completed, total, processing }: ShowingProgressProps) {
  if (processing > 0) {
    return (
      <div className="flex flex-col items-center text-white">
        <span className="text-xs -mb-1">Processing</span>
        <span className="text-sm font-semibold">{completed}/{total}</span>
      </div>
    );
  }
  
  return (
    <span className='text-white text-sm font-semibold'>
      {completed}/{total}
    </span>
  );
}

const isBoosterPage = (path: PageType): boolean => {
  const boosterPageSet = new Set<PageType>([
    PageType.FPS_BOOST,
    PageType.CONNECTION,
    PageType.PRECISION,
    PageType.GAMES,
    PageType.FLUSHER,
  ]);
  return boosterPageSet.has(path);
};

function BoosterStatus({
  path,
  executions,
  stats,
  stagedOperations,
  isExecuting,
  handleApply,
  circularLoaderProps
}: BoosterStatusProps) {
  const [currentState, setCurrentState] = useState<StatusState>('idle');
  const [showCompletedAnimation, setShowCompletedAnimation] = useState(false);

  const stagedCount = Object.keys(stagedOperations).length;
  const hasExecutions = executions.length > 0;
  const hasStaged = stagedCount > 0;

  // Determine current state
  useEffect(() => {
    if (!isBoosterPage(path)) {
      setCurrentState('idle');
      return;
    }

    if (showCompletedAnimation) {
      setCurrentState('completed_animation');
      return;
    }

    if (hasStaged && !isExecuting) {
      setCurrentState('apply');
    } else if (isExecuting || hasExecutions) {
      setCurrentState('showing_progress');
    } else {
      setCurrentState('idle');
    }
  }, [path, hasStaged, isExecuting, hasExecutions, showCompletedAnimation]);

  // Handle completion animation
  useEffect(() => {
    if (stats.completed > 0 && stats.processing === 0 && stats.queued === 0 && hasExecutions) {
      setShowCompletedAnimation(true);
      const timer = setTimeout(() => {
        setShowCompletedAnimation(false);
      }, 2000); // Show animation for 2 seconds
      
      return () => clearTimeout(timer);
    }
  }, [stats.completed, stats.processing, stats.queued, hasExecutions]);


  
  const renderContent = (): React.ReactNode => {
    switch (currentState) {
      case 'apply':
        return (
          <div className="flex items-center justify-center">
            <MdPlayArrow size={24} fill="#fff" className="text-white" />
          </div>
        );
      
      case 'showing_progress':
        return (
          <ShowingProgress 
            completed={stats.completed} 
            total={stats.total}
            processing={stats.processing}
          />
        );
      
      case 'completed_animation':
        return (
          <div className="flex items-center justify-center animate-pulse">
            <Check size={24} className="text-green-400" />
          </div>
        );
      
      default:
        return null;
    }
  };

  const shouldShow = currentState !== 'idle';
  const canClick = currentState === 'apply';

  if (!shouldShow) {
    return null;
  }

  return (
    <FloatElement
      id={`booster-status-${path}`}
      type="custom"
      position="bottom-right"
      priority={5}
      active={shouldShow}
    >
      <div
        className={`
          rounded-full w-16 h-16 shadow-lg border border-white/20 
          transition-all duration-200 
          ${canClick 
            ? 'hover:shadow-white/[0.05] bg-blue-700 hover:bg-blue-500/80 hover:shadow-xl cursor-pointer' 
            : 'bg-blue-600/80'
          }
          ${currentState === 'completed_animation' ? 'bg-green-600' : ''}
        `}
        onClick={canClick ? handleApply : undefined}
      >
        <CircularLoader
          items={circularLoaderProps.items}
          completed={circularLoaderProps.completed}
          currentProgress={circularLoaderProps.currentProgress}
          size={82}
          gap={14}
          strokeWidth={5}
          showProgress={hasExecutions || isExecuting}
        >
          <div className="flex flex-col items-center justify-center">
            {renderContent()}
          </div>
        </CircularLoader>
      </div>
    </FloatElement>
  );
}

export { BoosterStatus };
export default BoosterStatus;