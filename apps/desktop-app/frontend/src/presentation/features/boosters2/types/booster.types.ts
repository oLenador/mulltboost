import { BoosterCategory, BoosterLevel } from 'bindings/github.com/oLenador/mulltbost/internal/core/domain/entities';
import { LucideIcon } from 'lucide-react';
import { ReactElement } from 'react';

type BoosterOperationType = "apply" | "revert"

export interface BoosterItem {
  id: string;
  name: string;
  description: string;
  category: BoosterCategory;
  level: BoosterLevel;
  platform: string;
  dependencies: string[];
  conflicts: string[];
  reversible: boolean;
  riskLevel: 'low' | 'medium' | 'high';
  version: string;
  isApplied: boolean;
  appliedAt: Date;
  revertedAt: Date;
  tags: string[];
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