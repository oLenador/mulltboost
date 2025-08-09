import React, { useState } from 'react';
import {
  MessageSquare,
  BarChart3,
  Brain,
  Send,
  Bot,
  User,
  TrendingUp,
  Cpu,
  Activity,
  Sparkles
} from 'lucide-react';

// shadcn/ui components (presume sua stack já tem estes componentes)
import { Card, CardHeader, CardTitle, CardContent } from '@/presentation/components/ui/card';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/presentation/components/ui/tabs';
import { Badge } from '@/presentation/components/ui/badge';

interface MultiAIProps {
  activeSection?: string;
}

// --- Sub-páginas (preservando o estilo original, agora reaproveitando componentes shadcn) ---

const HomePage: React.FC = () => {
  return (
    <div className="p-8 bg-zinc-950 min-h-screen text-zinc-100">
      <div className="max-w-4xl mx-auto text-center">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-zinc-100 mb-1">Otimização Inteligente</h1>
          <p className="text-zinc-400 text-sm">IA avançada que aprende com seu uso e otimiza automaticamente</p>
        </div>

        <Card className="bg-zinc-900 border-zinc-800 mb-8">
          <CardContent className="p-8">
            <div className="w-16 h-16 bg-zinc-800 rounded-lg flex items-center justify-center mx-auto mb-6">
              <Brain className="w-8 h-8 text-zinc-400" />
            </div>
            <h2 className="text-xl font-medium text-zinc-100 mb-4">Sistema de IA Adaptativa</h2>
            <p className="text-zinc-400 mb-8 text-sm">
              Nossa IA monitora continuamente seu sistema, aprende seus padrões de uso e aplica otimizações personalizadas automaticamente.
            </p>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
              <div className="p-4 bg-zinc-800/50 border border-zinc-700/50 rounded-lg">
                <Sparkles className="w-6 h-6 text-zinc-400 mx-auto mb-3" />
                <h3 className="text-sm font-medium text-zinc-300 mb-2">Aprendizado Contínuo</h3>
                <p className="text-xs text-zinc-500">Analisa seu comportamento para otimizações personalizadas</p>
              </div>
              <div className="p-4 bg-zinc-800/50 border border-zinc-700/50 rounded-lg">
                <Activity className="w-6 h-6 text-zinc-400 mx-auto mb-3" />
                <h3 className="text-sm font-medium text-zinc-300 mb-2">Monitoramento 24/7</h3>
                <p className="text-xs text-zinc-500">Supervisiona o desempenho em tempo real</p>
              </div>
              <div className="p-4 bg-zinc-800/50 border border-zinc-700/50 rounded-lg">
                <TrendingUp className="w-6 h-6 text-zinc-400 mx-auto mb-3" />
                <h3 className="text-sm font-medium text-zinc-300 mb-2">Otimização Automática</h3>
                <p className="text-xs text-zinc-500">Aplica melhorias sem intervenção manual</p>
              </div>
            </div>

            <div className="text-center">
              <Button className="px-6 py-2.5 bg-zinc-800 text-zinc-300 rounded-lg font-medium hover:bg-zinc-700 transition-colors duration-200 border border-zinc-700 flex items-center space-x-2 mx-auto text-sm">
                <Brain className="w-4 h-4" />
                <span>Ativar IA Adaptativa</span>
              </Button>
            </div>
          </CardContent>
        </Card>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Card className="bg-zinc-900 border-zinc-800 p-6">
            <CardContent>
              <h3 className="text-lg font-medium text-zinc-100 mb-4">Status da IA</h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-zinc-400 text-sm">Modelo Treinado</span>
                  <span className="text-green-400 text-sm font-medium">Ativo</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-zinc-400 text-sm">Dados Coletados</span>
                  <span className="text-zinc-300 text-sm font-medium">2.4GB</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-zinc-400 text-sm">Precisão</span>
                  <span className="text-zinc-300 text-sm font-medium">94.3%</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-zinc-400 text-sm">Última Atualização</span>
                  <span className="text-zinc-300 text-sm font-medium">2h atrás</span>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="bg-zinc-900 border-zinc-800 p-6">
            <CardContent>
              <h3 className="text-lg font-medium text-zinc-100 mb-4">Próximas Ações</h3>
              <div className="space-y-3">
                <div className="p-3 bg-zinc-800/50 border border-zinc-700/50 rounded-lg text-left">
                  <p className="text-sm font-medium text-blue-400">Otimização de RAM agendada</p>
                  <p className="text-xs text-zinc-500">Em 2 horas</p>
                </div>
                <div className="p-3 bg-zinc-800/50 border border-zinc-700/50 rounded-lg text-left">
                  <p className="text-sm font-medium text-green-400">Análise de jogos detectados</p>
                  <p className="text-xs text-zinc-500">Quando iniciar jogo</p>
                </div>
                <div className="p-3 bg-zinc-800/50 border border-zinc-700/50 rounded-lg text-left">
                  <p className="text-sm font-medium text-purple-400">Relatório semanal</p>
                  <p className="text-xs text-zinc-500">Em 3 dias</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

const ChatPage: React.FC = () => {
  const [chatMessage, setChatMessage] = useState('');
  const [messages] = useState([
    {
      id: 1,
      type: 'bot',
      content: 'Olá! Sou sua assistente de otimização IA. Como posso ajudar você hoje?',
      timestamp: '10:30'
    },
    {
      id: 2,
      type: 'user',
      content: 'Meu jogo está com fps baixo, o que posso fazer?',
      timestamp: '10:32'
    },
    {
      id: 3,
      type: 'bot',
      content: 'Analisando seu sistema... Recomendo executar o FPS Boost e fechar aplicativos em background. Detectei que há 12 processos desnecessários consumindo 23% da CPU.',
      timestamp: '10:32'
    }
  ]);

  return (
    <div className="p-8 bg-zinc-950 min-h-screen text-zinc-100">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-zinc-100 mb-1">Chat IA</h1>
          <p className="text-zinc-400 text-sm">Converse com a assistente inteligente de otimização</p>
        </div>

        <Card className="bg-zinc-900 border-zinc-800 rounded-lg h-[600px] flex flex-col">
          <div className="p-4 border-b border-zinc-800">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-zinc-800 rounded-full flex items-center justify-center">
                <Bot className="w-4 h-4 text-zinc-400" />
              </div>
              <div>
                <h3 className="text-sm font-medium text-zinc-300">Assistente OptiMax IA</h3>
                <p className="text-xs text-green-400">● Online</p>
              </div>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {messages.map((message) => (
              <div key={message.id} className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div className={`max-w-xs lg:max-w-md flex items-start space-x-3 ${message.type === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`}>
                  <div className={`w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0 ${
                    message.type === 'user'
                      ? 'bg-zinc-700'
                      : 'bg-zinc-800'
                  }`}>
                    {message.type === 'user' ? (
                      <User className="w-3 h-3 text-zinc-400" />
                    ) : (
                      <Bot className="w-3 h-3 text-zinc-400" />
                    )}
                  </div>
                  <div className={`px-3 py-2 rounded-lg ${
                    message.type === 'user'
                      ? 'bg-zinc-800 text-zinc-300 border border-zinc-700'
                      : 'bg-zinc-800/50 text-zinc-300 border border-zinc-700/50'
                  }`}>
                    <p className="text-sm">{message.content}</p>
                    <p className="text-xs text-zinc-500 mt-1">{message.timestamp}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="p-4 border-t border-zinc-800">
            <div className="flex space-x-3">
              <Input
                value={chatMessage}
                onChange={(e: any) => setChatMessage(e.target.value)}
                placeholder="Digite sua mensagem..."
                className="flex-1 px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg text-zinc-300 placeholder-zinc-500 focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600 text-sm"
              />
              <Button className="px-4 py-2 bg-zinc-800 text-zinc-300 rounded-lg hover:bg-zinc-700 transition-colors duration-200 border border-zinc-700 flex items-center space-x-2 text-sm">
                <Send className="w-3 h-3" />
                <span>Enviar</span>
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

const AnalisesPage: React.FC = () => {
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
    <div className="p-8 bg-zinc-950 min-h-screen text-zinc-100">
      <div className="max-w-7xl mx-auto">
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
      </div>
    </div>
  );
};

// --- Componente principal que separa as páginas ---

const MultiAI: React.FC<MultiAIProps> = ({ activeSection = 'home' }) => {
  const [tab, setTab] = useState<string>(activeSection || 'home');

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100">
      <header className="p-4 border-b border-zinc-800">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-zinc-900 rounded-lg flex items-center justify-center">
              <Brain className="w-5 h-5 text-zinc-400" />
            </div>
            <div>
              <h1 className="text-lg font-semibold">OptiMax</h1>
              <p className="text-xs text-zinc-400">Otimização Inteligente</p>
            </div>
          </div>

          <nav>
            <Tabs value={tab} onValueChange={(v) => setTab(v)}>
              <TabsList className="bg-transparent border-none p-0">
                <TabsTrigger value="home" className={`px-3 py-1 rounded ${tab === 'home' ? 'bg-zinc-900 border border-zinc-800' : ''}`}>Início</TabsTrigger>
                <TabsTrigger value="chat" className={`px-3 py-1 rounded ${tab === 'chat' ? 'bg-zinc-900 border border-zinc-800' : ''}`}>Chat</TabsTrigger>
                <TabsTrigger value="analises" className={`px-3 py-1 rounded ${tab === 'analises' ? 'bg-zinc-900 border border-zinc-800' : ''}`}>Análises</TabsTrigger>
              </TabsList>
            </Tabs>
          </nav>
        </div>
      </header>

      <main>
        {tab === 'home' && <HomePage />}
        {tab === 'chat' && <ChatPage />}
        {tab === 'analises' && <AnalisesPage />}
      </main>

    </div>
  );
};

export default MultiAI;
