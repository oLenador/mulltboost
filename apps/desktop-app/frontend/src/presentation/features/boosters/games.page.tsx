// features/boosters/pages/GamesPage.tsx
import React from 'react';
import { Gamepad2 } from 'lucide-react';
import { BoosterPageConfig } from './types/booster.types';
import BaseBoosterPage from './base.componente';

const gamesConfig: BoosterPageConfig = {
  title: 'Games',
  description: 'Configurações específicas para jogos',
  icon: Gamepad2,
  category: "games"
  // boosters: [
  //   {
  //     id: 'game-mode',
  //   
  //     name: 'Modo Jogo Windows',
  //     description: 'Ativa otimizações nativas do Windows',
  //     enabled: true,
  //     impact: 'medium'
  //   },
  //   {
  //     id: 'fullscreen-opt',
  //     name: 'Otimização Tela Cheia',
  //     description: 'Melhora performance em modo fullscreen',
  //     enabled: true,
  //     impact: 'high'
  //   },
  //   {
  //     id: 'dwm-disable',
  //     name: 'Desabilitar DWM em Jogos',
  //     description: 'Remove composição de janelas durante jogos',
  //     enabled: false,
  //     impact: 'high',
  //     requiresRestart: true
  //   },
  //   {
  //     id: 'game-bar',
  //     name: 'Desabilitar Game Bar',
  //     description: 'Remove overlay e gravação automática',
  //     enabled: true,
  //     impact: 'low'
  //   },
  //   {
  //     id: 'focus-assist',
  //     name: 'Assistente de Foco Gaming',
  //     description: 'Bloqueia notificações durante jogos',
  //     enabled: true,
  //     impact: 'low'
  //   },
  //   {
  //     id: 'shader-cache',
  //     name: 'Otimizar Cache de Shaders',
  //     description: 'Pré-compila shaders para loading mais rápido',
  //     enabled: false,
  //     impact: 'medium'
  //   },
  //   {
  //     id: 'exclusive-fullscreen',
  //     name: 'Forçar Fullscreen Exclusivo',
  //     description: 'Força modo exclusivo para melhor performance',
  //     enabled: false,
  //     impact: 'high'
  //   },
  //   {
  //     id: 'gpu-scheduling',
  //     name: 'Hardware-Accelerated GPU Scheduling',
  //     description: 'Ativa agendamento de GPU por hardware',
  //     enabled: false,
  //     impact: 'medium',
  //     requiresRestart: true
  //   },
  //   {
  //     id: 'game-dvr',
  //     name: 'Desabilitar Game DVR',
  //     description: 'Remove gravação automática de jogos',
  //     enabled: true,
  //     impact: 'medium'
  //   },
  //   {
  //     id: 'steam-overlay',
  //     name: 'Otimizar Overlay Steam',
  //     description: 'Reduz impacto do overlay da Steam',
  //     enabled: false,
  //     impact: 'low'
  //   },
  //   {
  //     id: 'directx-booster',
  //     name: 'Otimizações DirectX',
  //     description: 'Aplica tweaks específicos para DirectX',
  //     enabled: false,
  //     impact: 'high',
  //     advanced: true
  //   },
  //   {
  //     id: 'vulkan-layers',
  //     name: 'Desabilitar Layers Vulkan',
  //     description: 'Remove layers desnecessários do Vulkan',
  //     enabled: false,
  //     impact: 'low',
  //     advanced: true
  //   }
  // ]
};

const GamesPage: React.FC = () => {
  const handleApplyBoosters = (enabledBoosters: any[]) => {
    console.log('Aplicando otimizações de Games:', enabledBoosters);
    // Implementar lógica específica para Games
  };

  return (
    <BaseBoosterPage 
      config={gamesConfig}
      onApplyBoosters={handleApplyBoosters}
    />
  );
};

export default GamesPage;