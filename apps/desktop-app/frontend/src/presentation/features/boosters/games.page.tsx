// features/boosters/pages/GamesPage.tsx
import React from 'react';
import { Gamepad2 } from 'lucide-react';
import { BoosterPageConfig } from './types/booster.types';
import BaseBoosterPage from './base.componente';

const gamesConfig: BoosterPageConfig = {
  title: 'pages.games.title',
  description: 'pages.games.description',
  icon: Gamepad2,
  category: 'games'
};

const GamesPage: React.FC = () => {
  const handleApplyBoosters = (enabledBoosters: any[]) => {
    console.log('Aplicando otimizações de Games:', enabledBoosters);
  };

  return (
    <BaseBoosterPage 
      config={gamesConfig}
      onApplyBoosters={handleApplyBoosters}
    />
  );
};

export default GamesPage;