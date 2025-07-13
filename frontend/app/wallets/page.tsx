'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/hooks/use-auth';
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
  TrendingUp,
  TrendingDown,
  DollarSign,
  Bitcoin,
  Copy,
  QrCode,
  History
} from 'lucide-react';
import { toast } from "sonner";

export default function Wallets() {
  const { user, loading, requireAuth } = useAuth();
  const [showBalances, setShowBalances] = useState(true);

  useEffect(() => {
    requireAuth();
  }, [requireAuth]);

  if (loading) {
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

  // Mock wallet data
  const wallets = [
    {
      id: 1,
      name: 'Bitcoin',
      symbol: 'BTC',
      balance: '0.12458902',
      usdValue: '$5,847.23',
      change24h: '+5.67%',
      changeType: 'positive' as const,
      icon: Bitcoin,
      color: 'text-orange-500',
    },
    {
      id: 2,
      name: 'Ethereum',
      symbol: 'ETH',
      balance: '2.45891234',
      usdValue: '$4,234.56',
      change24h: '+3.21%',
      changeType: 'positive' as const,
      icon: DollarSign,
      color: 'text-blue-500',
    },
    {
      id: 3,
      name: 'Cardano',
      symbol: 'ADA',
      balance: '1,250.00',
      usdValue: '$567.89',
      change24h: '-2.14%',
      changeType: 'negative' as const,
      icon: DollarSign,
      color: 'text-purple-500',
    },
  ];

  const recentTransactions = [
    { 
      id: 1, 
      type: 'receive', 
      asset: 'BTC', 
      amount: '0.00521', 
      usdValue: '$245.67', 
      from: '1A1zP1...eP2SH', 
      time: '2 min ago',
      status: 'completed'
    },
    { 
      id: 2, 
      type: 'send', 
      asset: 'ETH', 
      amount: '0.5', 
      usdValue: '$860.23', 
      to: '0x742d...35Cc', 
      time: '1 hour ago',
      status: 'completed'
    },
    { 
      id: 3, 
      type: 'receive', 
      asset: 'ADA', 
      amount: '100', 
      usdValue: '$45.50', 
      from: 'addr1q...xyz', 
      time: '3 hours ago',
      status: 'pending'
    },
  ];

  const totalPortfolioValue = '$10,649.68';
  const portfolioChange = '+4.23%';

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard!');
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <ProtectedNavbar user={user} currentPage="wallets" />

      {/* Main Content */}
      <main className="container mx-auto p-6 space-y-6">{/* Portfolio Summary */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold">Your Wallets</h2>
            <p className="text-muted-foreground">Manage your cryptocurrency holdings</p>
          </div>
          <Button className="bg-gradient-to-r from-primary to-primary/90">
            <Plus className="mr-2 h-4 w-4" />
            Add Wallet
          </Button>
        </div>

        {/* Total Portfolio Value */}
        <Card className="bg-gradient-to-r from-primary/5 to-primary/10 border-primary/20">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-lg">Total Portfolio Value</CardTitle>
                <CardDescription>All your assets combined</CardDescription>
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
                {showBalances ? totalPortfolioValue : '••••••'}
              </div>
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-4 w-4 text-green-500" />
                <span className="text-green-500 font-medium">{portfolioChange}</span>
                <span className="text-muted-foreground text-sm">24h</span>
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
                    key={wallet.id}
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center space-x-4">
                      <div className={`p-2 rounded-full bg-muted ${wallet.color}`}>
                        <wallet.icon className="h-5 w-5" />
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <h3 className="font-medium">{wallet.name}</h3>
                          <Badge variant="secondary">{wallet.symbol}</Badge>
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {showBalances ? wallet.balance : '••••••'} {wallet.symbol}
                        </p>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <p className="font-medium">
                        {showBalances ? wallet.usdValue : '••••••'}
                      </p>
                      <div className="flex items-center space-x-1">
                        {wallet.changeType === 'positive' ? (
                          <TrendingUp className="h-3 w-3 text-green-500" />
                        ) : (
                          <TrendingDown className="h-3 w-3 text-red-500" />
                        )}
                        <span className={`text-xs ${
                          wallet.changeType === 'positive' ? 'text-green-500' : 'text-red-500'
                        }`}>
                          {wallet.change24h}
                        </span>
                      </div>
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
                {recentTransactions.map((tx) => (
                  <div
                    key={tx.id}
                    className="flex items-center space-x-3 p-3 border rounded-lg"
                  >
                    <div className={`p-2 rounded-full ${
                      tx.type === 'receive' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'
                    }`}>
                      {tx.type === 'receive' ? (
                        <ArrowDownLeft className="h-4 w-4" />
                      ) : (
                        <ArrowUpRight className="h-4 w-4" />
                      )}
                    </div>
                    
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <p className="text-sm font-medium capitalize">
                          {tx.type} {tx.asset}
                        </p>
                        <div className="text-right">
                          <p className="text-sm font-medium">
                            {tx.type === 'receive' ? '+' : '-'}{tx.amount} {tx.asset}
                          </p>
                          <p className="text-xs text-muted-foreground">{tx.usdValue}</p>
                        </div>
                      </div>
                      
                      <div className="flex items-center justify-between mt-1">
                        <div className="flex items-center space-x-2">
                          <p className="text-xs text-muted-foreground">
                            {tx.type === 'receive' ? 'From:' : 'To:'} {tx.from || tx.to}
                          </p>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-auto p-0"
                            onClick={() => copyToClipboard(tx.from || tx.to || '')}
                          >
                            <Copy className="h-3 w-3" />
                          </Button>
                        </div>
                        <Badge variant={tx.status === 'completed' ? 'default' : 'secondary'}>
                          {tx.status}
                        </Badge>
                      </div>
                      
                      <p className="text-xs text-muted-foreground mt-1">{tx.time}</p>
                    </div>
                  </div>
                ))}
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
