import { FloatElement } from '@/presentation/components/floating-manager'
import React, { ReactElement, useState } from 'react'
// import { LoadingItem, useBoosterStatus } from '../hooks/booster-status.hook'
import { CircularLoader } from './circular-loader.component'
import { BoosterItem } from '../types/booster.types'
import { Check, Play } from 'lucide-react'
import { PageType } from '@/presentation/pages/dashboard/dashboard'

interface BoosterStatusProps {
  path: PageType
  boosterQueue: BoosterItem[]
  completed: number
  boostersSelected: boolean
  isLoading?: boolean
  progress?: number
  items: any
  onToggleVisibility?: () => void
}
type StatusStates = "apply" | 'showing_progress' | 'completed_animation'

function ShowingProgress({completed, totalItems}: {completed: number, totalItems: number}) {
 return (<span className='text-white'>{completed}/{totalItems}</span>)
}

const isBoosterPage = (path: PageType) => {
  const boosterPageSet = new Set<PageType>([
    PageType.FPS_BOOST,
    PageType.CONNECTION,
    PageType.PRECISION,
    PageType.GAMES,
    PageType.FLUSHER,
  ])
  return boosterPageSet.has(path)
}

function BoosterStatus({
  path,
  items = [],
  completed,
  isLoading = false,
  boostersSelected,
  progress,
  onToggleVisibility
}: BoosterStatusProps) {

  // const boosterStatus = useBoosterStatus()
  const possibleStates: Record<StatusStates, React.ReactNode> = {
    apply: <Play fill={"#fff"} className='text-white' />,
    showing_progress: <ShowingProgress completed={completed} totalItems={items.length} />,
    completed_animation: <Check fill={"#fff"} className='text-white' />,
  }
  const hasItems = true // items > 0
   // hasItems && (isLoading || completed < items)
  const [loaderChild, setLoaderChild] = useState<React.ReactNode | undefined>(possibleStates["apply"])
  const isVisible = !!loaderChild
  // if (!isVisible) {
  //   return null
  // }

  useState(() => {
    if (boostersSelected && isBoosterPage(path)) {
      setLoaderChild(possibleStates["apply"])
      return 
    }

    if (hasItems) {
      setLoaderChild(possibleStates["showing_progress"])
      return
    }
  })
  

  return (
    <FloatElement
      id={`booster-status-${path}`}
      type="custom"
      position="bottom-right"

      priority={5}
      active={isVisible}
    >
      <div
        className="rounded-full w-16 h-16 hover:shadow-white/[0.05] bg-blue-600 hover:bg-blue-500/80 shadow-lg border border-gray-200/10 hover:shadow-xl transition-shadow duration-200 cursor-pointer"
        onClick={onToggleVisibility}
      >
        <CircularLoader
          items={4}
          completed={0}
          currentProgress={42}
          size={82}
gap={14}
          strokeWidth={5}
        >
          <div className="flex flex-col items-center justify-center">
          {loaderChild}
          </div>
        </CircularLoader>
      </div>
    </FloatElement>
  )
}


export { BoosterStatus }
export default BoosterStatus