import React from 'react'
import SARSIcon from '../../components/SARSIcon'

function expired() {
    return (
        <div className='w-full h-[100vh] bg-black text-white flex flex-col items-center justify-center overflow-hidden'>
            <div className='flex flex-col gap-8 items-center justify-center w-[560px]'>
                <SARSIcon width={96} />
                <div className='w-full flex flex-col gap-16 bg-neutral-dark-0 rounded-lg py-16 px-16'>

                </div>
            </div>
        </div>
    )
}

export default expired