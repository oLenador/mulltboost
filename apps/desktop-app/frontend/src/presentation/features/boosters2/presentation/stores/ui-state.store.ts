import { atom } from 'jotai';

// UI filters and search
export const searchTermAtom = atom<string>('');
export const impactFilterAtom = atom<string>('all');
export const statusFilterAtom = atom<string>('all');
export const showAdvancedAtom = atom<boolean>(false);

// Loading states
export const loadingBoostersAtom = atom<boolean>(false);
export const executingBatchAtom = atom<boolean>(false);
export const syncingAtom = atom<boolean>(false);

// Actions
export const setSearchTermAtom = atom(
  null,
  (get, set, searchTerm: string) => {
    set(searchTermAtom, searchTerm);
  }
);

export const setImpactFilterAtom = atom(
  null,
  (get, set, filter: string) => {
    set(impactFilterAtom, filter);
  }
);

export const setStatusFilterAtom = atom(
  null,
  (get, set, filter: string) => {
    set(statusFilterAtom, filter);
  }
);

export const toggleAdvancedAtom = atom(
  null,
  (get, set) => {
    const current = get(showAdvancedAtom);
    set(showAdvancedAtom, !current);
  }
);

export const setLoadingBoostersAtom = atom(
  null,
  (get, set, loading: boolean) => {
    set(loadingBoostersAtom, loading);
  }
);

export const setExecutingBatchAtom = atom(
  null,
  (get, set, executing: boolean) => {
    set(executingBatchAtom, executing);
  }
);

export const setSyncingAtom = atom(
  null,
  (get, set, syncing: boolean) => {
    set(syncingAtom, syncing);
  }
);