import React from 'react';
import { Zap } from 'lucide-react';
import BaseBoosterPage from '../components/base-page/base.componente';
import { BoosterPageConfig } from '../../types/booster.types';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

const fpsBoostConfig: BoosterPageConfig = {
  title: 'pages.fpsboost.title',
  description: 'pages.fpsboost.description',
  icon: Zap,
  category: BoosterCategory.CategoryFPSBooster
};

const FpsBoostPage: React.FC = () => {
  return (
    <BaseBoosterPage 
      config={fpsBoostConfig}
    />
  );
};

export default FpsBoostPage;
