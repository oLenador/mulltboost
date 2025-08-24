import { BoosterEntity } from '../entities/booster.entity';

export interface BoosterDataRepository {
  // Booster data
  getAllBoosters(): BoosterEntity[];
  getBooster(id: string): BoosterEntity | undefined;
  setBoosters(boosters: BoosterEntity[]): void;
  updateBooster(id: string, updates: Partial<BoosterEntity>): void;
  clearBoosters(): void;

  // Filtering and search
  getBoostersByCategory(category: string): BoosterEntity[];
  searchBoosters(searchTerm: string): BoosterEntity[];
  getBoostersByStatus(applied: boolean): BoosterEntity[];
  getBoostersByRiskLevel(riskLevel: 'low' | 'medium' | 'high'): BoosterEntity[];
}