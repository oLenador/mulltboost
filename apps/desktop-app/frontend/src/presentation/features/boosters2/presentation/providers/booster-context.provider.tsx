import React, { createContext, useContext, ReactNode } from 'react';
import { BoosterDataRepository } from '../../core/repositories/booster-data.repository';
import { ExecutionStateRepository } from '../../core/repositories/execution-state.repository';
import { LoadBoostersUseCase } from '../../application/use-cases/load-boosters.use-case';
import { ManageStagingUseCase } from '../../application/use-cases/manage-staging.use-case';

export interface BoosterContextValue {
  boosterRepository: BoosterDataRepository;
  executionRepository: ExecutionStateRepository;
  loadBoostersUseCase: LoadBoostersUseCase;
  manageStagingUseCase: ManageStagingUseCase;
}

const BoosterContext = createContext<BoosterContextValue | null>(null);

export interface BoosterContextProviderProps {
  children: ReactNode;
  boosterRepository: BoosterDataRepository;
  executionRepository: ExecutionStateRepository;
  loadBoostersUseCase: LoadBoostersUseCase;
  manageStagingUseCase: ManageStagingUseCase;
}

export function BoosterContextProvider({
  children,
  boosterRepository,
  executionRepository,
  loadBoostersUseCase,
  manageStagingUseCase,
}: BoosterContextProviderProps) {
  const contextValue: BoosterContextValue = {
    boosterRepository,
    executionRepository,
    loadBoostersUseCase,
    manageStagingUseCase,
  };

  return (
    <BoosterContext.Provider value={contextValue}>
      {children}
    </BoosterContext.Provider>
  );
}

export function useBoosterContext(): BoosterContextValue {
  const context = useContext(BoosterContext);
  if (!context) {
    throw new Error('useBoosterContext must be used within BoosterContextProvider');
  }
  return context;
}
