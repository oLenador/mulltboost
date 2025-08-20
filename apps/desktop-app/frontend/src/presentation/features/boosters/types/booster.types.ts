// features/boosters/types/optimization.types.ts
import { BoosterCategory } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { LucideIcon } from 'lucide-react';
import { ReactElement } from 'react';

export interface ProcessableItem {
  id: string;
  name: string;
  type: string;
}

export interface BoosterItem extends ProcessableItem {
  description: string;
  enabled: boolean;
  impact: 'low' | 'medium' | 'high';
  advanced?: boolean;
  requiresRestart?: boolean;
  type: 'booster'; 
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
  icon: LucideIcon;
  category: BoosterCategory
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