// features/boosters/types/optimization.types.ts
import { LucideIcon } from 'lucide-react';

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