import React from 'react'
import { Heading, Text } from '@radix-ui/themes';

interface HeadOfSectionI {
  sectionTitle: string;
  subTitle?: string;
  buttonPrimary: React.ReactNode;
  buttonSecundary?: React.ReactNode;
}

function HeadSection({ sectionTitle, subTitle, buttonPrimary, buttonSecundary }: HeadOfSectionI): React.ReactNode {
  return (
    <div className='flex flex-col gap-4'>
      <div className='flex flex-row w-full gap-16 justify-between items-center'>
        <div  className='flex flex-col gap-4'>
          <h6 className='font-semibold'>{sectionTitle}</h6>

          {subTitle &&
          <p className='text-text-p3 text-neutral-light-0/60 max-w-[560px]'>
              {subTitle}</p>
          }
        </div>
        <div className='flex flex-row gap-4 items-center justify-end'>
          {buttonSecundary}
          {buttonPrimary}
        </div>
      </div>


    </div>
  )
}

export default HeadSection