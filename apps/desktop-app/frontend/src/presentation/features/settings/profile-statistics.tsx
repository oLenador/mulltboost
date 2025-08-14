import React from 'react';
import { Calendar, Target, Clock } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/presentation/components/ui/card';

interface Profile {
  stats: {
    joinedAt: string;
    optimizationsApplied: number;
    lastLogin: string;
  };
}

interface ProfileStatisticsProps {
  profile: Profile;
}

export const ProfileStatistics: React.FC<ProfileStatisticsProps> = ({ profile }) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  return (
    <Card variant="zinc">
      <CardHeader>
        <CardTitle className="text-zinc-100">Estatísticas</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center space-x-3">
          <div className="p-2 bg-zinc-800 rounded-lg">
            <Calendar className="w-4 h-4 text-zinc-400" />
          </div>
          <div>
            <p className="text-xs text-zinc-500">Membro desde</p>
            <p className="text-sm text-zinc-300">{formatDate(profile.stats.joinedAt)}</p>
          </div>
        </div>

        <div className="flex items-center space-x-3">
          <div className="p-2 bg-zinc-800 rounded-lg">
            <Target className="w-4 h-4 text-green-400" />
          </div>
          <div>
            <p className="text-xs text-zinc-500">Otimizações aplicadas</p>
            <p className="text-sm text-zinc-300">{profile.stats.optimizationsApplied}</p>
          </div>
        </div>

        <div className="flex items-center space-x-3">
          <div className="p-2 bg-zinc-800 rounded-lg">
            <Clock className="w-4 h-4 text-blue-400" />
          </div>
          <div>
            <p className="text-xs text-zinc-500">Último acesso</p>
            <p className="text-sm text-zinc-300">{formatDate(profile.stats.lastLogin)}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};