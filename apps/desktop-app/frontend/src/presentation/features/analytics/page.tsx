import React from 'react';
import {
  TrendingUp,
  Cpu,
  Activity,
} from 'lucide-react';

// shadcn/ui components (presume sua stack já tem estes componentes)
import { Card, CardContent } from '@/presentation/components/ui/card';
import { Badge } from '@/presentation/components/ui/badge';
import BasePage from '@/presentation/components/pages/base-page';


const AnalyticsPage: React.FC = () => {
  const recs = [
    {
      title: 'Otimizar inicialização do sistema',
      impact: 'Alto',
      description: 'Desativar 7 programas que iniciam com o Windows',
      color: 'text-red-400'
    },
    {
      title: 'Limpeza de arquivos temporários',
      impact: 'Médio',
      description: 'Liberar 2.3GB de espaço em disco',
      color: 'text-yellow-400'
    },
    {
      title: 'Atualizar drivers gráficos',
      impact: 'Alto',
      description: 'Nova versão disponível com melhorias de performance',
      color: 'text-red-400'
    },
    {
      title: 'Configurar modo de energia',
      impact: 'Baixo',
      description: 'Alterar para modo de alta performance',
      color: 'text-green-400'
    }
  ];

  return (
    <BasePage>
    <>
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-zinc-100 mb-1">Análises IA</h1>
          <p className="text-zinc-400 text-sm">Insights inteligentes sobre o desempenho do sistema</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-8">
          <Card className="bg-zinc-900 border-zinc-800 rounded-lg p-6">
            <CardContent>
              <div className="flex items-center space-x-3 mb-4">
                <div className="p-2 bg-zinc-800 rounded-lg">
                  <TrendingUp className="w-4 h-4 text-zinc-400" />
                </div>
                <h3 className="text-sm font-medium text-zinc-300">Score de Performance</h3>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold text-green-400 mb-2">87</div>
                <div className="w-full bg-zinc-800 rounded-full h-2">
                  <div className="bg-green-400 h-2 rounded-full" style={{ width: '87%' }}></div>
                </div>
                <p className="text-xs text-zinc-500 mt-2">Excelente desempenho</p>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-zinc-900 border-zinc-800 rounded-lg p-6">
            <CardContent>
              <div className="flex items-center space-x-3 mb-4">
                <div className="p-2 bg-zinc-800 rounded-lg">
                  <Cpu className="w-4 h-4 text-zinc-400" />
                </div>
                <h3 className="text-sm font-medium text-zinc-300">Otimizações Sugeridas</h3>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold text-blue-400 mb-2">5</div>
                <p className="text-xs text-zinc-500">Melhorias identificadas</p>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-zinc-900 border-zinc-800 rounded-lg p-6">
            <CardContent>
              <div className="flex items-center space-x-3 mb-4">
                <div className="p-2 bg-zinc-800 rounded-lg">
                  <Activity className="w-4 h-4 text-zinc-400" />
                </div>
                <h3 className="text-sm font-medium text-zinc-300">Economia de Energia</h3>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold text-purple-400 mb-2">23%</div>
                <p className="text-xs text-zinc-500">Redução no consumo</p>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="bg-zinc-900 border-zinc-800 rounded-lg p-6">
            <CardContent>
              <h3 className="text-lg font-medium text-zinc-100 mb-6">Recomendações da IA</h3>
              <div className="space-y-3">
                {recs.map((rec, index) => (
                  <div key={index} className="p-3 bg-zinc-800/50 border border-zinc-700/50 rounded-lg">
                    <div className="flex items-start justify-between mb-2">
                      <h4 className="text-sm font-medium text-zinc-300">{rec.title}</h4>
                      <Badge className={`px-2 py-0.5 text-xs font-medium rounded ${rec.color} bg-zinc-800`}>{rec.impact}</Badge>
                    </div>
                    <p className="text-xs text-zinc-500">{rec.description}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          <Card className="bg-zinc-900 border-zinc-800 rounded-lg p-6">
            <CardContent>
              <h3 className="text-lg font-medium text-zinc-100 mb-6">Tendências de Uso</h3>
              <div className="h-64 bg-zinc-800 border border-zinc-700 rounded-lg flex items-center justify-center">
                <p className="text-zinc-500 text-sm">Gráfico de Tendências IA</p>
              </div>
            </CardContent>
          </Card>
        </div>
        </>
        </BasePage>
  );
};

export default AnalyticsPage