

import React, { ReactElement } from 'react'
 import BoosterStatus from './booster-status.component'
import { LoadingItem } from '../hooks/booster-status.hook';
import { PageType } from '@/presentation/pages/dashboard/dashboard';




function BoosterStatusProvider({ children, path }: { path: PageType; children: ReactElement }) {



  const exampleItems: LoadingItem[] = [
    { id: '1', title: 'Configurar ambiente'},
    { id: '2', title: 'Instalar dependências'},
    { id: '3', title: 'Compilar projeto', status:},
    { id: '4', title: 'Executar testes'},
    { id: '5', title: 'Deploy para produção'},
    { id: '6', title: 'Validar deployment'},
    { id: '7', title: 'Notificar equipe'},
    { id: '8', title: 'Atualizar documentação'},
  ]

  return (
    <>
      {children}
      <BoosterStatus 
        path={path}
        items={exampleItems}
        onToggleVisibility={() => console.log('Toggle visibility')}
        boosterQueue={[]}
        completed={0} 
        boostersSelected={false}
        />
    </>
  )
}
export default BoosterStatusProvider