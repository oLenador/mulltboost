// features/profile/components/ProfilePage.tsx
import React, { useState, useRef } from 'react';
import {
    User,
    Mail,
    Crown,
    Calendar,
    Activity,
    Settings,
    Shield,
    Camera,
    Check,
    X,
    Loader2,
    AlertCircle,
    Star,
    Clock,
    Target
} from 'lucide-react';
import { Button } from '@/presentation/components/ui/button';
import { Input } from '@/presentation/components/ui/input';
import { Label } from '@/presentation/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card';
import { Badge } from '@/presentation/components/ui/badge';
import { Switch } from '@/presentation/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/presentation/components/ui/select';
import { useProfile } from './hooks/use-profile.hook';
import BasePage from '@/presentation/components/pages/base-page';

const ProfilePage: React.FC = () => {
    const { profile, isLoading, isUpdating, error, updateProfile, uploadAvatar } = useProfile();
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        notifications: false,
        theme: 'dark',
        language: 'pt-BR'
    });
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    // Atualiza form quando profile carrega
    React.useEffect(() => {
        if (profile && !isEditing) {
            setFormData({
                name: profile.name,
                email: profile.email,
                notifications: profile.preferences.notifications,
                theme: profile.preferences.theme,
                language: profile.preferences.language
            });
        }
    }, [profile, isEditing]);

    const handleEdit = () => {
        setIsEditing(true);
        setMessage(null);
    };

    const handleCancel = () => {
        if (profile) {
            setFormData({
                name: profile.name,
                email: profile.email,
                notifications: profile.preferences.notifications,
                theme: profile.preferences.theme,
                language: profile.preferences.language
            });
        }
        setIsEditing(false);
        setMessage(null);
    };

    const handleSave = async () => {
        const result = await updateProfile({
            name: formData.name,
            email: formData.email,
            preferences: {
                notifications: formData.notifications,
                theme: formData.theme as 'light' | 'dark' | 'system',
                language: formData.language
            }
        });

        setMessage({
            type: result.success ? 'success' : 'error',
            text: result.message || (result.success ? 'Perfil atualizado!' : 'Erro ao atualizar')
        });

        if (result.success) {
            setIsEditing(false);
        }
    };

    const handleAvatarChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        const result = await uploadAvatar(file);
        setMessage({
            type: result.success ? 'success' : 'error',
            text: result.message || (result.success ? 'Avatar atualizado!' : 'Erro ao atualizar avatar')
        });
    };

    const getSubscriptionIcon = () => {
        if (!profile) return Crown;
        switch (profile.subscription.plan) {
            case 'premium': return Crown;
            case 'pro': return Star;
            default: return Shield;
        }
    };

    const getSubscriptionColor = () => {
        if (!profile) return 'text-zinc-400';
        switch (profile.subscription.plan) {
            case 'premium': return 'text-yellow-400';
            case 'pro': return 'text-blue-400';
            default: return 'text-zinc-400';
        }
    };

    const getPlanName = () => {
        if (!profile) return 'Free';
        switch (profile.subscription.plan) {
            case 'premium': return 'Premium';
            case 'pro': return 'Pro';
            default: return 'Free';
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('pt-BR');
    };

    if (isLoading) {
        return (
            <BasePage>
                <div className="max-w-4xl mx-auto">
                    <div className="flex items-center justify-center py-12">
                        <Loader2 className="w-8 h-8 animate-spin text-zinc-400" />
                        <span className="ml-2 text-zinc-400">Carregando perfil...</span>
                    </div>
                </div>
            </BasePage>
        );
    }

    if (error || !profile) {
        return (
            <BasePage>
                <div className="max-w-4xl mx-auto">
                    <Card variant="zinc" className="text-center py-12">
                        <CardContent>
                            <AlertCircle className="w-12 h-12 text-red-400 mx-auto mb-4" />
                            <h2 className="text-xl font-semibold text-zinc-200 mb-2">
                                Erro ao carregar perfil
                            </h2>
                            <p className="text-zinc-400 mb-4">{error || 'Perfil não encontrado'}</p>
                            <Button variant="zinc" onClick={() => window.location.reload()}>
                                Tentar Novamente
                            </Button>
                        </CardContent>
                    </Card>
                </div>
            </BasePage>);
    }

    const SubscriptionIcon = getSubscriptionIcon();

    return (
        <BasePage>
            <div className="max-w-4xl  space-y-8">
                {/* Header */}
                <div>
                    <h1 className="text-3xl font-bold text-zinc-100 mb-2">Meu Perfil</h1>
                    <p className="text-zinc-400">Gerencie suas informações pessoais e preferências</p>
                </div>

                {/* Message */}
                {message && (
                    <div className={`p-4 rounded-lg border ${message.type === 'success'
                        ? 'bg-green-400/10 border-green-400/20 text-green-400'
                        : 'bg-red-400/10 border-red-400/20 text-red-400'
                        }`}>
                        <div className="flex items-center space-x-2">
                            {message.type === 'success' ? (
                                <Check className="w-4 h-4" />
                            ) : (
                                <X className="w-4 h-4" />
                            )}
                            <span className="text-sm">{message.text}</span>
                        </div>
                    </div>
                )}

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Profile Info */}
                    <div className="lg:col-span-2 space-y-6">
                        {/* Basic Info Card */}
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
                                        <Button variant="zinc" size="sm" onClick={handleEdit}>
                                            <Settings className="w-4 h-4 mr-2" />
                                            Editar
                                        </Button>
                                    )}
                                </div>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                {/* Avatar Section */}
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
                                            onChange={handleAvatarChange}
                                            className="hidden"
                                        />
                                    </div>
                                    <div>
                                        <h3 className="text-lg font-medium text-zinc-200">{profile.name}</h3>
                                        <p className="text-zinc-400 text-sm">{profile.email}</p>
                                    </div>
                                </div>

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
                                            onClick={handleSave}
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
                                        <Button
                                            variant="zinc"
                                            onClick={handleCancel}
                                            disabled={isUpdating}
                                        >
                                            <X className="w-4 h-4 mr-2" />
                                            Cancelar
                                        </Button>
                                    </div>
                                )}
                            </CardContent>
                        </Card>

                        {/* Preferences Card */}
                        <Card variant="zinc">
                            <CardHeader>
                                <CardTitle className="text-zinc-100">Preferências</CardTitle>
                                <CardDescription className="text-zinc-400">
                                    Configure suas preferências do sistema
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <div>
                                        <Label className="text-zinc-300">Notificações</Label>
                                        <p className="text-xs text-zinc-500">Receber notificações do sistema</p>
                                    </div>
                                    {isEditing ? (
                                        <Switch
                                            checked={formData.notifications}
                                            onCheckedChange={(checked) => setFormData({ ...formData, notifications: checked })}
                                        />
                                    ) : (
                                        <Badge variant={profile.preferences.notifications ? "default" : "secondary"}>
                                            {profile.preferences.notifications ? "Ativadas" : "Desativadas"}
                                        </Badge>
                                    )}
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <Label className="text-zinc-300">Tema</Label>
                                        {isEditing ? (
                                            <Select value={formData.theme} onValueChange={(value) => setFormData({ ...formData, theme: value })}>
                                                <SelectTrigger className="bg-zinc-800 border-zinc-700">
                                                    <SelectValue />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="light">Claro</SelectItem>
                                                    <SelectItem value="dark">Escuro</SelectItem>
                                                    <SelectItem value="system">Sistema</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        ) : (
                                            <div className="p-3 bg-zinc-800 rounded-lg">
                                                <span className="text-zinc-300 capitalize">{profile.preferences.theme}</span>
                                            </div>
                                        )}
                                    </div>

                                    <div className="space-y-2">
                                        <Label className="text-zinc-300">Idioma</Label>
                                        {isEditing ? (
                                            <Select value={formData.language} onValueChange={(value) => setFormData({ ...formData, language: value })}>
                                                <SelectTrigger className="bg-zinc-800 border-zinc-700">
                                                    <SelectValue />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="pt-BR">Português (BR)</SelectItem>
                                                    <SelectItem value="en-US">English (US)</SelectItem>
                                                    <SelectItem value="es-ES">Español</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        ) : (
                                            <div className="p-3 bg-zinc-800 rounded-lg">
                                                <span className="text-zinc-300">{profile.preferences.language}</span>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </div>

                    {/* Sidebar */}
                    <div className="space-y-6">
                        {/* Subscription Status */}
                        <Card variant="zinc">
                            <CardHeader>
                                <CardTitle className="flex items-center space-x-2 text-zinc-100">
                                    <SubscriptionIcon className={`w-5 h-5 ${getSubscriptionColor()}`} />
                                    <span>Assinatura</span>
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div>
                                    <div className="flex items-center justify-between mb-2">
                                        <span className="text-zinc-300">Plano</span>
                                        <Badge className={getSubscriptionColor()}>
                                            {getPlanName()}
                                        </Badge>
                                    </div>

                                    <div className="flex items-center justify-between mb-2">
                                        <span className="text-zinc-300">Status</span>
                                        <Badge variant={profile.subscription.status === 'active' ? 'default' : 'destructive'}>
                                            {profile.subscription.status === 'active' ? 'Ativo' : 'Inativo'}
                                        </Badge>
                                    </div>

                                    {profile.subscription.expiresAt && (
                                        <div className="flex items-center justify-between">
                                            <span className="text-zinc-300">Expira em</span>
                                            <span className="text-zinc-400 text-sm">
                                                {formatDate(profile.subscription.expiresAt)}
                                            </span>
                                        </div>
                                    )}
                                </div>

                                <div className="pt-4 border-t border-zinc-800">
                                    <h4 className="text-sm font-medium text-zinc-300 mb-2">Recursos inclusos:</h4>
                                    <ul className="space-y-1">
                                        {profile.subscription.features.map((feature, index) => (
                                            <li key={index} className="flex items-center space-x-2 text-xs">
                                                <Check className="w-3 h-3 text-green-400" />
                                                <span className="text-zinc-400">{feature}</span>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            </CardContent>
                        </Card>

                        {/* Statistics */}
                        <Card variant="zinc">
                            <CardHeader>
                                <CardTitle className="text-zinc-100">Estatísticas</CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="flex items-center space-x-3">
                                    <div className="p-2 bg-zinc-800 rounded-lg">
                                        <Calendar className="w-4 h-4 text-zinc-400" />
                                    </div>
                                    <div>
                                        <p className="text-xs text-zinc-500">Membro desde</p>
                                        <p className="text-sm text-zinc-300">{formatDate(profile.stats.joinedAt)}</p>
                                    </div>
                                </div>

                                <div className="flex items-center space-x-3">
                                    <div className="p-2 bg-zinc-800 rounded-lg">
                                        <Target className="w-4 h-4 text-green-400" />
                                    </div>
                                    <div>
                                        <p className="text-xs text-zinc-500">Otimizações aplicadas</p>
                                        <p className="text-sm text-zinc-300">{profile.stats.optimizationsApplied}</p>
                                    </div>
                                </div>

                                <div className="flex items-center space-x-3">
                                    <div className="p-2 bg-zinc-800 rounded-lg">
                                        <Clock className="w-4 h-4 text-blue-400" />
                                    </div>
                                    <div>
                                        <p className="text-xs text-zinc-500">Último acesso</p>
                                        <p className="text-sm text-zinc-300">
                                            {formatDate(profile.stats.lastLogin)}
                                        </p>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>
        </BasePage>
    );
};

export default ProfilePage;