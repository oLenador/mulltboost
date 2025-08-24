// src/presentation/features/boosters/components/circular-loader.component.tsx

import React, { useMemo } from "react"
import { ReactElement } from "react"

interface CircularLoaderProps {
  items: number
  completed: number
  currentProgress?: number
  size?: number
  strokeWidth?: number
  children?: ReactElement
  showPercentage?: boolean
  gap?: number
  showProgress: boolean
  itemStatuses?: ("pending" | "loading" | "completed" | "error")[] // Nova prop
}

interface SegmentData {
  index: number;
  startAngle: number;
  endAngle: number;
  segmentDegrees: number;
  progress: number;
  status: "pending" | "loading" | "completed" | "error";
}

export function CircularLoader({
  items,
  completed,
  currentProgress = 0,
  size = 64,
  strokeWidth = 4,
  children,
  showPercentage = false,
  gap = 4,
  showProgress,
  itemStatuses = [],
}: CircularLoaderProps) {
  const safeItems = Math.max(1, items || 1);
  const safeCompleted = Math.max(0, Math.min(completed || 0, safeItems));
  const safeCurrentProgress = Math.max(0, Math.min(currentProgress, 100));

  const radius = (size - strokeWidth) / 2;
  const center = size / 2;
const segments = useMemo((): SegmentData[] => {
  const currentGap = safeCompleted == 1 ? 1 : gap
  const totalGapDegrees = currentGap * safeItems;
  const availableDegrees = 360 - totalGapDegrees;
  const segmentDegrees = availableDegrees / safeItems;

  return Array.from({ length: safeItems }, (_, index) => {
    const startAngle = index * (segmentDegrees + currentGap);
    const endAngle = startAngle + segmentDegrees;

    let progress = 0;
    let status: "pending" | "loading" | "completed" | "error" = "pending";

    // Verificar se há status específico fornecido
    if (itemStatuses && itemStatuses[index]) {
      status = itemStatuses[index];
      if (status === "error") {
        progress = 100; // Mostra o segmento completo em vermelho
      } else if (status === "completed") {
        progress = 100;
      } else if (status === "loading" && index === safeCompleted) {
        progress = safeCurrentProgress;
      }
    } else {
      // Lógica original
      if (index < safeCompleted) {
        progress = 100;
        status = "completed";
      } else if (index === safeCompleted && safeCurrentProgress > 0) {
        progress = safeCurrentProgress;
        status = "loading";
      }
    }

    return {
      index,
      startAngle,
      endAngle,
      segmentDegrees,
      progress,
      status
    };
  });
}, [safeItems, safeCompleted, safeCurrentProgress, gap, itemStatuses]);

  const createArcPath = (startAngle: number, endAngle: number, progress: number = 100): string => {
    if (progress === 0) return "";

    const actualEndAngle = startAngle + ((endAngle - startAngle) * (progress / 100));
    const startAngleRad = (startAngle - 90) * (Math.PI / 180);
    const endAngleRad = (actualEndAngle - 90) * (Math.PI / 180);

    const startX = center + radius * Math.cos(startAngleRad);
    const startY = center + radius * Math.sin(startAngleRad);
    const endX = center + radius * Math.cos(endAngleRad);
    const endY = center + radius * Math.sin(endAngleRad);

    const largeArcFlag = actualEndAngle - startAngle > 180 ? 1 : 0;

    return `M ${startX} ${startY} A ${radius} ${radius} 0 ${largeArcFlag} 1 ${endX} ${endY}`;
  };

  const totalProgress = useMemo(() => {
    if (safeItems === 0) return "0.0";
    
    const completedItems = safeCompleted * 100;
    const currentItemProgress = safeCurrentProgress;
    return ((completedItems + currentItemProgress) / safeItems).toFixed(1);
  }, [safeCompleted, safeCurrentProgress, safeItems]);

  const containerStyle = {
    width: `${size}px`,
    height: `${size}px`,
  };

  const getSegmentColor = (status: SegmentData['status']): string => {
    switch (status) {
      case "completed": 
        return "text-green-500";
      case "loading": 
        return "text-blue-500";
      case "error":
        return "text-red-500";
      default: 
        return "text-gray-300";
    }
  };

  return (
    <div 
      className="relative -top-[10px] -left-[10px] inline-flex items-center justify-center"
      style={containerStyle}
      data-circular-loader
      data-items={safeItems}
      data-completed={safeCompleted}
      data-progress={totalProgress}
    >
      {showProgress && (
        <svg
          width={size}
          height={size}
          role="img"
          aria-label={`Progress: ${totalProgress}% complete`}
          data-circular-loader-svg
        >
          {/* Background segments */}
          {segments.map((segment) => (
            <path
              key={`bg-${segment.index}`}
              d={createArcPath(segment.startAngle, segment.endAngle, 100)}
              stroke="currentColor"
              strokeWidth={strokeWidth}
              fill="none"
              strokeLinecap="butt"
              className={segment.status === 'error' ? 'text-red-200' : 'text-white/20'}
              data-segment-bg
              data-index={segment.index}
              data-status={segment.status}
            />
          ))}
          
          {/* Progress segments */}
          {segments.map((segment) => {
            const colorClass = getSegmentColor(segment.status);

            return (
              <path
                key={`progress-${segment.index}`}
                d={createArcPath(segment.startAngle, segment.endAngle, segment.progress)}
                stroke="currentColor"
                strokeWidth={strokeWidth}
                fill="none"
                strokeLinecap="round"
                className={`transition-all duration-500 ease-out ${colorClass}`}
                style={{
                  transitionTimingFunction: "cubic-bezier(0.4, 0.0, 0.2, 1)",
                }}
                data-segment
                data-index={segment.index}
                data-status={segment.status}
                data-progress={segment.progress}
              />
            );
          })}
        </svg>
      )}
      
      <div 
        className="absolute inset-0 flex items-center justify-center pointer-events-none"
        data-circular-loader-content
      >
        {children}
        {showPercentage && (
          <div className="text-xs text-white/80 absolute -bottom-1">
            {totalProgress}%
          </div>
        )}
      </div>
    </div>
  );
}
