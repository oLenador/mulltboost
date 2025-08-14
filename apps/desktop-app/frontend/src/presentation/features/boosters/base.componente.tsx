// features/boosters/components/BaseBoosterPage.tsx
import React from 'react';
import { 
  Search,
  Filter,
  Play,
} from 'lucide-react';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/presentation/components/ui/select';
import { Switch } from '@/presentation/components/ui/switch';
import { Label } from '@/presentation/components/ui/label';
import { BoosterItem, BoosterPageConfig } from './types/booster.types';
import BasePage from '@/presentation/components/pages/base-page';
import { useBoosters } from './use-booster.hook';
import { useTranslation, Trans } from 'react-i18next';

interface BaseBoosterPageProps {
  config: BoosterPageConfig;
  onApplyBoosters?: (enabledBoosters: BoosterItem[]) => void;
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
    applySelectedBoosters,
    enabledCount
  } = useBoosters(config);

  const getImpactColor = (impact: string) => {
    switch (impact) {
      case 'high': return 'text-red-400 bg-red-400/10 border-red-400/20';
      case 'medium': return 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20';
      case 'low': return 'text-green-400 bg-green-400/10 border-green-400/20';
      default: return 'text-zinc-400 bg-zinc-400/10 border-zinc-400/20';
    }
  };

  const getImpactText = (impact: string) => {
    switch (impact) {
      case 'high': return t('impact.badge.high');
      case 'medium': return t('impact.badge.medium');
      case 'low': return t('impact.badge.low');
      default: return t('impact.badge.na');
    }
  };

  const handleApplyBoosters = () => {
    const enabledOpts = boosters.filter(opt => opt.enabled);
    onApplyBoosters?.(enabledOpts);
    applySelectedBoosters();
  };

  const Icon = config.icon;
  const totalCount = boosters.length;

  return (
    <BasePage>
      <>
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center space-x-4 mb-4">
            <div className="p-3 bg-zinc-800 rounded-lg">
              <Icon className="w-6 h-6 text-zinc-400" />
            </div>
            <div>
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
                values={{ enabled: enabledCount, total: totalCount }}
                components={{
                  b1: <span className="font-medium text-green-400" />,
                  b2: <span className="font-medium" />
                }}
              />
            </span>
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
                    <SelectValue placeholder={t('filters.impact.label')} />
                  </SelectTrigger>
                  <SelectContent className="bg-zinc-900 border border-zinc-800">
                    <SelectItem value="all">{t('filters.impact.all')}</SelectItem>
                    <SelectItem value="high">{t('filters.impact.high')}</SelectItem>
                    <SelectItem value="medium">{t('filters.impact.medium')}</SelectItem>
                    <SelectItem value="low">{t('filters.impact.low')}</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-40 bg-zinc-900 border border-zinc-800 rounded-lg text-zinc-300 text-sm focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600">
                  <SelectValue placeholder={t('filters.status.label')} />
                </SelectTrigger>
                <SelectContent className="bg-zinc-900 border border-zinc-800">
                  <SelectItem value="all">{t('filters.status.all')}</SelectItem>
                  <SelectItem value="enabled">{t('filters.status.enabled')}</SelectItem>
                  <SelectItem value="disabled">{t('filters.status.disabled')}</SelectItem>
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
          {filteredBoosters.map((opt) => (
            <div key={opt.id} className="p-4 bg-zinc-900 border border-zinc-800 rounded-lg hover:bg-zinc-800/50 transition-colors duration-200">
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1 mr-4">
                  <div className="flex items-center space-x-2 mb-1">
                    <h4 className="text-sm font-medium text-zinc-300">
                      {t(opt.name, { defaultValue: opt.name })}
                    </h4>
                    <span className={`px-2 py-0.5 text-xs font-medium rounded border ${getImpactColor(opt.impact)}`}>
                      {getImpactText(opt.impact)}
                    </span>
                    {opt.advanced && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-purple-400 bg-purple-400/10 border-purple-400/20">
                        {t('badges.advanced')}
                      </span>
                    )}
                    {opt.requiresRestart && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-orange-400 bg-orange-400/10 border-orange-400/20">
                        {t('badges.requiresRestart')}
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-zinc-500">
                    {t(opt.description, { defaultValue: opt.description })}
                  </p>
                </div>
                <button
                  onClick={() => toggleBooster(opt.id)}
                  className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-1 focus:ring-zinc-600 ${
                    opt.enabled ? 'bg-zinc-600' : 'bg-zinc-700'
                  }`}
                  aria-label={opt.enabled ? t('filters.status.enabled') : t('filters.status.disabled')}
                >
                  <span
                    className={`inline-block h-3 w-3 transform rounded-full bg-white transition-transform ${
                      opt.enabled ? 'translate-x-5' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>
            </div>
          ))}
        </div>

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

        {/* Action Buttons */}
        {filteredBoosters.length > 0 && (
          <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-6">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
              <div className="flex flex-wrap items-center gap-3">
                <Button 
                  onClick={handleApplyBoosters}
                  className="px-6 py-2.5 bg-zinc-800 text-zinc-300 rounded-lg hover:bg-zinc-700 transition-colors duration-200 border border-zinc-700 text-sm font-medium flex items-center space-x-2"
                >
                  <Play className="w-4 h-4" />
                  <span>{t('actions.applySelected')}</span>
                </Button>
                <button 
                  onClick={() => toggleAllBoosters(true)}
                  className="px-4 py-2.5 text-zinc-400 hover:text-zinc-300 transition-colors duration-200 text-sm font-medium"
                >
                  {t('actions.enableAll')}
                </button>
                <button 
                  onClick={() => toggleAllBoosters(false)}
                  className="px-4 py-2.5 text-zinc-400 hover:text-zinc-300 transition-colors duration-200 text-sm font-medium"
                >
                  {t('actions.disableAll')}
                </button>
              </div>
              
              <div className="text-right">
                <div className="text-sm font-medium text-zinc-300">
                  {t('footer.toApply', { count: enabledCount })}
                </div>
                {boosters.filter(opt => opt.enabled && opt.requiresRestart).length > 0 && (
                  <div className="text-xs text-zinc-500">
                    {t('footer.someRequireRestart')}
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </>
    </BasePage>
  );
};

export default BaseBoosterPage;
