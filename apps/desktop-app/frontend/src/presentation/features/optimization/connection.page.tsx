import React from 'react';
import { Wifi } from 'lucide-react';
import BaseOptimizationPage from './base.componente';
import { OptimizationPageConfig } from './optimization.types';

const connectionConfig: OptimizationPageConfig = {
  title: 'Conexão',
  description: 'Melhorias de velocidade e estabilidade de rede',
  icon: Wifi,
  optimizations: [
    {
      id: 'dns-cloudflare',
      name: 'DNS Cloudflare (1.1.1.1)',
      description: 'Configura DNS mais rápido e seguro',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'dns-google',
      name: 'DNS Google (8.8.8.8)',
      description: 'Utiliza servidores DNS do Google',
      enabled: true,
      impact: 'medium'
    },
    {
      id: 'tcp-optimizer',
      name: 'Otimização TCP/IP',
      description: 'Ajusta parâmetros de rede para menor latência',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'nagle-disable',
      name: 'Desabilitar Algoritmo de Nagle',
      description: 'Reduz delay em conexões TCP',
      enabled: false,
      impact: 'medium'
    },
    {
      id: 'qos-gaming',
      name: 'QoS para Gaming',
      description: 'Prioriza tráfego de jogos na rede',
      enabled: true,
      impact: 'high'
    },
    {
      id: 'ipv6-disable',
      name: 'Desabilitar IPv6',
      description: 'Remove overhead do protocolo IPv6',
      enabled: false,
      impact: 'low'
    },
    {
      id: 'network-throttling',
      name: 'Desabilitar Throttling de Rede',
      description: 'Remove limitações artificiais de velocidade',
      enabled: false,
      impact: 'high'
    },
    {
      id: 'tcp-chimney',
      name: 'TCP Chimney Offload',
      description: 'Transfere processamento TCP para placa de rede',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'rss-scaling',
      name: 'Receive Side Scaling (RSS)',
      description: 'Distribui processamento de rede entre CPUs',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'interrupt-moderation',
      name: 'Moderação de Interrupções de Rede',
      description: 'Otimiza interrupções da placa de rede',
      enabled: false,
      impact: 'low',
      advanced: true
    },
    {
      id: 'buffer-sizes',
      name: 'Otimizar Tamanhos de Buffer',
      description: 'Ajusta buffers de recepção e transmissão',
      enabled: false,
      impact: 'medium',
      advanced: true
    },
    {
      id: 'dns-cache-size',
      name: 'Aumentar Cache DNS',
      description: 'Expande cache de resolução DNS',
      enabled: false,
      impact: 'low'
    }
  ]
};

const ConnectionPage: React.FC = () => {
  const handleApplyOptimizations = (enabledOptimizations: any[]) => {
    console.log('Aplicando otimizações de Conexão:', enabledOptimizations);
    // Implementar lógica específica para Conexão
  };

  return (
    <BaseOptimizationPage 
      config={connectionConfig}
      onApplyOptimizations={handleApplyOptimizations}
    />
  );
};

export default ConnectionPage;