
import { PrimaryLightCTABoxes } from '@/presentation/components/Buttons/Primary/PrimaryLightCTABoxes'
import { PrimaryButton } from '@/presentation/components/Buttons/Primary/PrimaryLightWeight'
import { OneRealCoin } from '@/presentation/components/coin/one-real-coin'
import { ArrowUp } from 'lucide-react'
import React, { Suspense } from 'react'


function BoxesPage() {

  let participants = 23
  return (
    <section className=" flex flex-col justify-between h-full w-full max-w-[556px] px-16  pt-[40px]   bg-background-200 border-r border-neutral-light-0/5">
      <div className='flex flex-col items-center w-full'>
        <div className='w-44 h-28 bg-neutral-dark-400 rounded-md'>
          <img />
        </div>

        <div className='flex flex-col items-center mt-2'>
          <span className="text-text-p4 text-[#778EFF] uppercase font-semibold">
            Heeph Boxes
          </span>


          <span className='text-white text-center text-text-pl'>Você Não precisa pagar <br /><span className='text-white'>R$ 127,98</span> no Mine!</span>
        </div>
        <div className='flex flex-col mt-6 w-full'>
          <div className='flex flex-col gap-1'>
            <span className=' text-white/60 text-sm'>
              Concorra por 1 REAL
            </span>
            <div className='flex flex-col gap-2 rounded-md bg-black/60 border border-white/10  p-4'>
              <span>- Minecraft Original (FULL ACESSO)</span>
              <span>- Mouse Bloody A91</span>
              <span>- Teclado Gamer Logitech G213</span>
            </div>
          </div>

          <div className='flex flex-col gap-1 mt-6'>
            <span className='uppercase text-white/60 text-sm'>
              Prêmios GARANTIDOS
            </span>
            <div className='rounded-md bg-black/60 border border-white/10  p-4'>
              <span>- Capa Optifine (Ganha Sempre)</span>
            </div>
          </div >
          <div className='w-full flex items-center justify-center text-green-300'>
            <span className='w-fit items-center mt-2 flex flex-row gap-1 -ml-3'>Você <span className='font-medium'>SEMPRE</span> ganha!</span>
          </div>
          <hr className='border border-white/10' />

            <div className='w-full mt-6'>
              <PrimaryLightCTABoxes title='Participar' onClickFn={() => { }} />
            </div>
        </div>
      </div>
    </section>
  )
}

export default BoxesPage



