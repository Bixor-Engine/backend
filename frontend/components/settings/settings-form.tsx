'use client';

import * as React from 'react';
import { useState, useEffect } from 'react';
import { AuthService, UpdateSettingsRequest } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { AlertCircle, CheckCircle2, Settings as SettingsIcon } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

export function SettingsForm() {
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const [formData, setFormData] = useState<UpdateSettingsRequest>({
        language: 'en',
        timezone: 'UTC',
    });

    useEffect(() => {
        const loadUser = async () => {
            try {
                const currentUser = await AuthService.getCurrentUser();
                if (currentUser) {
                    setFormData({
                        language: currentUser.language || 'en',
                        timezone: currentUser.timezone || 'UTC',
                    });
                }
            } catch {
                setMessage({ type: 'error', text: 'Failed to load settings' });
            } finally {
                setLoading(false);
            }
        };

        loadUser();
    }, []);

    const handleSelectChange = (name: keyof UpdateSettingsRequest, value: string) => {
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setMessage(null);

        try {
            await AuthService.updateSettings(formData);
            setMessage({ type: 'success', text: 'Settings updated successfully' });
        } catch (error: unknown) {
            setMessage({ type: 'error', text: (error as Error).message || 'Failed to update settings' });
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return <div className="flex justify-center p-8">Loading settings...</div>;
    }

    return (
        <Card className="w-full max-w-2xl mx-auto">
            <CardHeader>
                <div className="flex items-center gap-2">
                    <SettingsIcon className="h-6 w-6 text-primary" />
                    <CardTitle>Preferences</CardTitle>
                </div>
                <CardDescription>
                    Customize your experience by changing your language and timezone.
                </CardDescription>
            </CardHeader>
            <CardContent>
                {message && (
                    <Alert variant={message.type === 'error' ? 'destructive' : 'default'} className="mb-6 bg-primary/5 border-primary/20 text-primary">
                        {message.type === 'success' ? <CheckCircle2 className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
                        <AlertTitle>{message.type === 'success' ? 'Success' : 'Error'}</AlertTitle>
                        <AlertDescription>{message.text}</AlertDescription>
                    </Alert>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="space-y-2">
                        <Label htmlFor="language">Language</Label>
                        <Select
                            value={formData.language}
                            onValueChange={(value) => handleSelectChange('language', value)}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select Language" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="en">English (US)</SelectItem>
                                <SelectItem value="es">Español</SelectItem>
                                <SelectItem value="fr">Français</SelectItem>
                                <SelectItem value="de">Deutsch</SelectItem>
                                <SelectItem value="zh">中文</SelectItem>
                                <SelectItem value="ja">日本語</SelectItem>
                            </SelectContent>
                        </Select>
                        <p className="text-xs text-muted-foreground">Select your preferred language interface.</p>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="timezone">Timezone</Label>
                        <Select
                            value={formData.timezone}
                            onValueChange={(value) => handleSelectChange('timezone', value)}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select Timezone" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="UTC">UTC (Universal Coordinated Time)</SelectItem>
                                <SelectItem value="America/New_York">Eastern Time (US & Canada)</SelectItem>
                                <SelectItem value="America/Los_Angeles">Pacific Time (US & Canada)</SelectItem>
                                <SelectItem value="Europe/London">London</SelectItem>
                                <SelectItem value="Europe/Paris">Paris</SelectItem>
                                <SelectItem value="Asia/Tokyo">Tokyo</SelectItem>
                                <SelectItem value="Asia/Singapore">Singapore</SelectItem>
                                {/* Add more timezones as needed */}
                            </SelectContent>
                        </Select>
                        <p className="text-xs text-muted-foreground">Your timezone is used for notifications and transaction logs.</p>
                    </div>

                    <Button type="submit" className="w-full" disabled={saving}>
                        {saving ? 'Saving...' : 'Save Preferences'}
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}
