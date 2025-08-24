
import React from 'react';
import { useTranslation } from 'react-i18next';
import { BoosterItem } from '../types/booster.types';
import { BoosterExecution, ExecutionStatus } from '../domain/booster-queue.types';

interface EffectiveBoosterItem extends BoosterItem {
  effectiveState: 'original' | 'staged-apply' | 'staged-revert';
  hasChanges: boolean;
}

interface BoosterCardProps {
  booster: EffectiveBoosterItem;
  execution?: BoosterExecution;
  onToggle: (id: string) => void;
}

export const BoosterCard: React.FC<BoosterCardProps> = ({ 
  booster, 
  execution, 
  onToggle 
}) => {
  const { t } = useTranslation('boosters');

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

  const getExecutionStatusText = (status: ExecutionStatus) => {
    switch (status) {
      case 'queued': return t('status.queued');
      case 'processing': 
        return execution?.operation === 'apply' 
          ? t('status.applying') 
          : t('status.reverting');
      case 'completed': return t('status.completed');
      case 'error': return t('status.error');
      case 'cancelled': return t('status.cancelled');
      default: return '';
    }
  };

  const isExecuting = execution && ['queued', 'processing'].includes(execution.status);
  const showToggle = !isExecuting;

  return (
    <div
      className={`px-5 py-4 bg-zinc-900 border rounded-lg hover:bg-zinc-800/50 transition-colors duration-200 ${
        booster.hasChanges
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

          <p className="pt-1 text-sm text-zinc-500">
            {t(booster.description, { defaultValue: booster.description })}
          </p>

          {booster.tags.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-2">
              {booster.tags.map((tag, index) => (
                <span 
                  key={index} 
                  className="px-1.5 py-0.5 text-xs text-zinc-400 bg-zinc-800 rounded"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Toggle Switch or Status */}
        {showToggle ? (
          <button
            onClick={() => onToggle(booster.id)}
            className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-1 focus:ring-zinc-600 ${
              booster.isApplied ? 'bg-zinc-600' : 'bg-zinc-700'
            } ${booster.hasChanges ? 'ring-1 ring-amber-400/50' : ''}`}
            aria-label={booster.isApplied ? t('filters.status.applied') : t('filters.status.notApplied')}
          >
            <span
              className={`inline-block h-3 w-3 transform rounded-full transition-transform bg-white ${
                booster.isApplied ? 'translate-x-5' : 'translate-x-1'
              }`}
            />
          </button>
        ) : (
          <div className="px-2 py-1 text-xs font-medium rounded bg-blue-400/10 border border-blue-400/20 text-blue-400">
            {getExecutionStatusText(execution!.status)}
          </div>
        )}
      </div>
    </div>
  );
};