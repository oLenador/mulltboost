import React from 'react';
import { Wifi } from 'lucide-react';
import { BoosterPageConfig } from '../../types/booster.types';
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import BaseBoosterPage from '../components/base-page/base.componente';

const connectionConfig: BoosterPageConfig = {
  title: 'pages.connection.title',
  description: 'pages.connection.description',
  icon: Wifi,
  category: BoosterCategory.CategoryConnection
};

const ConnectionPage: React.FC = () => {
  return (
    <BaseBoosterPage 
      config={connectionConfig}
    />
  );
};

export default ConnectionPage;