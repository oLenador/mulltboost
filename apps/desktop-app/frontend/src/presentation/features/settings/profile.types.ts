// features/profile/types/profile.types.ts
export interface UserProfile {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    subscription: {
      plan: 'free' | 'pro' | 'premium';
      status: 'active' | 'expired' | 'cancelled';
      expiresAt?: string;
      features: string[];
    };
    preferences: {
      theme: 'light' | 'dark' | 'system';
      language: string;
      notifications: boolean;
    };
    stats: {
      joinedAt: string;
      optimizationsApplied: number;
      lastLogin: string;
    };
  }
  
  export interface UpdateProfileRequest {
    name?: string;
    email?: string;
    avatar?: string;
    preferences?: Partial<UserProfile['preferences']>;
  }
  
  export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    message?: string;
  }