// src/presentation/features/boosters/components/execution-loader.component.tsx

import React, { useMemo, useEffect, useState } from "react"
import { ReactElement } from "react"
import { BoosterExecution, ExecutionStatus } from "../../domain/booster-queue.types"

interface ExecutionLoaderProps {
  executions: BoosterExecution[]
  size?: number
  strokeWidth?: number
  children?: ReactElement
  showPercentage?: boolean
  gap?: number
  showProgress?: boolean
  onAnimationComplete?: () => void
}

interface SegmentData {
  index: number
  startAngle: number
  endAngle: number
  segmentDegrees: number
  progress: number
  status: ExecutionStatus
  execution: BoosterExecution
}

const getSegmentColor = (status: ExecutionStatus): string => {
  switch (status) {
    case "completed":
      return "text-blue-500" // Azul para aplicada
    case "error":
      return "text-red-500" // Vermelho para erro
    case "cancelled":
      return "text-yellow-500" // Amarelo para cancelada
    case "processing":
    case "queued":
      return "text-blue-400" // Azul claro para em progresso
    case "idle":
    default:
      return "text-gray-400" // Cinza para não executada
  }
}

const getBackgroundColor = (status: ExecutionStatus): string => {
  switch (status) {
    case "error":
      return "text-red-200"
    case "cancelled":
      return "text-yellow-200"
    case "completed":
      return "text-blue-200"
    default:
      return "text-white/20"
  }
}

export function ExecutionLoader({
  executions,
  size = 64,
  strokeWidth = 4,
  children,
  showPercentage = false,
  gap = 4,
  showProgress = true,
  onAnimationComplete,
}: ExecutionLoaderProps) {
  const safeExecutions = executions || []
  const itemCount = Math.max(1, safeExecutions.length)

  const radius = (size - strokeWidth) / 2
  const center = size / 2

  const segments = useMemo((): SegmentData[] => {
    const currentGap = itemCount === 1 ? 1 : gap
    const totalGapDegrees = currentGap * itemCount
    const availableDegrees = 360 - totalGapDegrees
    const segmentDegrees = availableDegrees / itemCount

    return safeExecutions.map((execution, index) => {
      const startAngle = index * (segmentDegrees + currentGap)
      const endAngle = startAngle + segmentDegrees

      let progress = 0
      
      // Determinar progresso baseado no status
      switch (execution.status) {
        case "completed":
        case "error":
        case "cancelled":
          progress = 100
          break
        case "processing":
          progress = Math.max(10, execution.progress) // Mínimo 10% para mostrar que está processando
          break
        case "queued":
          progress = 5 // Pequeno progresso para mostrar que está na fila
          break
        case "idle":
        default:
          progress = 0
          break
      }

      return {
        index,
        startAngle,
        endAngle,
        segmentDegrees,
        progress,
        status: execution.status,
        execution
      }
    })
  }, [safeExecutions, itemCount, gap])

  const createArcPath = (startAngle: number, endAngle: number, progress: number = 100): string => {
    if (progress === 0) return ""

    const actualEndAngle = startAngle + ((endAngle - startAngle) * (progress / 100))
    const startAngleRad = (startAngle - 90) * (Math.PI / 180)
    const endAngleRad = (actualEndAngle - 90) * (Math.PI / 180)

    const startX = center + radius * Math.cos(startAngleRad)
    const startY = center + radius * Math.sin(startAngleRad)
    const endX = center + radius * Math.cos(endAngleRad)
    const endY = center + radius * Math.sin(endAngleRad)

    const largeArcFlag = actualEndAngle - startAngle > 180 ? 1 : 0

    return `M ${startX} ${startY} A ${radius} ${radius} 0 ${largeArcFlag} 1 ${endX} ${endY}`
  }

  const totalProgress = useMemo(() => {
    if (safeExecutions.length === 0) return "0.0"
    
    const totalProgress = safeExecutions.reduce((sum, exec) => {
      if (['completed', 'error', 'cancelled'].includes(exec.status)) {
        return sum + 100
      } else if (exec.status === 'processing') {
        return sum + exec.progress
      } else if (exec.status === 'queued') {
        return sum + 5
      }
      return sum
    }, 0)
    
    return (totalProgress / safeExecutions.length).toFixed(1)
  }, [safeExecutions])

  const containerStyle = {
    width: `${size}px`,
    height: `${size}px`,
  }

  // Estatísticas para debug/info
  const stats = useMemo(() => {
    return safeExecutions.reduce((acc, exec) => {
      acc[exec.status] = (acc[exec.status] || 0) + 1
      return acc
    }, {} as Record<ExecutionStatus, number>)
  }, [safeExecutions])

  return (
    <div 
      className="relative -top-[10px] -left-[10px] inline-flex items-center justify-center"
      style={containerStyle}
      data-execution-loader
      data-items={itemCount}
      data-progress={totalProgress}
      data-stats={JSON.stringify(stats)}
    >
      {showProgress && (
        <svg
          width={size}
          height={size}
          role="img"
          aria-label={`Execution Progress: ${totalProgress}% complete`}
          data-execution-loader-svg
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
              className={getBackgroundColor(segment.status)}
              data-segment-bg
              data-index={segment.index}
              data-status={segment.status}
              data-booster-id={segment.execution.boosterId}
            />
          ))}
          
          {/* Progress segments */}
          {segments.map((segment) => {
            const colorClass = getSegmentColor(segment.status)

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
                data-booster-id={segment.execution.boosterId}
                data-operation={segment.execution.operation}
              />
            )
          })}
        </svg>
      )}
      
      <div 
        className="absolute inset-0 flex items-center justify-center pointer-events-none"
        data-execution-loader-content
      >
        {children}
        {showPercentage && (
          <div className="text-xs text-white/80 absolute -bottom-1">
            {totalProgress}%
          </div>
        )}
      </div>
    </div>
  )
}

// Hook personalizado para usar com os átomos Jotai
export const useExecutionLoader = (executions: BoosterExecution[]) => {
  const [isVisible, setIsVisible] = useState(true)
  
  const handleAnimationComplete = () => {
    setIsVisible(false)
  }

  const resetVisibility = () => {
    setIsVisible(true)
  }

  return {
    isVisible,
    handleAnimationComplete,
    resetVisibility,
  }
}