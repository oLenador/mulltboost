import { UserData, userDataAtom } from '@/core/store/user-data.store';
import { useAtom } from 'jotai';

export const useUserData = () => {
  const [userData, setUserData] = useAtom(userDataAtom);

  const updateUserData = (updates: Partial<UserData>) => {
    setUserData(prev => ({ ...prev, ...updates }));
  };

  const updatePreferences = (preferences: Partial<Pick<UserData, 'notifications' | 'theme' | 'language'>>) => {
    setUserData(prev => ({ ...prev, ...preferences }));
  };

  const updateStats = (stats: Partial<UserData['stats']>) => {
    setUserData(prev => ({
      ...prev,
      stats: { ...prev.stats, ...stats }
    }));
  };

  return {
    userData,
    setUserData,
    updateUserData,
    updatePreferences,
    updateStats
  };
}