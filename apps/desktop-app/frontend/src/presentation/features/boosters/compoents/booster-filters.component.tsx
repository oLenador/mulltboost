// src/presentation/features/boosters/components/BoosterFilters.tsx

import React from 'react';
import { Search, Filter } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Input } from '@/presentation/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/presentation/components/ui/select';
import { Switch } from '@/presentation/components/ui/switch';
import { Label } from '@/presentation/components/ui/label';

interface BoosterFiltersProps {
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  impactFilter: string;
  setImpactFilter: (filter: string) => void;
  statusFilter: string;
  setStatusFilter: (filter: string) => void;
  showAdvanced: boolean;
  setShowAdvanced: (show: boolean) => void;
}

export const BoosterFilters: React.FC<BoosterFiltersProps> = ({
  searchTerm,
  setSearchTerm,
  impactFilter,
  setImpactFilter,
  statusFilter,
  setStatusFilter,
  showAdvanced,
  setShowAdvanced,
}) => {
  const { t } = useTranslation('boosters');

  return (
    <div className="mb-6 space-y-4">
      <div className="flex flex-col lg:flex-row lg:items-center lg:space-x-4 space-y-4 lg:space-y-0">
        {/* Search */}
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-zinc-500" />
          <Input
            type="text"
            placeholder={t('search.placeholder')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-zinc-900 border border-zinc-800 rounded-lg text-zinc-300 placeholder-zinc-500 focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600 text-sm"
          />
        </div>

        {/* Filters */}
        <div className="flex items-center space-x-3">
          <div className="flex items-center space-x-2">
            <Filter className="w-4 h-4 text-zinc-500" />
            <Select value={impactFilter} onValueChange={setImpactFilter}>
              <SelectTrigger className="w-40 bg-zinc-900 border border-zinc-800 rounded-lg text-zinc-300 text-sm focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600">
                <SelectValue placeholder={t('filters.risk.label')} />
              </SelectTrigger>
              <SelectContent className="bg-zinc-900 border border-zinc-800">
                <SelectItem value="all">{t('filters.risk.all')}</SelectItem>
                <SelectItem value="high">{t('filters.risk.high')}</SelectItem>
                <SelectItem value="medium">{t('filters.risk.medium')}</SelectItem>
                <SelectItem value="low">{t('filters.risk.low')}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-40 bg-zinc-900 border border-zinc-800 rounded-lg text-zinc-300 text-sm focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600">
              <SelectValue placeholder={t('filters.status.label')} />
            </SelectTrigger>
            <SelectContent className="bg-zinc-900 border border-zinc-800">
              <SelectItem value="all">{t('filters.status.all')}</SelectItem>
              <SelectItem value="enabled">{t('filters.status.applied')}</SelectItem>
              <SelectItem value="disabled">{t('filters.status.notApplied')}</SelectItem>
            </SelectContent>
          </Select>

          <label className="flex items-center space-x-2 text-sm">
            <Switch
              id="show-advanced"
              checked={showAdvanced}
              onCheckedChange={setShowAdvanced}
              className="data-[state=checked]:bg-zinc-600 data-[state=unchecked]:bg-zinc-700"
            />
            <Label htmlFor="show-advanced" className="text-zinc-400">
              {t('filters.showAdvanced')}
            </Label>
          </label>
        </div>
      </div>
    </div>
  );
};