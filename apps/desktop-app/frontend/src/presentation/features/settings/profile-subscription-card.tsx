import React from 'react';
import { Crown, Star, Shield, Check } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/presentation/components/ui/card';
import { Badge } from '@/presentation/components/ui/badge';

interface Profile {
  subscription: {
    plan: string;
    status: string;
    expiresAt?: string;
    features: string[];
  };
}

interface ProfileSubscriptionProps {
  profile: Profile;
}

export const ProfileSubscription: React.FC<ProfileSubscriptionProps> = ({ profile }) => {
  const getSubscriptionIcon = () => {
    switch (profile.subscription.plan) {
      case 'premium': return Crown;
      case 'pro': return Star;
      default: return Shield;
    }
  };

  const getSubscriptionColor = () => {
    switch (profile.subscription.plan) {
      case 'premium': return 'text-yellow-400';
      case 'pro': return 'text-blue-400';
      default: return 'text-zinc-400';
    }
  };

  const getPlanName = () => {
    switch (profile.subscription.plan) {
      case 'premium': return 'Premium';
      case 'pro': return 'Pro';
      default: return 'Free';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  const SubscriptionIcon = getSubscriptionIcon();

  return (
    <Card variant="zinc">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2 text-zinc-100">
          <SubscriptionIcon className={`w-5 h-5 ${getSubscriptionColor()}`} />
          <span>Assinatura</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <div className="flex items-center justify-between mb-2">
            <span className="text-zinc-300">Plano</span>
            <Badge className={getSubscriptionColor()}>{getPlanName()}</Badge>
          </div>
          <div className="flex items-center justify-between mb-2">
            <span className="text-zinc-300">Status</span>
            <Badge variant={profile.subscription.status === 'active' ? 'default' : 'destructive'}>
              {profile.subscription.status === 'active' ? 'Ativo' : 'Inativo'}
            </Badge>
          </div>
          {profile.subscription.expiresAt && (
            <div className="flex items-center justify-between">
              <span className="text-zinc-300">Expira em</span>
              <span className="text-zinc-400 text-sm">
                {formatDate(profile.subscription.expiresAt)}
              </span>
            </div>
          )}
        </div>

        <div className="pt-4 border-t border-zinc-800">
          <h4 className="text-sm font-medium text-zinc-300 mb-2">Recursos inclusos:</h4>
          <ul className="space-y-1">
            {profile.subscription.features.map((feature, index) => (
              <li key={index} className="flex items-center space-x-2 text-xs">
                <Check className="w-3 h-3 text-green-400" />
                <span className="text-zinc-400">{feature}</span>
              </li>
            ))}
          </ul>
        </div>
      </CardContent>
    </Card>
  );
};