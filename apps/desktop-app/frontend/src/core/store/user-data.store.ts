import { atomWithStorage } from 'jotai/utils';

export type UserData = {
  awl: string;
  name: string;
  email: string;
  language: string;
  avatar?: string;
  notifications: boolean;
  theme: 'light' | 'dark' | 'system';
  subscription: {
    plan: 'free' | 'pro' | 'premium';
    status: 'active' | 'inactive';
    expiresAt?: string;
    features: string[];
  };
  stats: {
    joinedAt: string;
    optimizationsApplied: number;
    lastLogin: string;
  };
};

export const userDataAtom = atomWithStorage<UserData>('user_data', {
  awl: '', // Identificador único do backend ou PC
  name: '',
  email: '',
  language: navigator.language || 'en',
  avatar: '',
  notifications: false,
  theme: 'dark',
  subscription: {
    plan: 'free',
    status: 'active',
    expiresAt: undefined,
    features: ['Acesso básico', 'Suporte por email']
  },
  stats: {
    joinedAt: new Date().toISOString(),
    optimizationsApplied: 0,
    lastLogin: new Date().toISOString()
  }
});