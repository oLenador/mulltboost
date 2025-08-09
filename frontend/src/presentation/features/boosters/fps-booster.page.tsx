// features/boosters/pages/FpsBoostPage.tsx
import React from 'react';
import { Zap } from 'lucide-react';
import { OptimizationPageConfig } from './optimization.types';
import BaseOptimizationPage from './base.componente';

const fpsBoostConfig: OptimizationPageConfig = {
  title: 'FPS Boost',
  description: 'Otimizações para máximo desempenho em jogos',
  icon: Zap,
  optimizations: [
    {
      id: 'gpu-priority',
      name: 'Prioridade Alta da GPU',
      description: 'Define prioridade máxima para processos gráficos',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'cpu-cores',
      name: 'Otimização de Núcleos CPU',
      description: 'Distribui carga entre núcleos de forma eficiente',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'ram-cleanup',
      name: 'Limpeza Automática de RAM',
      description: 'Libera memória não utilizada automaticamente',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'vsync-disable',
      name: 'Desabilitar V-Sync Global',
      description: 'Remove limitação de FPS do V-Sync',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'power-plan',
      name: 'Plano de Energia Alto Desempenho',
      description: 'Configura Windows para máxima performance',
      enabled: true,
      impact: 'high',
      requiresRestart: true
    },
    {
      id: 'background-apps',
      name: 'Suspender Apps em Background',
      description: 'Pausa aplicativos desnecessários durante jogos',
      enabled: true,
      impact: 'medium'
    },
    {
      id: 'cpu-affinity',
      name: 'Otimizar Afinidade de CPU',
      description: 'Reserva núcleos específicos para jogos',
      enabled: false,
      impact: 'high',
      advanced: true
    },
    {
      id: 'memory-compression',
      name: 'Desabilitar Compressão de Memória',
      description: 'Remove overhead da compressão de RAM',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'timer-resolution',
      name: 'Otimizar Resolução de Timer',
      description: 'Melhora precisão de timing do sistema',
      enabled: false,
      impact: 'high',
      advanced: true
    },
    {
      id: 'pci-latency',
      name: 'Reduzir Latência PCI',
      description: 'Otimiza comunicação entre componentes',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'interrupt-policy',
      name: 'Otimizar Política de Interrupções',
      description: 'Melhora resposta do sistema',
      enabled: false,
      impact: 'high',
      advanced: true
    },
    {
      id: 'prefetch-disable',
      name: 'Desabilitar Prefetch',
      description: 'Remove pré-carregamento desnecessário',
      enabled: false,
      impact: 'low'
    }
  ]
};

const FpsBoostPage: React.FC = () => {
  const handleApplyOptimizations = (enabledOptimizations: any[]) => {
    console.log('Aplicando otimizações de FPS Boost:', enabledOptimizations);
    // Implementar lógica específica para FPS Boost
  };

  return (
    <BaseOptimizationPage 
      config={fpsBoostConfig}
      onApplyOptimizations={handleApplyOptimizations}
    />
  );
};

export default FpsBoostPage;