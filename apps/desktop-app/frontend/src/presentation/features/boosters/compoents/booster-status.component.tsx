import { FloatElement } from '@/presentation/components/floating-manager'
import React from 'react'
import { LoadingItem, useBoosterStatus } from '../hooks/booster-status.hook'
import { CircularLoader } from './circular-loader.component'
import { BoosterItem } from '../types/booster.types'

interface BoosterStatusProps {
  path: string
  boosterQueue: BoosterItem[]
  completed: number
  isLoading?: boolean
  progress?: number
  onToggleVisibility?: () => void
}

function BoosterStatus({ 
  path, 
  items, 
  completed, 
  isLoading = false,
  progress,
  onToggleVisibility 
}: BoosterStatusProps) {

  const boosterStatus = useBoosterStatus()

  // Condições de visibilidade mais inteligentes
  const hasItems = items > 0
  const isVisible = hasItems && (isLoading || completed < items)
  
  if (!isVisible) {
    return null
  }

  return (
    <FloatElement
      id={`booster-status-${path}`}
      type="custom"
      position="bottom-right"
      priority={5}
      active={isVisible}
    >
      <div 
        className="rounded-full bg-white shadow-lg border border-gray-200 hover:shadow-xl transition-shadow duration-200 cursor-pointer"
        onClick={onToggleVisibility}
      >
        <CircularLoader 
          items={items} 
          completed={completed}
          progress={progress}
          size={64}
          strokeWidth={3}
        >
          <div className="flex flex-col items-center justify-center">
            <span className="text-xs font-semibold text-blue-700">
              {completed}/{items}
            </span>
            {isLoading && (
              <div className="w-1 h-1 bg-blue-500 rounded-full animate-pulse mt-0.5" />
            )}
          </div>
        </CircularLoader>
      </div>
    </FloatElement>
  )
}


export { BoosterStatus }
export default BoosterStatus