// features/boosters/pages/FlusherPage.tsx
import React from 'react';
import { Trash2 } from 'lucide-react';
import { BoosterPageConfig } from '../../types/booster.types';
import BaseBoosterPage from '../components/base-page/base.componente';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

const flusherConfig: BoosterPageConfig = {
  title: 'pages.flusher.title',
  description: 'pages.flusher.description',
  icon: Trash2,
  category: BoosterCategory.CategoryFlusher
};

const FlusherPage: React.FC = () => {
  return (
    <BaseBoosterPage 
      config={flusherConfig}
    />
  );
};

export default FlusherPage;