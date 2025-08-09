import { UserProfile, UpdateProfileRequest, ApiResponse } from "./profile.types";

const mockUserProfile: UserProfile = {
  id: 'user_123',
  name: 'Carlos Silva',
  email: 'carlos.silva@email.com',
  avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
  subscription: {
    plan: 'pro',
    status: 'active',
    expiresAt: '2024-12-31',
    features: [
      'Otimizações Avançadas',
      'Backup de Configurações',
      'Suporte Prioritário',
      'Relatórios Detalhados'
    ]
  },
  preferences: {
    theme: 'dark',
    language: 'pt-BR',
    notifications: true
  },
  stats: {
    joinedAt: '2023-03-15',
    optimizationsApplied: 127,
    lastLogin: '2024-08-09T10:30:00Z'
  }
};

// Simula delay de rede
const simulateDelay = (ms: number = 1000) => 
  new Promise(resolve => setTimeout(resolve, ms));

export class ProfileAPI {
  static async getProfile(): Promise<ApiResponse<UserProfile>> {
    await simulateDelay(800);
    
    return {
      success: true,
      data: { ...mockUserProfile },
      message: 'Perfil carregado com sucesso'
    };
  }

  static async updateProfile(updates: UpdateProfileRequest): Promise<ApiResponse<UserProfile>> {
    await simulateDelay(1200);
    
    // Simula validação
    if (updates.email && !updates.email.includes('@')) {
      return {
        success: false,
        message: 'Email inválido'
      };
    }

    if (updates.name && updates.name.trim().length < 2) {
      return {
        success: false,
        message: 'Nome deve ter pelo menos 2 caracteres'
      };
    }

    // Simula atualização
    const updatedProfile = {
      ...mockUserProfile,
      ...updates,
      preferences: {
        ...mockUserProfile.preferences,
        ...updates.preferences
      }
    };

    return {
      success: true,
      data: updatedProfile,
      message: 'Perfil atualizado com sucesso'
    };
  }

  static async uploadAvatar(file: File): Promise<ApiResponse<string>> {
    await simulateDelay(1500);
    
    // Simula validação de arquivo
    if (!file.type.startsWith('image/')) {
      return {
        success: false,
        message: 'Arquivo deve ser uma imagem'
      };
    }

    if (file.size > 5 * 1024 * 1024) { // 5MB
      return {
        success: false,
        message: 'Arquivo deve ter menos de 5MB'
      };
    }

    // Simula upload e retorna URL
    const mockUrl = `https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face&_t=${Date.now()}`;
    
    return {
      success: true,
      data: mockUrl,
      message: 'Avatar atualizado com sucesso'
    };
  }
}