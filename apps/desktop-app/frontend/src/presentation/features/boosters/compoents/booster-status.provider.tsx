import React, { ReactElement } from 'react'
import BoosterStatus from './booster-status.component'
import { PageType } from '@/presentation/pages/dashboard/dashboard'
import { useBatchManager } from '../hooks/use-batch-manager.hook'

function BoosterStatusProvider({ children, path }: { path: PageType; children: ReactElement }) {
  const { startStagedBatch, stagedItems, items, isProcessing } = useBatchManager()

  // completed: quantos items com status === 'completed'
  console.log(items, stagedItems)
  const completed = items.filter((it) => it.status === 'completed').length

  return (
    <>
      {children}
      <BoosterStatus
        path={path}
        items={items}
        handleApply={startStagedBatch}
        boosterQueue={[]} 
        completed={completed}
        boostersSelected={Object.keys(stagedItems)}
        isLoading={isProcessing}
      />
    </>
  )
}

export default BoosterStatusProvider
