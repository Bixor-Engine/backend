import { useState, useEffect, useCallback } from 'react';
import { AuthService } from '@/lib/auth';

export interface Wallet {
    id: string | null;
    // user_id removed as it's not returned for virtual wallets
    coin_id: number;
    balance: string;
    frozen_balance: string;
    coin_name: string;
    coin_ticker: string;
    coin_logo?: string;
    coin_decimal: number;
    usd_value?: string;
}

export interface Transaction {
    id: string;
    type: 'deposit' | 'withdraw' | 'transfer';
    amount: string;
    fee: string;
    description?: string;
    status: 'pending' | 'completed' | 'failed' | 'cancelled' | 'processing';
    created_at: string;
    coin_ticker: string;
    coin_name: string;
}

interface TransactionsResponse {
    data: Transaction[];
    total: number;
    page: number;
    limit: number;
}

export function useWallet() {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchWallets = useCallback(async () => {
        try {
            const response = await AuthService.fetch('/wallets');
            if (response.ok) {
                const data = await response.json();
                setWallets(data || []);
            } else {
                throw new Error('Failed to fetch wallets');
            }
        } catch (_err: unknown) {
            console.error('Error fetching wallets:', _err);
            // Don't set global error for this, just log it
        }
    }, []);

    const fetchTransactions = useCallback(async (page = 1, limit = 10) => {
        try {
            const response = await AuthService.fetch(`/transactions?page=${page}&limit=${limit}`);
            if (response.ok) {
                const data: TransactionsResponse = await response.json();
                setTransactions(data.data || []);
            } else {
                throw new Error('Failed to fetch transactions');
            }
        } catch (err: unknown) {
            console.error('Error fetching transactions:', err);
        }
    }, []);

    const refreshData = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            await Promise.all([fetchWallets(), fetchTransactions()]);
        } catch (err: any) {
            setError(err.message || 'Failed to load wallet data');
        } finally {
            setLoading(false);
        }
    }, [fetchWallets, fetchTransactions]);

    useEffect(() => {
        if (AuthService.isAuthenticated()) {
            refreshData();
        } else {
            setLoading(false);
        }
    }, [refreshData]);

    return {
        wallets,
        transactions,
        loading,
        error,
        refreshData,
        fetchTransactions, // Expose for pagination if needed
    };
}
