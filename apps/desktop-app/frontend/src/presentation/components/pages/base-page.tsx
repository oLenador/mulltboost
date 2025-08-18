import React, { ReactElement } from 'react'

function BasePage({ children }: { children: ReactElement}) {
  return (
    <div className="w-full px-8 py-4 space-y-8 bg-zinc-950 min-h-screen text-zinc-100">
        {children}
      </div>
  )
}

export default BasePage