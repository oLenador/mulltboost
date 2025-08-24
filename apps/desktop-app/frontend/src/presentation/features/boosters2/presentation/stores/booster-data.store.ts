import { atom } from 'jotai';
import { BoosterEntity } from '../../core/entities/booster.entity';

export const boostersAtom = atom<BoosterEntity[]>([]);

export const boosterByIdAtom = atom(
  (get) => {
    const boosters = get(boostersAtom);
    return new Map(boosters.map(booster => [booster.id, booster]));
  }
);

// Actions
export const setBoostersAtom = atom(
  null,
  (get, set, boosters: BoosterEntity[]) => {
    set(boostersAtom, boosters);
  }
);

export const updateBoosterAtom = atom(
  null,
  (get, set, { id, updates }: { id: string; updates: Partial<BoosterEntity> }) => {
    const boosters = get(boostersAtom);
    const updatedBoosters = boosters.map(booster =>
      booster.id === id ? { ...booster, ...updates } : booster
    );
    set(boostersAtom, updatedBoosters);
  }
);

export const clearBoostersAtom = atom(
  null,
  (get, set) => {
    set(boostersAtom, []);
  }
);