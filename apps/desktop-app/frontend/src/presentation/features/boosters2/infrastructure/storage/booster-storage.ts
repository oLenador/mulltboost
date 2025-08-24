import { BoosterDataRepository } from '../../core/repositories/booster-data.repository';
import { BoosterEntity } from '../../core/entities/booster.entity';

export function createInMemoryBoosterRepository(): BoosterDataRepository {
  let boosters: BoosterEntity[] = [];

  return {
    getAllBoosters(): BoosterEntity[] {
      return [...boosters];
    },

    getBooster(id: string): BoosterEntity | undefined {
      return boosters.find(b => b.id === id);
    },

    setBoosters(newBoosters: BoosterEntity[]): void {
      boosters = [...newBoosters];
    },

    updateBooster(id: string, updates: Partial<BoosterEntity>): void {
      const index = boosters.findIndex(b => b.id === id);
      if (index !== -1) {
        boosters[index] = { ...boosters[index], ...updates };
      }
    },

    clearBoosters(): void {
      boosters = [];
    },

    getBoostersByCategory(category: string): BoosterEntity[] {
      return boosters.filter(b => b.category === category);
    },

    searchBoosters(searchTerm: string): BoosterEntity[] {
      const searchLower = searchTerm.toLowerCase();
      return boosters.filter(b => 
        b.name.toLowerCase().includes(searchLower) ||
        b.description.toLowerCase().includes(searchLower) ||
        b.tags.some(tag => tag.toLowerCase().includes(searchLower))
      );
    },

    getBoostersByStatus(applied: boolean): BoosterEntity[] {
      return boosters.filter(b => b.isApplied === applied);
    },

    getBoostersByRiskLevel(riskLevel: 'low' | 'medium' | 'high'): BoosterEntity[] {
      return boosters.filter(b => b.riskLevel === riskLevel);
    },
  };
}