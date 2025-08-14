import React from 'react';
import { User, Camera, Loader2 } from 'lucide-react';

interface Profile {
  name: string;
  email: string;
  avatar?: string;
}

interface ProfileAvatarProps {
  profile: Profile;
  isEditing: boolean;
  isUpdating: boolean;
  fileInputRef: React.RefObject<HTMLInputElement>;
}

export const ProfileAvatar: React.FC<ProfileAvatarProps> = ({
  profile,
  isEditing,
  isUpdating,
  fileInputRef,
}) => {
  return (
    <div className="flex items-center space-x-4">
      <div className="relative">
        <div className="w-20 h-20 rounded-full bg-zinc-800 overflow-hidden">
          {profile.avatar ? (
            <img
              src={profile.avatar}
              alt="Avatar"
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <User className="w-10 h-10 text-zinc-500" />
            </div>
          )}
        </div>
        {isEditing && (
          <button
            onClick={() => fileInputRef.current?.click()}
            disabled={isUpdating}
            className="absolute -bottom-1 -right-1 w-6 h-6 bg-zinc-700 rounded-full flex items-center justify-center hover:bg-zinc-600 transition-colors"
          >
            {isUpdating ? (
              <Loader2 className="w-3 h-3 animate-spin" />
            ) : (
              <Camera className="w-3 h-3" />
            )}
          </button>
        )}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          className="hidden"
        />
      </div>
      <div>
        <h3 className="text-lg font-medium text-zinc-200">{profile.name}</h3>
        <p className="text-zinc-400 text-sm">{profile.email}</p>
      </div>
    </div>
  );
};