import { FloatElement } from '@/presentation/components/floating-manager'
import React, { ReactElement, useEffect, useState } from 'react'
// import { LoadingItem, useBoosterStatus } from '../hooks/booster-status.hook'
import { CircularLoader } from './circular-loader.component'
import { BoosterItem } from '../types/booster.types'
import { Check, Play } from 'lucide-react'
import { PageType } from '@/presentation/pages/dashboard/dashboard'
import { MdPlayArrow} from "react-icons/md"
interface BoosterStatusProps {
  path: PageType
  boosterQueue: BoosterItem[]
  completed: number
  boostersSelected: string[]
  isLoading?: boolean
  progress?: number
  items: any
  handleApply?: () => void
}
type StatusStates = "apply" | 'showing_progress' | 'completed_animation'

function ShowingProgress({ completed, totalItems }: { completed: number, totalItems: number }) {
  return (<span className='text-white -mb-1'>{completed}/{totalItems}</span>)
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
  items,
  completed,
  isLoading = false,
  boostersSelected,
  progress,
  handleApply
}: BoosterStatusProps) {

  const possibleStates: Record<StatusStates, React.ReactNode> = {
    apply: <MdPlayArrow size={24} fill={"#fff"} className='text-white' />,
    showing_progress: <ShowingProgress completed={completed} totalItems={items.length} />,
    completed_animation: <Check fill={"#fff"} className='text-white' />,
  }
  const hasItems = items.length > 0
  const [loaderChild, setLoaderChild] = useState<React.ReactNode | undefined>(null)
  const isVisible = !!loaderChild


  useEffect(() => {
    
    if (boostersSelected.length > 0 && isBoosterPage(path)) {
      setLoaderChild(possibleStates["apply"])
      return
    }

    if (hasItems) {
      setLoaderChild(possibleStates["showing_progress"])
      return
    }
    if (!isBoosterPage(path)) {
      setLoaderChild(null)
    } 

  }, [hasItems, boostersSelected, path])
  if (!isVisible) {
    return
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
        className="rounded-full w-16 h-16 hover:shadow-white/[0.05] bg-blue-600 hover:bg-blue-500/80 shadow-lg border border-white/20 hover:shadow-xl transition-shadow duration-200 cursor-pointer"
        onClick={boostersSelected.length > 0 ? handleApply : () => {}}
      >
        <CircularLoader
          items={4}
          completed={0}
          currentProgress={42}
          size={82}
          gap={14}
          strokeWidth={5}
          showProgress={hasItems}
        >
          <div className=" flex flex-col items-center justify-center">
            {loaderChild}
          </div>
        </CircularLoader>
      </div>
    </FloatElement>
  )
}


export { BoosterStatus }
export default BoosterStatus