import { apiClient } from '../../../core/api/client';
import type { 
  OptimizationCategory, 
  OptimizationResult, 
  BatchResult 
} from '../../../core/api/types';

export class OptimizationService {
  // Convenience methods for specific categories
  static async applyPrecisionOptimizations(ids: string[]): Promise<BatchResult> {
    return apiClient.applyOptimizationBatch(ids);
  }

  static async applyPerformanceOptimizations(ids: string[]): Promise<BatchResult> {
    return apiClient.applyOptimizationBatch(ids);
  }

  static async applyNetworkOptimizations(ids: string[]): Promise<BatchResult> {
    return apiClient.applyOptimizationBatch(ids);
  }

  // Validation helpers
  static async validateOptimizationBatch(ids: string[]): Promise<{ valid: string[], invalid: string[] }> {
    const optimizations = await apiClient.getAvailableOptimizations();
    const optimizationMap = new Map(optimizations.map(opt => [opt.id, opt]));
    
    const valid: string[] = [];
    const invalid: string[] = [];
    
    for (const id of ids) {
      const opt = optimizationMap.get(id);
      if (opt) {
        const state = await apiClient.getOptimizationState(id);
        if (!state.applied) {
          valid.push(id);
        } else {
          invalid.push(id);
        }
      } else {
        invalid.push(id);
      }
    }
    
    return { valid, invalid };
  }

  // Risk assessment
  static assessBatchRisk(optimizations: Array<{ id: string, riskLevel: string }>): {
    totalRisk: 'low' | 'medium' | 'high',
    riskBreakdown: Record<string, number>
  } {
    const riskBreakdown = { low: 0, medium: 0, high: 0 };
    
    optimizations.forEach(opt => {
      riskBreakdown[opt.riskLevel as keyof typeof riskBreakdown]++;
    });
    
    let totalRisk: 'low' | 'medium' | 'high' = 'low';
    if (riskBreakdown.high > 0) {
      totalRisk = 'high';
    } else if (riskBreakdown.medium > 0) {
      totalRisk = 'medium';
    }
    
    return { totalRisk, riskBreakdown };
  }
}