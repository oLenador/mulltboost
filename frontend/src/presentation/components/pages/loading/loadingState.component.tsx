import React from 'react'
import "./loadingState.component.scss"

function LoadingState() {
    return (
        <div className='w-full h-[100vh] bg-black text-white flex flex-col items-center justify-center overflow-hidden'>
                <div id="loadingMainDot"></div>
        </div>
    )
}

export default LoadingState