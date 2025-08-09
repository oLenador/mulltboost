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
import { Card } from '@/presentation/components/ui/card';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import BasePage from '@/presentation/components/pages/base-page';


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
    <BasePage>
    <div className='max-w-3xl'>
        <div className="mb-8">
          <h1 className="text-2xl font-semibold text-zinc-100 mb-1">Chat IA</h1>
          <p className="text-zinc-400 text-sm">Converse com a assistente inteligente de otimização</p>
        </div>

        <Card className="!bg-zinc-900 border-zinc-800 rounded-lg h-[600px] p-0 flex flex-col">
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
        </BasePage>
  );
};

export default ChatPage


