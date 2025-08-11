import { OptimizationItem } from "./optimization.types";

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

export const validateOptimization = (optimization: OptimizationItem): boolean => {
  return !!(optimization.id && optimization.name && optimization.description);
};

export const groupOptimizationsByImpact = (optimizations: OptimizationItem[]) => {
  return optimizations.reduce((acc, opt) => {
    if (!acc[opt.impact]) {
      acc[opt.impact] = [];
    }
    acc[opt.impact].push(opt);
    return acc;
  }, {} as Record<string, OptimizationItem[]>);
};

export const getOptimizationStats = (optimizations: OptimizationItem[]) => {
  const total = optimizations.length;
  const enabled = optimizations.filter(opt => opt.enabled).length;
  const advanced = optimizations.filter(opt => opt.advanced).length;
  const requiresRestart = optimizations.filter(opt => opt.requiresRestart).length;
  
  const byImpact = {
    high: optimizations.filter(opt => opt.impact === 'high').length,
    medium: optimizations.filter(opt => opt.impact === 'medium').length,
    low: optimizations.filter(opt => opt.impact === 'low').length
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