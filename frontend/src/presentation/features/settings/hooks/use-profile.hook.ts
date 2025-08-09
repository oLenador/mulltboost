// features/profile/hooks/useProfile.ts
import { useState, useEffect } from 'react';
import { UpdateProfileRequest, UserProfile } from '../profile.types';
import { ProfileAPI } from '../profile-api.mock';


export const useProfile = () => {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadProfile = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await ProfileAPI.getProfile();
      if (response.success && response.data) {
        setProfile(response.data);
      } else {
        setError(response.message || 'Erro ao carregar perfil');
      }
    } catch (err) {
      setError('Erro de conexão');
    } finally {
      setIsLoading(false);
    }
  };

  const updateProfile = async (updates: UpdateProfileRequest) => {
    setIsUpdating(true);
    setError(null);
    
    try {
      const response = await ProfileAPI.updateProfile(updates);
      if (response.success && response.data) {
        setProfile(response.data);
        return { success: true, message: response.message };
      } else {
        setError(response.message || 'Erro ao atualizar perfil');
        return { success: false, message: response.message };
      }
    } catch (err) {
      const message = 'Erro de conexão';
      setError(message);
      return { success: false, message };
    } finally {
      setIsUpdating(false);
    }
  };

  const uploadAvatar = async (file: File) => {
    setIsUpdating(true);
    setError(null);
    
    try {
      const response = await ProfileAPI.uploadAvatar(file);
      if (response.success && response.data) {
        // Atualiza o perfil com o novo avatar
        if (profile) {
          setProfile({ ...profile, avatar: response.data });
        }
        return { success: true, message: response.message };
      } else {
        setError(response.message || 'Erro ao fazer upload do avatar');
        return { success: false, message: response.message };
      }
    } catch (err) {
      const message = 'Erro de conexão';
      setError(message);
      return { success: false, message };
    } finally {
      setIsUpdating(false);
    }
  };

  useEffect(() => {
    loadProfile();
  }, []);

  return {
    profile,
    isLoading,
    isUpdating,
    error,
    updateProfile,
    uploadAvatar,
    refetch: loadProfile
  };
};