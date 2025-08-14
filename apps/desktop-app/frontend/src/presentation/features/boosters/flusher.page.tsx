// features/boosters/pages/FlusherPage.tsx
import React from 'react';
import { Trash2 } from 'lucide-react';
import { BoosterPageConfig } from './types/booster.types';
import BaseBoosterPage from './base.componente';

const flusherConfig: BoosterPageConfig = {
  title: 'pages.flusher.title',
  description: 'pages.flusher.description',
  icon: Trash2,
  category: 'flusher'
};

const FlusherPage: React.FC = () => {
  const handleApplyBoosters = (enabledBoosters: any[]) => {
    console.log('Aplicando otimizações de Flusher:', enabledBoosters);
    // Implementar lógica específica para Flusher
  };

  return (
    <BaseBoosterPage 
      config={flusherConfig}
      onApplyBoosters={handleApplyBoosters}
    />
  );
};

export default FlusherPage;