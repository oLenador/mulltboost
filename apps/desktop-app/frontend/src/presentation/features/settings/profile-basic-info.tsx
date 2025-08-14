import React from 'react';
import { User, Mail, Settings, Check, X, Loader2 } from 'lucide-react';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import { Label } from '@/presentation/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card';
import { UserData } from '@/core/store/user-data.store';
import { ProfileAvatar } from './profile-avatar';


interface ProfileBasicInfoProps {
  profile: UserData;
  isEditing: boolean;
  isUpdating: boolean;
  formData: {
    name: string;
    email: string;
    notifications: boolean;
    theme: string;
    language: string;
  };
  setFormData: React.Dispatch<React.SetStateAction<any>>;
  fileInputRef: React.RefObject<HTMLInputElement>;
  onEdit: () => void;
  onSave: () => void;
  onCancel: () => void;
}

export const ProfileBasicInfo: React.FC<ProfileBasicInfoProps> = ({
  profile,
  isEditing,
  isUpdating,
  formData,
  setFormData,
  fileInputRef,
  onEdit,
  onSave,
  onCancel,
}) => {
  return (
    <Card variant="zinc">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-zinc-100">Informações Básicas</CardTitle>
            <CardDescription className="text-zinc-400">
              Suas informações pessoais
            </CardDescription>
          </div>
          {!isEditing && (
            <Button variant="zinc" size="sm" onClick={onEdit}>
              <Settings className="w-4 h-4 mr-2" />
              Editar
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Avatar Section */}
        <ProfileAvatar
          profile={profile}
          isEditing={isEditing}
          isUpdating={isUpdating}
          fileInputRef={fileInputRef}
        />

        {/* Form Fields */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-zinc-300">Nome</Label>
            {isEditing ? (
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="bg-zinc-800 border-zinc-700 text-zinc-200"
              />
            ) : (
              <div className="flex items-center space-x-2 p-3 bg-zinc-800 rounded-lg">
                <User className="w-4 h-4 text-zinc-500" />
                <span className="text-zinc-300">{profile.name}</span>
              </div>
            )}
          </div>
          
          <div className="space-y-2">
            <Label className="text-zinc-300">Email</Label>
            {isEditing ? (
              <Input
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="bg-zinc-800 border-zinc-700 text-zinc-200"
              />
            ) : (
              <div className="flex items-center space-x-2 p-3 bg-zinc-800 rounded-lg">
                <Mail className="w-4 h-4 text-zinc-500" />
                <span className="text-zinc-300">{profile.email}</span>
              </div>
            )}
          </div>
        </div>

        {/* Action Buttons */}
        {isEditing && (
          <div className="flex items-center space-x-3 pt-4 border-t border-zinc-800">
            <Button
              onClick={onSave}
              disabled={isUpdating}
              className="bg-zinc-700 hover:bg-zinc-600"
            >
              {isUpdating ? (
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              ) : (
                <Check className="w-4 h-4 mr-2" />
              )}
              Salvar Alterações
            </Button>
            <Button variant="zinc" onClick={onCancel} disabled={isUpdating}>
              <X className="w-4 h-4 mr-2" />
              Cancelar
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
};