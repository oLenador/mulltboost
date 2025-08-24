import { BoosterEntity } from '../entities/booster.entity';
import { ApiResponse, QueueStateData } from '../types/api.types';

export interface BoosterApiService {
  executeBooster(boosterId: string, operation: 'apply' | 'revert'): Promise<ApiResponse>;
  getQueueState(): Promise<QueueStateData>;
  getBoostersByCategory(category: string, language: string): Promise<BoosterEntity[]>;
}