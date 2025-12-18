import { z } from 'zod';

// Profile update validation schema
export const updateProfileSchema = z.object({
    first_name: z.string().min(2, 'First name must be at least 2 characters').max(50, 'First name must be less than 50 characters'),
    last_name: z.string().min(2, 'Last name must be at least 2 characters').max(50, 'Last name must be less than 50 characters'),
    phone_number: z.string().optional(),
    address: z.string().optional(),
    city: z.string().optional(),
    country: z.string().optional(),
});

// Settings update validation schema
export const updateSettingsSchema = z.object({
    language: z.string().length(2, 'Language code must be exactly 2 characters'),
    timezone: z.string().min(1, 'Timezone is required'),
});

// Change password validation schema
export const changePasswordSchema = z.object({
    current_password: z.string().min(1, 'Current password is required'),
    new_password: z.string()
        .min(8, 'Password must be at least 8 characters')
        .max(128, 'Password must be less than 128 characters')
        .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
        .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
        .regex(/[0-9]/, 'Password must contain at least one number')
        .regex(/[^A-Za-z0-9]/, 'Password must contain at least one special character'),
});

// Toggle 2FA validation schema
export const toggleTwoFASchema = z.object({
    enable: z.boolean(),
    code: z.string().length(6, 'Code must be exactly 6 digits').regex(/^\d+$/, 'Code must contain only digits'),
});

export type UpdateProfileData = z.infer<typeof updateProfileSchema>;
export type UpdateSettingsData = z.infer<typeof updateSettingsSchema>;
export type ChangePasswordData = z.infer<typeof changePasswordSchema>;
export type ToggleTwoFAData = z.infer<typeof toggleTwoFASchema>;
