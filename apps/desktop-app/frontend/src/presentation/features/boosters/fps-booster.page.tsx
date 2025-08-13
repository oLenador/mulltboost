import React from 'react';
import { Zap } from 'lucide-react';
import BaseBoosterPage from './base.componente';
import { BoosterPageConfig } from './types/booster.types';

const fpsBoostConfig: Omit<BoosterPageConfig, 'boosters'> & { category: string } = {
  title: 'FPS Boost',
  description: 'Otimizações para máximo desempenho em jogos',
  icon: Zap,
  category: 'performance' 
};

const FpsBoostPage: React.FC = () => {
  const handleApplyBoosters = (enabledBoosters: any[]) => {
    console.log('Aplicando otimizações de FPS Boost:', enabledBoosters);
  };

  return (
    <BaseBoosterPage 
      config={fpsBoostConfig}
      onApplyBoosters={handleApplyBoosters}
    />
  );
};

export default FpsBoostPage;
