import React from 'react';
import {
    Brain,
    TrendingUp,
    Activity,
    Sparkles
} from 'lucide-react';

// shadcn/ui components (presume sua stack já tem estes componentes)
import { Card, CardContent } from '@/presentation/components/ui/card';
import { Button } from '@/presentation/components/ui/button';
import BasePage from '@/presentation/components/pages/base-page';


const SmartBoost: React.FC = () => {
    return (
        <BasePage>
            <div className='max-w-6xl'>
                <div className="mb-8">
                    <h1 className="text-2xl font-semibold text-zinc-100 mb-1">Smart Boost</h1>
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
        </BasePage>
    );
};

export default SmartBoost