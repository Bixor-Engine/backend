'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';
import { Check, X } from 'lucide-react';

interface PasswordStrengthProps {
    password: string;
    className?: string;
}

interface PasswordRequirement {
    label: string;
    test: (password: string) => boolean;
}

const requirements: PasswordRequirement[] = [
    {
        label: 'At least 8 characters',
        test: (password) => password.length >= 8,
    },
    {
        label: 'Contains uppercase letter',
        test: (password) => /[A-Z]/.test(password),
    },
    {
        label: 'Contains lowercase letter',
        test: (password) => /[a-z]/.test(password),
    },
    {
        label: 'Contains number',
        test: (password) => /[0-9]/.test(password),
    },
    {
        label: 'Contains special character',
        test: (password) => /[^A-Za-z0-9]/.test(password),
    },
];

export function PasswordStrength({ password, className }: PasswordStrengthProps) {
    const passedRequirements = requirements.filter((req) => req.test(password));
    const strength = (passedRequirements.length / requirements.length) * 100;

    const getStrengthColor = () => {
        if (strength === 0) return 'bg-gray-200';
        if (strength < 40) return 'bg-red-500';
        if (strength < 80) return 'bg-yellow-500';
        return 'bg-green-500';
    };

    const getStrengthText = () => {
        if (strength === 0) return '';
        if (strength < 40) return 'Weak';
        if (strength < 80) return 'Medium';
        return 'Strong';
    };

    return (
        <div className={cn('space-y-2', className)}>
            {password && (
                <>
                    <div className="space-y-1">
                        <div className="flex justify-between items-center text-xs">
                            <span className="text-muted-foreground">Password strength</span>
                            <span className={cn('font-medium',
                                strength < 40 ? 'text-red-500' :
                                    strength < 80 ? 'text-yellow-500' :
                                        'text-green-500'
                            )}>
                                {getStrengthText()}
                            </span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-1.5">
                            <div
                                className={cn('h-1.5 rounded-full transition-all duration-300', getStrengthColor())}
                                style={{ width: `${strength}%` }}
                            />
                        </div>
                    </div>
                    <ul className="space-y-1">
                        {requirements.map((req, index) => {
                            const passed = req.test(password);
                            return (
                                <li
                                    key={index}
                                    className={cn(
                                        'flex items-center gap-2 text-xs',
                                        passed ? 'text-green-600' : 'text-muted-foreground'
                                    )}
                                >
                                    {passed ? (
                                        <Check className="h-3 w-3" />
                                    ) : (
                                        <X className="h-3 w-3 text-gray-400" />
                                    )}
                                    {req.label}
                                </li>
                            );
                        })}
                    </ul>
                </>
            )}
        </div>
    );
}
