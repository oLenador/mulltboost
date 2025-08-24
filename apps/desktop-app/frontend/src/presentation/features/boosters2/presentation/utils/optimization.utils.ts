import { BoosterItem } from "./types/booster.types";

export const getImpactVariant = (impact: string): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (impact) {
    case 'high': return 'destructive';
    case 'medium': return 'default';
    case 'low': return 'secondary';
    default: return 'outline';
  }
};

export const getImpactText = (impact: string): string => {
  switch (impact) {
    case 'high': return 'Alto';
    case 'medium': return 'Médio';
    case 'low': return 'Baixo';
    default: return 'N/A';
  }
};

export const validateBooster = (booster: BoosterItem): boolean => {
  return !!(booster.id && booster.name && booster.description);
};

export const groupBoostersByImpact = (boosters: BoosterItem[]) => {
  return boosters.reduce((acc, opt) => {
    if (!acc[opt.impact]) {
      acc[opt.impact] = [];
    }
    acc[opt.impact].push(opt);
    return acc;
  }, {} as Record<string, BoosterItem[]>);
};

export const getBoosterStats = (boosters: BoosterItem[]) => {
  const total = boosters.length;
  const enabled = boosters.filter(opt => opt.enabled).length;
  const advanced = boosters.filter(opt => opt.advanced).length;
  const requiresRestart = boosters.filter(opt => opt.requiresRestart).length;
  
  const byImpact = {
    high: boosters.filter(opt => opt.impact === 'high').length,
    medium: boosters.filter(opt => opt.impact === 'medium').length,
    low: boosters.filter(opt => opt.impact === 'low').length
  };

  return {
    total,
    enabled,
    advanced,
    requiresRestart,
    byImpact,
    enabledPercentage: total > 0 ? Math.round((enabled / total) * 100) : 0
  };
};