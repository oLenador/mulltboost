// features/boosters/types/optimization.types.ts
import { LucideIcon } from 'lucide-react';

export interface OptimizationItem {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  impact: 'low' | 'medium' | 'high';
  advanced?: boolean;
  requiresRestart?: boolean;
}

export interface OptimizationPageConfig {
  title: string;
  description: string;
  icon: LucideIcon;
  optimizations: OptimizationItem[];
}

export interface OptimizationFilters {
  searchTerm: string;
  impactFilter: 'all' | 'low' | 'medium' | 'high';
  statusFilter: 'all' | 'enabled' | 'disabled';
  showAdvanced: boolean;
}

export interface OptimizationStats {
  total: number;
  enabled: number;
  filtered: number;
  requiresRestart: number;
}