import React from "react"
import { ReactElement } from "react"

interface CircularLoaderProps {
    items: number
    completed: number
    progress?: number
    size?: number
    strokeWidth?: number
    children: ReactElement
    showPercentage?: boolean
  }
  
 export function CircularLoader({
    items,
    completed,
    progress,
    size = 64,
    strokeWidth = 4,
    children,
    showPercentage = false
  }: CircularLoaderProps) {
    const calculatedProgress = progress !== undefined ? progress : (items > 0 ? (completed / items) * 100 : 0)
    const radius = (size - strokeWidth) / 2
    const circumference = radius * 2 * Math.PI
    const strokeDasharray = circumference
    const strokeDashoffset = circumference - (calculatedProgress / 100) * circumference
  
    return (
      <div className="relative inline-flex items-center justify-center">
        <svg width={size} height={size} className="transform -rotate-90">
          <circle
            cx={size / 2}
            cy={size / 2}
            r={radius}
            stroke="rgb(229, 231, 235)"
            strokeWidth={strokeWidth}
            fill="none"
          />
          <circle
            cx={size / 2}
            cy={size / 2}
            r={radius}
            stroke="rgb(59, 130, 246)"
            strokeWidth={strokeWidth}
            fill="none"
            strokeDasharray={strokeDasharray}
            strokeDashoffset={strokeDashoffset}
            strokeLinecap="round"
            className="transition-all duration-300 ease-in-out"
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          {showPercentage ? (
            <span className="text-sm font-medium text-gray-700">
              {Math.round(calculatedProgress)}%
            </span>
          ) : (
            children
          )}
        </div>
      </div>
    )
  }