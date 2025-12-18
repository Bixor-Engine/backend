'use client';

import { useEffect } from 'react';
import { migrateTokenStorage } from '@/lib/migrate-storage';

/**
 * TokenMigration Component
 * 
 * Runs once on app initialization to migrate users from old localStorage
 * token storage to the new secure in-memory + sessionStorage approach.
 * 
 * This component should be placed in the root layout to ensure it runs
 * on every page load (migration is idempotent and safe to run multiple times).
 */
export function TokenMigration() {
    useEffect(() => {
        // Run migration on client-side only, once per session
        migrateTokenStorage();
    }, []);

    // This component doesn't render anything
    return null;
}
