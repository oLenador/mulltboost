import React from 'react'

interface NoRequirementsI {
    title: string;
    text: string;
}

function NoRequirements({ title, text }: NoRequirementsI) {
    return (
        <section className='flex flex-col items-center w-full h-full pt-12 pb-32 overflow-scroll'>
            <div className='flex flex-col gap-8 w-full max-w-[1200px] items-center justify-center'>

                <div className='flex flex-col gap-4 shadow-lg bg-neutral-light-0/5 rounded-sm p-8 w-2/4'>
                    <h2 className=' text-xl'>{title}</h2>
                    <p className='text-base text-neutral-light-0/60'>{text}</p>
                    {
                        // Todo botão para enviaro usuário para o auto responder
                    }
                </div>
            </div>
        </section>
    )
}

export default NoRequirements