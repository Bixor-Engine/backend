'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { useWallet } from '@/hooks/use-wallet';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ProtectedNavbar } from '@/components/protected-navbar';
import {
  Plus,
  Send,
  ArrowDownLeft,
  ArrowUpRight,
  Eye,
  EyeOff,
  DollarSign,
  History,
  QrCode
} from 'lucide-react';

export default function Wallets() {
  const { user, loading: authLoading, requireAuth, requireEmailVerification } = useAuth();
  const { wallets, transactions, loading: walletLoading } = useWallet();
  const [showBalances, setShowBalances] = useState(true);

  useEffect(() => {
    if (!requireAuth()) return;
    requireEmailVerification();
  }, [requireAuth, requireEmailVerification]);

  if (authLoading || walletLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-muted-foreground">Loading wallets...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <ProtectedNavbar user={user} currentPage="wallets" />

      {/* Main Content */}
      <main className="container mx-auto p-6 space-y-6">
        {/* Portfolio Summary */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold">Your Wallets</h2>
            <p className="text-muted-foreground">Manage your cryptocurrency holdings</p>
          </div>
          <div className="flex space-x-2">
            {/* Action buttons could go here if needed in future */}
          </div>
        </div>

        {/* Total Portfolio Value - Placeholder for now */}
        <Card className="bg-gradient-to-r from-primary/5 to-primary/10 border-primary/20">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Total Portfolio Value</CardTitle>
                <CardDescription>All your assets combined (Estimated)</CardDescription>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setShowBalances(!showBalances)}
              >
                {showBalances ? <Eye className="h-4 w-4" /> : <EyeOff className="h-4 w-4" />}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex items-center space-x-4">
              <div className="text-3xl font-bold">
                {showBalances ? '$0.00' : '••••••'} <span className="text-sm text-muted-foreground font-normal">(Prices unavailable)</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Wallets and Transactions */}
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Wallets List */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Your Wallets</CardTitle>
                <CardDescription>Manage your cryptocurrency balances</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {wallets.map((wallet) => (
                  <div
                    key={wallet.id || wallet.coin_ticker}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center space-x-4">
                      <div className="flex items-center justify-center flex-shrink-0 w-8 h-8">
                        {wallet.coin_logo ? (
                          <img
                            src={wallet.coin_logo}
                            alt={wallet.coin_ticker}
                            className="h-full w-full object-contain"
                            onError={(e) => {
                              // Fallback if image fails to load
                              (e.target as HTMLImageElement).src = '';
                              (e.target as HTMLImageElement).className = 'hidden';
                            }}
                          />
                        ) : (
                          <DollarSign className="h-6 w-6 text-primary" />
                        )}
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <h3 className="font-medium">{wallet.coin_name}</h3>
                          <Badge variant="secondary">{wallet.coin_ticker}</Badge>
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {showBalances ? wallet.balance : '••••••'} {wallet.coin_ticker}
                        </p>
                      </div>
                    </div>

                    <div className="text-right">
                      {/* USD Value placeholder */}
                    </div>

                    <div className="flex space-x-2">
                      <Button variant="outline" size="sm">
                        <Send className="h-4 w-4 mr-1" />
                        Send
                      </Button>
                      <Button variant="outline" size="sm">
                        <ArrowDownLeft className="h-4 w-4 mr-1" />
                        Receive
                      </Button>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>

          {/* Recent Transactions */}
          <div>
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Recent Transactions</CardTitle>
                    <CardDescription>Your latest activity</CardDescription>
                  </div>
                  <Button variant="ghost" size="sm">
                    <History className="h-4 w-4 mr-1" />
                    View All
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {transactions.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No transactions found.
                  </div>
                ) : (
                  transactions.map((tx) => (
                    <div
                      key={tx.id}
                      className="flex items-center space-x-3 p-3 border rounded-lg"
                    >
                      <div className={`p-2 rounded-full ${tx.type === 'deposit' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
                        }`}>
                        {tx.type === 'deposit' ? (
                          <ArrowDownLeft className="h-4 w-4" />
                        ) : (
                          <ArrowUpRight className="h-4 w-4" />
                        )}
                      </div>

                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <p className="text-sm font-medium capitalize">
                            {tx.type} {tx.coin_ticker}
                          </p>
                          <div className="text-right">
                            <p className="text-sm font-medium">
                              {tx.type === 'deposit' ? '+' : '-'}{tx.amount} {tx.coin_ticker}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center justify-between mt-1">
                          <Badge variant={tx.status === 'completed' ? 'default' : 'secondary'}>
                            {tx.status}
                          </Badge>
                          <p className="text-xs text-muted-foreground">{new Date(tx.created_at).toLocaleDateString()}</p>
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Common wallet operations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <Button variant="outline" className="h-20 flex-col">
                <Send className="h-6 w-6 mb-2" />
                Send Crypto
              </Button>
              <Button variant="outline" className="h-20 flex-col">
                <ArrowDownLeft className="h-6 w-6 mb-2" />
                Receive Crypto
              </Button>
              <Button variant="outline" className="h-20 flex-col">
                <QrCode className="h-6 w-6 mb-2" />
                Scan QR Code
              </Button>
              <Button variant="outline" className="h-20 flex-col">
                <Plus className="h-6 w-6 mb-2" />
                Add Token
              </Button>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
