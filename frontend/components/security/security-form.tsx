import * as React from 'react';
import { useState, useEffect } from 'react';
import { AuthService, User, ChangePasswordRequest } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { AlertCircle, CheckCircle2, Lock, Smartphone } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { PasswordStrength } from '@/components/ui/password-strength';

export function SecurityForm() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadUser = async () => {
            try {
                const currentUser = await AuthService.getCurrentUser();
                setUser(currentUser);
            } catch {
                console.error('Failed to load user');
            } finally {
                setLoading(false);
            }
        };
        loadUser();
    }, []);

    if (loading) {
        return <div className="flex justify-center p-8">Loading security settings...</div>;
    }

    return (
        <div className="space-y-6 max-w-2xl mx-auto">
            <ChangePasswordSection />
            <TwoFactorSection user={user} />
        </div>
    );
}

function ChangePasswordSection() {
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
    const [formData, setFormData] = useState<ChangePasswordRequest>({
        current_password: '',
        new_password: '',
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setMessage(null);

        try {
            await AuthService.changePassword(formData);
            setMessage({ type: 'success', text: 'Password changed successfully' });
            setFormData({ current_password: '', new_password: '' });
        } catch (error: unknown) {
            setMessage({ type: 'error', text: (error as Error).message || 'Failed to change password' });
        } finally {
            setSaving(false);
        }
    };

    return (
        <Card>
            <CardHeader>
                <div className="flex items-center gap-2">
                    <Lock className="h-6 w-6 text-primary" />
                    <CardTitle>Change Password</CardTitle>
                </div>
                <CardDescription>Ensure your account is using a long, random password to stay secure.</CardDescription>
            </CardHeader>
            <CardContent>
                {message && (
                    <Alert variant={message.type === 'error' ? 'destructive' : 'default'} className="mb-6 bg-primary/5 border-primary/20 text-primary">
                        {message.type === 'success' ? <CheckCircle2 className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
                        <AlertTitle>{message.type === 'success' ? 'Success' : 'Error'}</AlertTitle>
                        <AlertDescription>{message.text}</AlertDescription>
                    </Alert>
                )}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="current_password">Current Password</Label>
                        <Input
                            id="current_password"
                            name="current_password"
                            type="password"
                            value={formData.current_password}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="new_password">New Password</Label>
                        <Input
                            id="new_password"
                            name="new_password"
                            type="password"
                            value={formData.new_password}
                            onChange={handleChange}
                            required
                            minLength={8}
                        />
                        <PasswordStrength password={formData.new_password} />
                    </div>
                    <Button type="submit" disabled={saving}>
                        {saving ? 'Updating...' : 'Update Password'}
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}

function TwoFactorSection({ user }: { user: User | null }) {
    const [isEnabled, setIsEnabled] = useState(user?.twofa_enabled || false);
    const [verifying, setVerifying] = useState(false);
    const [otpCode, setOtpCode] = useState('');
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    // If we are currently in the process of toggling (waiting for OTP)
    const [confirmingToggle, setConfirmingToggle] = useState(false);
    const [targetState, setTargetState] = useState(false); // TRUE = enabling, FALSE = disabling

    const handleToggle = async (checked: boolean) => {
        setMessage(null);
        setTargetState(checked);
        setConfirmingToggle(true);

        // In a real app, if enabling, we would show a QR code step here.
        // For now, we simulate requesting an OTP (e.g., via email if not using authenticator app yet, 
        // or just assume user has setup). 
        // But typically '2fa' type OTP implies Authenticator App code.
        // Since we don't have QR code generation yet, let's assume 'email-verification' style for the "Code"
        // OR just ask for the code if they have it.

        // If enabling: "Scan QR (not implemented), then enter code."
        // If disabling: "Enter code from app to confirm."

        if (checked) {
            setMessage({ type: 'success', text: 'Please enter a code to confirm enabling 2FA.' });
        } else {
            setMessage({ type: 'success', text: 'Please enter a code to confirm disabling 2FA.' });
        }
    };

    const verifyAndToggle = async () => {
        setVerifying(true);
        setMessage(null);
        try {
            await AuthService.toggleTwoFA({ enable: targetState, code: otpCode });
            setIsEnabled(targetState);
            setConfirmingToggle(false);
            setOtpCode('');
            setMessage({ type: 'success', text: `Two-factor authentication ${targetState ? 'enabled' : 'disabled'} successfully.` });
        } catch (error: unknown) {
            setMessage({ type: 'error', text: (error as Error).message || 'Verification failed.' });
        } finally {
            setVerifying(false);
        }
    };

    const cancelToggle = () => {
        setConfirmingToggle(false);
        setOtpCode('');
        setMessage(null);
    }

    return (
        <Card>
            <CardHeader>
                <div className="flex items-center gap-2">
                    <Smartphone className="h-6 w-6 text-primary" />
                    <CardTitle>Two-Factor Authentication (2FA)</CardTitle>
                </div>
                <CardDescription>Add an extra layer of security to your account.</CardDescription>
            </CardHeader>
            <CardContent>
                {message && (
                    <Alert variant={message.type === 'error' ? 'destructive' : 'default'} className="mb-6 bg-primary/5 border-primary/20 text-primary">
                        {message.type === 'success' ? <CheckCircle2 className="h-4 w-4" /> : <AlertCircle className="h-4 w-4" />}
                        <AlertTitle>{message.type === 'success' ? 'Info' : 'Error'}</AlertTitle>
                        <AlertDescription>{message.text}</AlertDescription>
                    </Alert>
                )}

                <div className="flex items-center justify-between space-x-2 mb-4">
                    <div className="flex flex-col space-y-1">
                        <Label htmlFor="2fa-mode" className="font-medium">Enable 2FA</Label>
                        <span className="text-sm text-muted-foreground">
                            {isEnabled ? 'Your account is secured with 2FA.' : 'Protect your account with 2FA.'}
                        </span>
                    </div>
                    <Switch id="2fa-mode" checked={isEnabled} onCheckedChange={handleToggle} disabled={confirmingToggle} />
                </div>

                {confirmingToggle && (
                    <div className="mt-4 p-4 border rounded-md bg-secondary/10 space-y-4">
                        <p className="text-sm">
                            {targetState
                                ? "To enable 2FA, please enter the code from your authenticator app (mock: use a valid OTP)."
                                : "To disable 2FA, please confirm with a code from your authenticator app."}
                        </p>
                        <div className="flex gap-2">
                            <Input
                                placeholder="000000"
                                value={otpCode}
                                onChange={(e) => setOtpCode(e.target.value)}
                                maxLength={6}
                                className="w-32 tracking-widest text-center"
                            />
                            <Button onClick={verifyAndToggle} disabled={verifying || otpCode.length !== 6}>
                                {verifying ? 'Verifying...' : 'Confirm'}
                            </Button>
                            <Button variant="ghost" onClick={cancelToggle} disabled={verifying}>
                                Cancel
                            </Button>
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
