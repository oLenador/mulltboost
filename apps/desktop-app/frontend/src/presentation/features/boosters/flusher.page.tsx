// features/boosters/pages/FlusherPage.tsx
import React from 'react';
import { Trash2 } from 'lucide-react';
import { BoosterPageConfig } from './types/booster.types';
import BaseBoosterPage from './compoents/base.componente';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

const flusherConfig: BoosterPageConfig = {
  title: 'pages.flusher.title',
  description: 'pages.flusher.description',
  icon: Trash2,
  category: BoosterCategory.CategoryFlusher
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