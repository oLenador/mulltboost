import React, { useState, useRef } from 'react';

import BasePage from '@/presentation/components/pages/base-page';
import { ProfileHeader } from './profile-header.component';
import { MessageAlert } from './profile-message-alert';
import { ProfilePreferences } from './profile-preferences-card';
import { ProfileSubscription } from './profile-subscription-card';
import { ProfileStatistics } from './profile-statistics';
import { ProfileBasicInfo } from './profile-basic-info';
import { useAtom } from 'jotai';
import { userDataAtom } from '@/core/store/user-data.store';


const ProfilePage: React.FC = () => {
    const [userData, setUserData] = useAtom(userDataAtom);
    const [isEditing, setIsEditing] = useState(false);
    const [isUpdating, setIsUpdating] = useState(false);
    const [formData, setFormData] = useState({
      name: userData.name,
      email: userData.email,
      notifications: userData.notifications,
      theme: userData.theme,
      language: userData.language
    });
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);
  
    // Atualiza form quando userData muda e não está editando
    React.useEffect(() => {
      if (!isEditing) {
        setFormData({
          name: userData.name,
          email: userData.email,
          notifications: userData.notifications,
          theme: userData.theme,
          language: userData.language
        });
      }
    }, [userData, isEditing]);
  
    const handleEdit = () => {
      setIsEditing(true);
      setMessage(null);
    };
  
    const handleCancel = () => {
      setFormData({
        name: userData.name,
        email: userData.email,
        notifications: userData.notifications,
        theme: userData.theme,
        language: userData.language
      });
      setIsEditing(false);
      setMessage(null);
    };
  
    const handleSave = async () => {
      setIsUpdating(true);
      
      try {
        // Simula um delay de API
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Atualiza o atom com os novos dados
        setUserData(prev => ({
          ...prev,
          name: formData.name,
          email: formData.email,
          notifications: formData.notifications,
          theme: formData.theme,
          language: formData.language,
          stats: {
            ...prev.stats,
            lastLogin: new Date().toISOString()
          }
        }));
  
        setMessage({
          type: 'success',
          text: 'Perfil atualizado com sucesso!'
        });
  
        setIsEditing(false);
      } catch (error) {
        setMessage({
          type: 'error',
          text: 'Erro ao atualizar perfil. Tente novamente.'
        });
      } finally {
        setIsUpdating(false);
      }
    };
  
    React.useEffect(() => {
      if (message) {
        const timer = setTimeout(() => setMessage(null), 5000);
        return () => clearTimeout(timer);
      }
    }, [message]);
  
    return (
      <BasePage>
        <div className="max-w-4xl space-y-8">
          <ProfileHeader />
          
          <MessageAlert message={message} />
  
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              <ProfileBasicInfo
                profile={userData}
                isEditing={isEditing}
                isUpdating={isUpdating}
                formData={formData}
                setFormData={setFormData}
                fileInputRef={fileInputRef}
                onEdit={handleEdit}
                onSave={handleSave}
                onCancel={handleCancel}
              />
              
              <ProfilePreferences
                user={userData}
                setUserData={setUserData}
              />
            </div>
  
            {/* Sidebar */}
            <div className="space-y-6">
              <ProfileSubscription profile={userData} />
              <ProfileStatistics profile={userData} />
            </div>
          </div>
        </div>
      </BasePage>
    );
  };
  
  export default ProfilePage;