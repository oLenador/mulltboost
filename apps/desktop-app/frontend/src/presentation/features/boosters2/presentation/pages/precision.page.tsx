// features/boosters/pages/PrecisaoPage.tsx
import React from 'react';
import { Target } from 'lucide-react';
import { BoosterPageConfig } from '../../types/booster.types';
import BaseBoosterPage from '../components/base-page/base.componente';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

const precisionConfig: BoosterPageConfig = {
  title: 'Precisão',
  description: 'Ajustes finos para mouse e teclado',
  icon: Target,
category: BoosterCategory.CategoryPrecision
};

const PrecisionPage: React.FC = () => {
  return (
    <BaseBoosterPage 
      config={precisionConfig}
    />
  );
};

export default PrecisionPage;