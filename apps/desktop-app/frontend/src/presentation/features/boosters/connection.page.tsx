import React from 'react';
import { Wifi } from 'lucide-react';
import BaseBoosterPage from './base.componente';
import { BoosterPageConfig } from './types/booster.types';

const connectionConfig: BoosterPageConfig = {
  title: 'pages.connection.title',
  description: 'pages.connection.description',
  icon: Wifi,
  category: 'connection'
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