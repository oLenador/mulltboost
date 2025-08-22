import React from 'react';
import {
  Search,
  Filter,
  RotateCcw,
  AlertCircle,
} from 'lucide-react';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/presentation/components/ui/select';
import { Switch } from '@/presentation/components/ui/switch';
import { Label } from '@/presentation/components/ui/label';
import { BoosterItem, BoosterPageConfig } from '../types/booster.types';
import BasePage from '@/presentation/components/pages/base-page';
import { useBoosters } from '../hooks/use-booster.hook';
import { useTranslation, Trans } from 'react-i18next';
import { FloatElement } from '@/presentation/components/floating-manager';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/presentation/components/ui/tooltip';

interface BaseBoosterPageProps {
  config: BoosterPageConfig;
  onApplyBoosters?: (appliedBoosters: BoosterItem[]) => void;
}

const BaseBoosterPage: React.FC<BaseBoosterPageProps> = ({ config, onApplyBoosters }) => {
  const { t } = useTranslation('boosters');

  const {
    filteredBoosters,
    boosters,
    searchTerm, setSearchTerm,
    impactFilter, setImpactFilter,
    statusFilter, setStatusFilter,
    showAdvanced, setShowAdvanced,
    toggleBooster,
    toggleAllBoosters,
    resetChanges,
    appliedCount,
    changesCount,
    hasChanges,
    loading
  } = useBoosters(config);

  const getRiskLevelColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'high': return 'text-red-400 bg-red-400/10 border-red-400/20';
      case 'medium': return 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20';
      case 'low': return 'text-green-400 bg-green-400/10 border-green-400/20';
      default: return 'text-zinc-400 bg-zinc-400/10 border-zinc-400/20';
    }
  };

  const getRiskLevelText = (riskLevel: string) => {
    switch (riskLevel) {
      case 'high': return t('impact.badge.high');
      case 'medium': return t('impact.badge.medium');
      case 'low': return t('impact.badge.low');
      default: return t('impact.badge.na');
    }
  };

  // const getEffectiveStateIndicator = (booster: any) => {
  //   if (!booster.hasChanges) return null;
  // 
  //   return (
  //     <div className="flex items-center space-x-1 text-xs">
  //       <AlertCircle className="w-3 h-3 text-amber-400" />
  //       <span className="text-amber-400 font-medium">
  //         {booster.effectiveState === 'staged-apply' ? t('status.willApply') : t('status.willRevert')}
  //       </span>
  //     </div>
  //   );
  // };

  const Icon = config.icon;
  const totalCount = boosters.length;

  if (loading) {
    return (
      <BasePage>
        <div className="flex items-center justify-center py-12">
          <div className="text-zinc-400">Loading boosters...</div>
        </div>
      </BasePage>
    );
  }

  return (
    <BasePage>
      <>
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center space-x-4 mb-4">
            <div className="p-3 bg-zinc-800 rounded-lg">
              <Icon className="w-6 h-6 text-zinc-400" />
            </div>
            <div className="flex-1">
              <h1 className="text-2xl font-semibold text-zinc-100">
                {t(config.title, { defaultValue: config.title })}
              </h1>
              <p className="text-zinc-400 text-sm">
                {t(config.description, { defaultValue: config.description })}
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-4 text-sm">
            <span className="text-zinc-300">
              <Trans
                i18nKey="header.enabledOfTotalActive"
                values={{ enabled: appliedCount ?? 0, total: totalCount ?? 0 }}
                components={{
                  b1: <span className="font-medium text-green-400" />,
                  b2: <span className="font-medium" />
                }}
              />
            </span>

            {hasChanges && (
              <>
                <span className="text-zinc-500">•</span>
                <span className="text-amber-400 font-medium">
                  {t('header.pendingChanges', { count: changesCount })}
                </span>
              </>
            )}

            <span className="text-zinc-500">•</span>
            <span className="text-zinc-400">
              {t('header.results', { count: filteredBoosters.length })}
            </span>
          </div>
        </div>

        {/* Search and Filters */}
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

        {/* Boosters Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-8">
          {filteredBoosters.map((booster) => (
            <div
              key={booster.id}
              className={`px-5 py-4 bg-zinc-900 border rounded-lg hover:bg-zinc-800/50 transition-colors duration-200 ${booster.hasChanges
                ? 'border-amber-400/30 bg-amber-400/5'
                : 'border-zinc-800'
                }`}
            >
              <div className="flex items-start justify-between mb-1">
                <div className="flex-1 mr-4">
                  <div className="flex items-center space-x-2 mb-1">
                    <h4 className="text-sm font-medium text-zinc-300">
                      {t(booster.name, { defaultValue: booster.name })}
                    </h4>
                    <span className={`px-2 py-0.5 text-xs font-medium rounded border ${getRiskLevelColor(booster.riskLevel)}`}>
                      {getRiskLevelText(booster.riskLevel)}
                    </span>
                    {!booster.reversible && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-orange-400 bg-orange-400/10 border-orange-400/20">
                        {t('badges.nonReversible')}
                      </span>
                    )}
                    {booster.dependencies.length > 0 && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-blue-400 bg-blue-400/10 border-blue-400/20">
                        {t('badges.hasDependencies')}
                      </span>
                    )}
                  </div>

                  {/* Effective State Indicator */}
                  {
                    // getEffectiveStateIndicator(booster)
                  }

                  <p className="pt-1 text-sm text-zinc-500">
                    {t(booster.description, { defaultValue: booster.description })}
                  </p>

                  {booster.tags.length > 0 && (
                    <div className="flex flex-wrap gap-1 mt-2">
                      {booster.tags.map((tag, index) => (
                        <span key={index} className="px-1.5 py-0.5 text-xs text-zinc-400 bg-zinc-800 rounded">
                          {tag}
                        </span>
                      ))}
                    </div>
                  )}
                </div>

                {/* Toggle Switch */}
                <button
                  onClick={() => toggleBooster(booster.id)}
                  className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-1 focus:ring-zinc-600 ${booster.isApplied ? 'bg-zinc-600' : 'bg-zinc-700'
                    } ${booster.hasChanges ? 'ring-1 ring-amber-400/50' : ''}`}
                  aria-label={booster.isApplied ? t('filters.status.applied') : t('filters.status.notApplied')}
                >
                  <span
                    className={`inline-block h-3 w-3 transform rounded-full transition-transform ${booster.isApplied ? 'translate-x-5 bg-white' : 'translate-x-1 bg-white'
                      } ${booster.hasChanges ? 'bg-white' : 'bg-white'}`}
                  />
                </button>
              </div>
            </div>
          ))}
        </div>
        {/* Reset Changes Button */}
        {hasChanges && (
          <FloatElement
            id={'selector-shurtcut'}
            type={'custom'}
            position='bottom-right'
            className='right-[112px]'
          >
            <div className='flex flex-row gap-2 bg-zinc-800/50 p-3 w-fit rounded-xl border border-zinc-700/60'>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <Button
                      onClick={resetChanges}
                      variant="zinc"
                      size="icon"
                      className="group border-zinc-700/60 hover:bg-zinc-800"
                    >
                      <RotateCcw size={24} className='min-w-5 min-h-5 text-zinc-300/60 group-hover:text-zinc-300'/>
                    </Button>
                    <TooltipContent>
                      {t('actions.resetChanges')}
                    </TooltipContent>
                  </TooltipTrigger>
                </Tooltip>
              </TooltipProvider>

              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger>
                    <Button
                      onClick={resetChanges}
                      variant="zinc"
                      size="icon"
                      className="group border-zinc-700/60 hover:bg-zinc-800"
                    >
                      <RotateCcw size={24} className='min-w-5 min-h-5 text-zinc-300/60 group-hover:text-zinc-300'/>
                    </Button>
                    <TooltipContent>
                      {t('actions.resetChanges')}
                    </TooltipContent>
                  </TooltipTrigger>
                </Tooltip>
              </TooltipProvider>
              
            </div>
          </FloatElement>
        )}
        {/* Empty State */}
        {filteredBoosters.length === 0 && (
          <div className="text-center py-12">
            <div className="w-16 h-16 bg-zinc-800 rounded-lg flex items-center justify-center mx-auto mb-4">
              <Search className="w-8 h-8 text-zinc-500" />
            </div>
            <h3 className="text-lg font-medium text-zinc-300 mb-2">
              {t('empty.title')}
            </h3>
            <p className="text-zinc-500 text-sm">
              {t('empty.subtitle')}
            </p>
          </div>
        )}
      </>
    </BasePage>
  );
};

export default BaseBoosterPage;