import { BoosterApiService } from '../../core/services/booster-api.service';
import { BoosterDataRepository } from '../../core/repositories/booster-data.repository';
import { BoosterEntity } from '../../core/entities/booster.entity';

export interface LoadBoostersRequest {
  category: string;
  language: string;
}

export interface LoadBoostersResponse {
  success: boolean;
  boosters: BoosterEntity[];
  error?: string;
}

export interface LoadBoostersUseCase {
  execute(request: LoadBoostersRequest): Promise<LoadBoostersResponse>;
}

export function createLoadBoostersUseCase(
  apiService: BoosterApiService,
  repository: BoosterDataRepository
): LoadBoostersUseCase {
  return {
    async execute(request: LoadBoostersRequest): Promise<LoadBoostersResponse> {
      const { category, language } = request;

      try {
        const boosters = await apiService.getBoostersByCategory(category, language);
        repository.setBoosters(boosters);

        return {
          success: true,
          boosters,
        };
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to load boosters';
        
        return {
          success: false,
          boosters: [],
          error: errorMessage,
        };
      }
    },
  };
}