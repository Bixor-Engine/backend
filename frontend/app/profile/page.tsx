'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { AuthService, User } from '@/lib/auth';
import { ProtectedNavbar } from '@/components/protected-navbar';
import { ProfileForm } from '@/components/profile/profile-form';

export default function ProfilePage() {
    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const checkAuth = async () => {
            if (!AuthService.isAuthenticated()) {
                router.push('/auth/signin');
                return;
            }

            try {
                const currentUser = await AuthService.getCurrentUser();
                if (!currentUser) {
                    router.push('/auth/signin');
                    return;
                }
                setUser(currentUser);
            } catch (error) {
                console.error(error);
                router.push('/auth/signin');
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, [router]);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-background">
                <div className="animate-pulse text-primary">Loading...</div>
            </div>
        );
    }

    if (!user) return null;

    return (
        <div className="min-h-screen bg-background">
            <ProtectedNavbar user={user} currentPage="profile" />
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-foreground">Profile Settings</h1>
                    <p className="text-muted-foreground mt-2">Manage your account information</p>
                </div>
                <ProfileForm />
            </main>
        </div>
    );
}
