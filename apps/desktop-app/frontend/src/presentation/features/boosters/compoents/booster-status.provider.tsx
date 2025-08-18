

import React, { ReactElement } from 'react'
import BoosterStatus from './booster-status.component'
import { LoadingItem } from '../hooks/booster-status.hook';




function BoosterStatusProvider({ children, path }: { path: string; children: ReactElement }) {
  // Exemplo de items para demonstração
  const exampleItems: LoadingItem[] = [
    { id: '1', title: 'Configurar ambiente', status: 'applied' },
    { id: '2', title: 'Instalar dependências', status: 'applied' },
    { id: '3', title: 'Compilar projeto', status: 'applying', progress: 65 },
    { id: '4', title: 'Executar testes', status: 'idle' },
    { id: '5', title: 'Deploy para produção', status: 'idle' },
    { id: '6', title: 'Validar deployment', status: 'error' },
    { id: '7', title: 'Notificar equipe', status: 'idle' },
    { id: '8', title: 'Atualizar documentação', status: 'idle' }
  ]

  return (
    <>
      {children}
      <BoosterStatus 
        path={path} 
        items={exampleItems}
        onToggleVisibility={() => console.log('Toggle visibility')}
      />
    </>
  )
}
export default BoosterStatusProvider