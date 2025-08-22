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
}: CircularLoaderProps) {
  const safeItems = Math.max(1, items || 1)
  const safeCompleted = Math.max(0, Math.min(completed || 0, safeItems))
  const safeCurrentProgress = Math.max(0, Math.min(currentProgress, 100))

  const radius = (size - strokeWidth) / 2
  const center = size / 2

  const segments = useMemo(() => {
    const totalGapDegrees = gap * safeItems
    const availableDegrees = 360 - totalGapDegrees
    const segmentDegrees = availableDegrees / safeItems

    return Array.from({ length: safeItems }, (_, index) => {
      const startAngle = index * (segmentDegrees + gap)
      const endAngle = startAngle + segmentDegrees

      let progress = 0
      let status: "pending" | "loading" | "completed" = "pending"

      if (index < safeCompleted) {
        progress = 100
        status = "completed"
      } else if (index === safeCompleted) {
        progress = safeCurrentProgress
        status = progress > 0 ? "loading" : "pending"
      }

      return {
        index,
        startAngle,
        endAngle,
        segmentDegrees,
        progress,
        status
      }
    })
  }, [safeItems, safeCompleted, safeCurrentProgress, gap])

  const createArcPath = (startAngle: number, endAngle: number, progress: number = 100) => {
    const actualEndAngle = startAngle + ((endAngle - startAngle) * (progress / 100))
    const startAngleRad = (startAngle - 90) * (Math.PI / 180)
    const endAngleRad = (actualEndAngle - 90) * (Math.PI / 180)

    const startX = center + radius * Math.cos(startAngleRad)
    const startY = center + radius * Math.sin(startAngleRad)
    const endX = center + radius * Math.cos(endAngleRad)
    const endY = center + radius * Math.sin(endAngleRad)

    const largeArcFlag = actualEndAngle - startAngle > 180 ? 1 : 0

    if (progress === 0) return ""

    return `M ${startX} ${startY} A ${radius} ${radius} 0 ${largeArcFlag} 1 ${endX} ${endY}`
  }

  const totalProgress = useMemo(() => {
    const completedItems = safeCompleted * 100
    const currentItemProgress = safeCurrentProgress
    return ((completedItems + currentItemProgress) / safeItems).toFixed(1)
  }, [safeCompleted, safeCurrentProgress, safeItems])

  // SOLUÇÃO: Sempre definir o tamanho da div container
  const containerStyle = {
    width: `${size}px`,
    height: `${size}px`,
  }

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
          data-circular-loader-svg
        >
          {segments.map((segment) => (
            <path
              key={`bg-${segment.index}`}
              d={createArcPath(segment.startAngle, segment.endAngle, 100)}
              stroke="currentColor"
              strokeWidth={strokeWidth}
              fill="none"
              strokeLinecap="round"
              rx="0.1"
              ry="0.1"
              className="text-white/20"
              data-segment-bg
              data-index={segment.index}
            />
          ))}
          {segments.map((segment) => {
            let colorClass = "text-gray-300"
            if (segment.status === "completed") colorClass = "text-green-500"
            else if (segment.status === "loading") colorClass = "text-blue-500"

            return (
              <path
                key={`progress-${segment.index}`}
                d={createArcPath(segment.startAngle, segment.endAngle, segment.progress)}
                stroke="currentColor"
                strokeWidth={strokeWidth}
                fill="none"
                strokeLinecap="square"
                rx="1"
                ry="1"
                className={`transition-all duration-500 ease-out ${colorClass}`}
                style={{
                  transitionTimingFunction: "cubic-bezier(0.4, 0.0, 0.2, 1)",
                }}
                data-segment
                data-index={segment.index}
                data-status={segment.status}
                data-progress={segment.progress}
              />
            )
          })}
        </svg>
      )}
      
      <div 
        className="absolute inset-0 flex items-center justify-center pointer-events-none"
        data-circular-loader-content
      >
        {children}
      </div>
    </div>
  )
}