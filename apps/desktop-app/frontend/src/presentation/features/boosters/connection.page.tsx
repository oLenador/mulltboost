import React from 'react';
import { Wifi } from 'lucide-react';
import BaseBoosterPage from './compoents/base.componente';
import { BoosterPageConfig } from './types/booster.types';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';

const connectionConfig: BoosterPageConfig = {
  title: 'pages.connection.title',
  description: 'pages.connection.description',
  icon: Wifi,
  category: BoosterCategory.CategoryConnection
};

const ConnectionPage: React.FC = () => {
  const handleApplyBoosters = (enabledBoosters: any[]) => {
    // Lógica específica opcional
  };

  return (
    <BaseBoosterPage 
      config={connectionConfig}
      onApplyBoosters={handleApplyBoosters}
    />
  );
};

export default ConnectionPage;