import React from 'react';
import { Label } from '@/presentation/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card';
import { Badge } from '@/presentation/components/ui/badge';
import { Switch } from '@/presentation/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/presentation/components/ui/select';
import { UserData } from '@/core/store/user-data.store';


interface ProfilePreferencesProps {
    user: UserData;
  setUserData: React.Dispatch<React.SetStateAction<any>>;
}

export const ProfilePreferences: React.FC<ProfilePreferencesProps> = ({
  user,
  setUserData,
}) => {
  return (
    <Card variant="zinc">
      <CardHeader>
        <CardTitle className="text-zinc-100">Preferências</CardTitle>
        <CardDescription className="text-zinc-400">
          Configure suas preferências do sistema
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-zinc-300">Idioma</Label>
              <Select 
                value={user.language} 
                onValueChange={(value) => setUserData({ ...user, language: value })}
              >
                <SelectTrigger className="bg-zinc-800 border-zinc-700">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pt">Português</SelectItem>
                  <SelectItem value="pt-BR">Português (BR)</SelectItem>
                  <SelectItem value="en">English (US)</SelectItem>
                  <SelectItem value="es">Español</SelectItem>
                  <SelectItem value="ru">Russian</SelectItem>
                </SelectContent>
              </Select>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
