import { useCallback } from 'react';
import { useAtom } from 'jotai';
import { BoosterEntity } from '../../core/entities/booster.entity';
import { 
  boostersAtom, 
  boosterByIdAtom, 
  setBoostersAtom, 
  updateBoosterAtom, 
  clearBoostersAtom 
} from '../stores/booster-data.store';

export const useBoosterData = () => {
  const [boosters] = useAtom(boostersAtom);
  const [boosterById] = useAtom(boosterByIdAtom);
  const [, setBoosters] = useAtom(setBoostersAtom);
  const [, updateBooster] = useAtom(updateBoosterAtom);
  const [, clearBoosters] = useAtom(clearBoostersAtom);

  const getBooster = useCallback((id: string): BoosterEntity | undefined => {
    return boosterById.get(id);
  }, [boosterById]);

  const getAppliedBoosters = useCallback((): BoosterEntity[] => {
    return boosters.filter(booster => booster.isApplied);
  }, [boosters]);

  const getBoostersByCategory = useCallback((category: string): BoosterEntity[] => {
    return boosters.filter(booster => booster.category === category);
  }, [boosters]);

  const searchBoosters = useCallback((searchTerm: string): BoosterEntity[] => {
    if (!searchTerm.trim()) return boosters;
    
    const searchLower = searchTerm.toLowerCase();
    return boosters.filter(booster =>
      booster.name.toLowerCase().includes(searchLower) ||
      booster.description.toLowerCase().includes(searchLower) ||
      booster.tags.some(tag => tag.toLowerCase().includes(searchLower))
    );
  }, [boosters]);

  return {
    boosters,
    boosterById,
    getBooster,
    getAppliedBoosters,
    getBoostersByCategory,
    searchBoosters,
    setBoosters,
    updateBooster,
    clearBoosters,
  };
};