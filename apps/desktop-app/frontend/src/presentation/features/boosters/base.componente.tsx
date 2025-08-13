// features/boosters/components/BaseBoosterPage.tsx
import React, { useEffect, useState } from 'react';
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
import { GetSystemMetrics } from 'wailsjs/go/handlers/MonitoringHandler';
import { entities } from 'wailsjs/go/models';

interface BaseBoosterPageProps {
  config: BoosterPageConfig;
  onApplyBoosters?: (enabledBoosters: BoosterItem[]) => void;
}

function mapBoosterToItem(opt: entities.BoosterDto): BoosterItem {
  return {
    id: opt.ID,
    name: opt.Name,
    description: opt.Description,
    enabled: false,
    impact: mapRiskLevelToImpact(opt.RiskLevel),
    advanced: opt.Level?.toLowerCase() === 'advanced',
    requiresRestart: opt.Dependencies?.includes('restart') || false,
  };
}

function mapRiskLevelToImpact(riskLevel: string): 'low' | 'medium' | 'high' {
  switch ((riskLevel || '').toLowerCase()) {
    case 'low': return 'low';
    case 'medium': return 'medium';
    case 'high': return 'high';
    default: return 'low';
  }
}

const BaseBoosterPage: React.FC<BaseBoosterPageProps> = ({ 
  config, 
  onApplyBoosters 
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [impactFilter, setImpactFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [boosters, setBoosters] = useState<BoosterItem[]>([]);

  useEffect(() => {
    const loadBoosters = async () => {
      try {
        const res = await GetBoostersByCategorBooster(config.category);
        setBoosters(res.map(mapBoosterToItem));
      } catch (err) {
        console.error('Erro ao buscar otimizações:', err);
      }
    };
    loadBoosters();
  }, [config.category]);


  useEffect(() => {
    
    console.log(GetSystemMetrics())
  }, [])

  const toggleBooster = (id: string) => {
    setBoosters(prev => 
      prev.map(opt => 
        opt.id === id ? { ...opt, enabled: !opt.enabled } : opt
      )
    );
  };

  const toggleAllBoosters = (enable: boolean) => {
    setBoosters(prev => 
      prev.map(opt => ({ ...opt, enabled: enable }))
    );
  };

  const handleApplyBoosters = () => {
    const enabledOpts = boosters.filter(opt => opt.enabled);
    if (onApplyBoosters) {
      onApplyBoosters(enabledOpts);
    }
    console.log('Aplicando otimizações:', enabledOpts);
  };

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
      case 'high': return 'Alto';
      case 'medium': return 'Médio';
      case 'low': return 'Baixo';
      default: return 'N/A';
    }
  };

  // Filter boosters
  const filteredBoosters = boosters.filter(opt => {
    // Search filter
    if (searchTerm && !opt.name.toLowerCase().includes(searchTerm.toLowerCase()) && 
        !opt.description.toLowerCase().includes(searchTerm.toLowerCase())) {
      return false;
    }
    
    // Impact filter
    if (impactFilter !== 'all' && opt.impact !== impactFilter) {
      return false;
    }
    
    // Status filter
    if (statusFilter === 'enabled' && !opt.enabled) return false;
    if (statusFilter === 'disabled' && opt.enabled) return false;
    
    // Advanced filter
    if (!showAdvanced && opt.advanced) return false;
    
    return true;
  });

  const Icon = config.icon;
  const enabledCount = boosters.filter(opt => opt.enabled).length;
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
              <h1 className="text-2xl font-semibold text-zinc-100">{config.title}</h1>
              <p className="text-zinc-400 text-sm">{config.description}</p>
            </div>
          </div>
          
          <div className="flex items-center space-x-4 text-sm">
            <span className="text-zinc-300">
              <span className="font-medium text-green-400">{enabledCount}</span> de <span className="font-medium">{totalCount}</span> otimizações ativas
            </span>
            <span className="text-zinc-500">•</span>
            <span className="text-zinc-400">
              {filteredBoosters.length} {filteredBoosters.length === 1 ? 'resultado' : 'resultados'}
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
                placeholder="Pesquisar otimizações..."
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
                    <SelectValue placeholder="Impacto" />
                  </SelectTrigger>
                  <SelectContent className="bg-zinc-900 border border-zinc-800">
                    <SelectItem value="all">Todos os Impactos</SelectItem>
                    <SelectItem value="high">Alto Impacto</SelectItem>
                    <SelectItem value="medium">Médio Impacto</SelectItem>
                    <SelectItem value="low">Baixo Impacto</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-40 bg-zinc-900 border border-zinc-800 rounded-lg text-zinc-300 text-sm focus:ring-1 focus:ring-zinc-600 focus:border-zinc-600">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent className="bg-zinc-900 border border-zinc-800">
                  <SelectItem value="all">Todos os Status</SelectItem>
                  <SelectItem value="enabled">Apenas Ativas</SelectItem>
                  <SelectItem value="disabled">Apenas Inativas</SelectItem>
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
                  Mostrar Avançadas
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
                    <h4 className="text-sm font-medium text-zinc-300">{opt.name}</h4>
                    <span className={`px-2 py-0.5 text-xs font-medium rounded border ${getImpactColor(opt.impact)}`}>
                      {getImpactText(opt.impact)}
                    </span>
                    {opt.advanced && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-purple-400 bg-purple-400/10 border-purple-400/20">
                        Avançada
                      </span>
                    )}
                    {opt.requiresRestart && (
                      <span className="px-2 py-0.5 text-xs font-medium rounded border text-orange-400 bg-orange-400/10 border-orange-400/20">
                        Requer Reinício
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-zinc-500">{opt.description}</p>
                </div>
                <button
                  onClick={() => toggleBooster(opt.id)}
                  className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none focus:ring-1 focus:ring-zinc-600 ${
                    opt.enabled ? 'bg-zinc-600' : 'bg-zinc-700'
                  }`}
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
            <h3 className="text-lg font-medium text-zinc-300 mb-2">Nenhuma otimização encontrada</h3>
            <p className="text-zinc-500 text-sm">Tente ajustar os filtros ou termo de pesquisa</p>
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
                  <span>Aplicar Selecionadas</span>
                </Button>
                <button 
                  onClick={() => toggleAllBoosters(true)}
                  className="px-4 py-2.5 text-zinc-400 hover:text-zinc-300 transition-colors duration-200 text-sm font-medium"
                >
                  Ativar Todas
                </button>
                <button 
                  onClick={() => toggleAllBoosters(false)}
                  className="px-4 py-2.5 text-zinc-400 hover:text-zinc-300 transition-colors duration-200 text-sm font-medium"
                >
                  Desativar Todas
                </button>
              </div>
              
              <div className="text-right">
                <div className="text-sm font-medium text-zinc-300">
                  {enabledCount} otimizações serão aplicadas
                </div>
                {boosters.filter(opt => opt.enabled && opt.requiresRestart).length > 0 && (
                  <div className="text-xs text-zinc-500">
                    Algumas otimizações requerem reinício
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