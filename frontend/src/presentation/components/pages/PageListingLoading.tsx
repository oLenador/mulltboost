import React, { memo } from 'react'

function PageListingLoading() {
    return (
        <section className='flex flex-col items-center w-full h-full py-16  overflow-scroll'>

            <div className='flex flex-col gap-8 w-full max-w-[1200px] h-full'>
                <div className='flex flex-row justify-between items-center w-full'>
                    <div className='animate-pulse min-h-12 h-12 w-56 bg-neutral-light-0/[0.07] rounded-lg'></div>
                    <div className='animate-pulse min-h-12 h-12 w-32 bg-neutral-light-0/[0.07] rounded-lg'></div>
                </div>

                <hr className='border-neutral-light-0/10 animate-pulse' />

                <div className='flex flex-col gap-4'>
                    <div className='animate-pulse min-h-12 h-12 w-full bg-neutral-light-0/[0.07] rounded-lg'></div>
                    <div className='animate-pulse min-h-12 h-12 w-full bg-neutral-light-0/[0.07] rounded-lg'></div>
                    <div className='animate-pulse min-h-12 h-12 w-full bg-neutral-light-0/[0.07] rounded-lg'></div>
                    <div className='animate-pulse min-h-12 h-12 w-full bg-neutral-light-0/[0.07] rounded-lg'></div>
                    <div className='animate-pulse min-h-12 h-12 w-full bg-neutral-light-0/[0.07] rounded-lg'></div>
                </div>

            </div>
        </section>
    )
}

export default memo(PageListingLoading)