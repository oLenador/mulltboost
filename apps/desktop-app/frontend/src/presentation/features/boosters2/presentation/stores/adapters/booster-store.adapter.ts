import { PrimitiveAtom } from 'jotai';
import { BoosterEntity } from '../../../core/entities/booster.entity';

export interface BoosterStoreAdapter {
  updateBoosterStatus: (boosterId: string, isApplied: boolean) => void;
  syncBoosters: (boosters: BoosterEntity[]) => void;
  updateBooster: (boosterId: string, updates: Partial<BoosterEntity>) => void;
  clearBoosters: () => void;
}

export interface BoosterStoreAtoms {
  setBoostersAtom: PrimitiveAtom<(boosters: BoosterEntity[]) => void>;
  updateBoosterAtom: PrimitiveAtom<(boosterId: string, updates: Partial<BoosterEntity>) => void>;
  clearBoostersAtom: PrimitiveAtom<() => void>;
}

export function createBoosterStoreAdapter(
  atoms: BoosterStoreAtoms,
  getAtomValue: <T>(atom: PrimitiveAtom<T>) => T,
  setAtomValue: <T>(atom: PrimitiveAtom<T>, value: T) => void
): BoosterStoreAdapter {
  
  const updateBoosterStatus = (boosterId: string, isApplied: boolean) => {
    const updateFn = getAtomValue(atoms.updateBoosterAtom);
    
    const updates: Partial<BoosterEntity> = {
      isApplied,
      appliedAt: isApplied ? new Date() : undefined,
      revertedAt: !isApplied ? new Date() : undefined,
    };

    updateFn(boosterId, updates);
  };

  const syncBoosters = (boosters: BoosterEntity[]) => {
    const setBoosters = getAtomValue(atoms.setBoostersAtom);
    setBoosters(boosters);
  };

  const updateBooster = (boosterId: string, updates: Partial<BoosterEntity>) => {
    const updateFn = getAtomValue(atoms.updateBoosterAtom);
    updateFn(boosterId, updates);
  };

  const clearBoosters = () => {
    const clearFn = getAtomValue(atoms.clearBoostersAtom);
    clearFn();
  };

  return {
    updateBoosterStatus,
    syncBoosters,
    updateBooster,
    clearBoosters,
  };
}