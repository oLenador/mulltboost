// features/boosters/pages/PrecisaoPage.tsx
import React from 'react';
import { Target } from 'lucide-react';
import { OptimizationPageConfig } from './optimization.types';
import BaseOptimizationPage from './base.componente';

const precisionConfig: OptimizationPageConfig = {
  title: 'Precisão',
  description: 'Ajustes finos para mouse e teclado',
  icon: Target,
  optimizations: [
    {
      id: 'mouse-accel',
      name: 'Desabilitar Aceleração do Mouse',
      description: 'Remove aceleração para precisão consistente',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'pointer-precision',
      name: 'Desabilitar Precisão do Ponteiro',
      description: 'Remove ajuste automático de sensibilidade',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'polling-rate',
      name: 'Otimizar Taxa de Polling',
      description: 'Configura taxa máxima de atualização do mouse',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'input-lag',
      name: 'Reduzir Input Lag',
      description: 'Minimiza delay entre input e resposta',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'raw-input',
      name: 'Forçar Raw Input',
      description: 'Bypass do processamento do Windows',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'mouse-smoothing',
      name: 'Desabilitar Suavização do Mouse',
      description: 'Remove filtros de movimento',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'keyboard-repeat',
      name: 'Otimizar Repetição de Teclado',
      description: 'Ajusta delay e taxa de repetição',
      enabled: false,
      impact: 'low'
    },
    {
      id: 'usb-polling',
      name: 'Forçar Polling USB 1000Hz',
      description: 'Força taxa máxima para dispositivos USB',
      enabled: false,
      impact: 'high',
      advanced: true
    },
    {
      id: 'hid-input',
      name: 'Otimizar Driver HID',
      description: 'Melhora comunicação com dispositivos de entrada',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'mouse-fix',
      name: 'MarkC Mouse Fix',
      description: 'Aplica correção avançada de aceleração',
      enabled: false,
      impact: 'high',
      advanced: true
    },
    {
      id: 'input-buffer',
      name: 'Reduzir Buffer de Input',
      description: 'Minimiza buffer de entrada do sistema',
      enabled: false,
      impact: 'medium',
      advanced: true
    }
  ]
};

const PrecisionPage: React.FC = () => {
  const handleApplyOptimizations = (enabledOptimizations: any[]) => {
    console.log('Aplicando otimizações de Precisão:', enabledOptimizations);
    // Implementar lógica específica para Precisão
  };

  return (
    <BaseOptimizationPage 
      config={precisionConfig}
      onApplyOptimizations={handleApplyOptimizations}
    />
  );
};

export default PrecisionPage;