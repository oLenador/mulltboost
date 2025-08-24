export interface BoosterEntity {
    readonly id: string;
    readonly name: string;
    readonly description: string;
    readonly category: string;
    readonly level: string;
    readonly platform: string;
    readonly dependencies: string[];
    readonly conflicts: string[];
    readonly reversible: boolean;
    readonly riskLevel: 'low' | 'medium' | 'high';
    readonly version: string;
    readonly isApplied: boolean;
    readonly appliedAt: Date;
    readonly revertedAt: Date;
    readonly tags: string[];
  }
  
  export function createBoosterEntity(data: Partial<BoosterEntity> & { id: string }): BoosterEntity {
    return {
      id: data.id,
      name: data.name || '',
      description: data.description || '',
      category: data.category || '',
      level: data.level || '',
      platform: data.platform || '',
      dependencies: data.dependencies || [],
      conflicts: data.conflicts || [],
      reversible: data.reversible ?? true,
      riskLevel: data.riskLevel || 'low',
      version: data.version || '1.0.0',
      isApplied: data.isApplied ?? false,
      appliedAt: data.appliedAt || new Date(0),
      revertedAt: data.revertedAt || new Date(0),
      tags: data.tags || [],
    };
  }
  