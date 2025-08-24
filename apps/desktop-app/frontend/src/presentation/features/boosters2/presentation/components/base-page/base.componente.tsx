// src/presentation/features/boosters/components/BaseBoosterPage.tsx

import React, { useCallback } from 'react';
import { Search } from 'lucide-react';
import { useTranslation, Trans } from 'react-i18next';
import BasePage from '@/presentation/components/pages/base-page';
import { BoosterFilters } from '../booster-filters.component';

import { BoosterPageConfig } from '../../../types/booster.types';
import { VirtualizedBoosterList } from './virtualized-booster-list.component';
import { useBoosters } from './use-base-component.hook';
import { useBoosterExecution } from '../status/use-booster-status-component.hook';
import { BoosterExecutionControls } from './booster-execution-control.component';


interface BaseBoosterPageProps {
  config: BoosterPageConfig;
}

const BaseBoosterPage: React.FC<BaseBoosterPageProps> = ({ config }) => {
  const { t } = useTranslation('boosters');

  // UI state and booster management
  const {
    filteredBoosters,
    boosters,
    searchTerm,
    setSearchTerm,
    impactFilter,
    setImpactFilter,
    statusFilter,
    setStatusFilter,
    showAdvanced,
    setShowAdvanced,
    toggleBooster,
    resetChanges,
    appliedCount,
    changesCount,
    hasChanges,
    loading,
  } = useBoosters(config);

  // Execution management
  const {
    executions,
    stats,
    isExecuting,
    canExecute,
    stagedCount,
    executeStagedBatch,
    cancelStagedExecutions,
    resetExecution,
  } = useBoosterExecution();

  // Convert executions array to map for efficient lookups
  const executionsMap = React.useMemo(() => {
    return new Map(executions.map(exec => [exec.boosterId, exec]));
  }, [executions]);

  const handleReset = useCallback(() => {
    resetChanges();
    resetExecution();
  }, [resetChanges, resetExecution]);

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

          {/* Stats */}
          <div className="flex items-center space-x-4 text-sm">
            <span className="text-zinc-300">
              <Trans
                i18nKey="header.enabledOfTotalActive"
                values={{ enabled: appliedCount ?? 0, total: totalCount ?? 0 }}
                components={{
                  b1: <span className="font-medium text-green-400" />,
                  b2: <span className="font-medium" />,
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

            {isExecuting && (
              <>
                <span className="text-zinc-500">•</span>
                <span className="text-blue-400 font-medium">
                  {t('header.executing', {
                    processing: stats.processing,
                    queued: stats.queued
                  })}
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
        <BoosterFilters
          searchTerm={searchTerm}
          setSearchTerm={setSearchTerm}
          impactFilter={impactFilter}
          setImpactFilter={setImpactFilter}
          statusFilter={statusFilter}
          setStatusFilter={setStatusFilter}
          showAdvanced={showAdvanced}
          setShowAdvanced={setShowAdvanced}
        />

        {/* Boosters List */}
        {filteredBoosters.length > 0 ? (
          <div className="mb-8">
            <VirtualizedBoosterList
              boosters={filteredBoosters}
              executions={executionsMap}
              onToggleBooster={toggleBooster}
              height={Math.min(600, window.innerHeight - 400)}
            />
          </div>
        ) : (
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

        <BoosterExecutionControls
          hasChanges={hasChanges}
          onReset={handleReset}
        />
      </>
    </BasePage>
  );
};

export default BaseBoosterPage;