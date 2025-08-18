// features/boosters/types/optimization.types.ts
import { LucideIcon } from 'lucide-react';
import { ReactElement } from 'react';

interface BoosterPipeline {
  pipeline
  startAt
}

interface BoosterPipeItem {
  id: 
}

export interface BoosterItem {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  impact: 'low' | 'medium' | 'high';
  advanced?: boolean;
  requiresRestart?: boolean;
}

export interface BoosterPageConfig {
  title: string;
  description: string;
  icon: ReactElement;
  category: string
}

export interface BoosterFilters {
  searchTerm: string;
  impactFilter: 'all' | 'low' | 'medium' | 'high';
  statusFilter: 'all' | 'enabled' | 'disabled';
  showAdvanced: boolean;
}

export interface BoosterStats {
  total: number;
  enabled: number;
  filtered: number;
  requiresRestart: number;
}