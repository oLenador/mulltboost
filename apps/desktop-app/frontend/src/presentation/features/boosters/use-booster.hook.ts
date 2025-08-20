// features/boosters/hooks/useBoosters.ts
import { useEffect, useState, useMemo, useCallback } from 'react';
import { BoosterItem, BoosterPageConfig } from './types/booster.types';
import { useAtom } from 'jotai';
import { userDataAtom } from '@/core/store/user-data.store';
import { GetBoostersByCategory } from 'bindings/github.com/oLenador/mulltbost/internal/app/handlers/boosterhandler';
import { Language } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/services/i18n';
import { BoosterDto } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/dto';

function mapRiskLevelToImpact(riskLevel: string): 'low' | 'medium' | 'high' {
  switch ((riskLevel || '').toLowerCase()) {
    case 'low': return 'low';
    case 'medium': return 'medium';
    case 'high': return 'high';
    default: return 'low';
  }
}

function mapBoosterToItem(opt: BoosterDto): BoosterItem {
  return {
    id: opt.ID,
    name: opt.Name,
    description: opt.Description,
    enabled: false,
    impact: mapRiskLevelToImpact(opt.RiskLevel),
    advanced: opt.Level?.toLowerCase() === 'advanced',
    requiresRestart: opt.Dependencies?.includes('restart') || false,
  };
}

export function useBoosters(config: BoosterPageConfig) {
  const [searchTerm, setSearchTerm] = useState('');
  const [impactFilter, setImpactFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [boosters, setBoosters] = useState<BoosterItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [userData] = useAtom(userDataAtom)


  useEffect(() => {
    const loadBoosters = async () => {
      setLoading(true);
      try {
        
        const res = await GetBoostersByCategory(config.category, userData.language as Language);
        console.log(res)
        setBoosters(res.map(mapBoosterToItem));
      } catch (err) {
        console.error('Erro ao buscar otimizações:', err);
      } finally {
        setLoading(false);
      }
    };
    loadBoosters();
  }, [config.category]);

  const toggleBooster = useCallback((id: string) => {
    setBoosters(prev =>
      prev.map(opt =>
        opt.id === id ? { ...opt, enabled: !opt.enabled } : opt
      )
    );
  }, []);

  const toggleAllBoosters = useCallback((enable: boolean) => {
    setBoosters(prev => prev.map(opt => ({ ...opt, enabled: enable })));
  }, []);

  const applySelectedBoosters = useCallback(async () => {
    console.log("trying to apply boosters ")
    const ids = boosters.filter(opt => opt.enabled).map(opt => opt.id);
    if (ids.length === 0) return null;

    try {
      const result = {} // await ApplyBoosterBatch(ids);
      console.log('Resultado da aplicação:', result);
      return result;
    } catch (err) {
      console.error('Erro ao aplicar boosters:', err);
      throw err;
    }
  }, [boosters]);

  const filteredBoosters = useMemo(() => {
    return boosters.filter(opt => {
      if (searchTerm && !opt.name.toLowerCase().includes(searchTerm.toLowerCase()) &&
          !opt.description.toLowerCase().includes(searchTerm.toLowerCase())) {
        return false;
      }
      if (impactFilter !== 'all' && opt.impact !== impactFilter) {
        return false;
      }
      if (statusFilter === 'enabled' && !opt.enabled) return false;
      if (statusFilter === 'disabled' && opt.enabled) return false;
      if (!showAdvanced && opt.advanced) return false;
      return true;
    });
  }, [boosters, searchTerm, impactFilter, statusFilter, showAdvanced]);

  const enabledCount = useMemo(
    () => boosters.filter(opt => opt.enabled).length,
    [boosters]
  );

  return {
    boosters,
    filteredBoosters,
    searchTerm, setSearchTerm,
    impactFilter, setImpactFilter,
    statusFilter, setStatusFilter,
    showAdvanced, setShowAdvanced,
    toggleBooster,
    toggleAllBoosters,
    enabledCount,
    loading,
    applySelectedBoosters
  };
}
