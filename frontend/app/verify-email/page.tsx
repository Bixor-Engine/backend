'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { AuthService } from '@/lib/auth';
import { toast } from "sonner";
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Loader2, Mail, Shield, CheckCircle2, ArrowLeft } from 'lucide-react';

export default function VerifyEmail() {
  const router = useRouter();
  const [code, setCode] = useState(['', '', '', '', '', '']);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [requesting, setRequesting] = useState(false);
  const [codeSent, setCodeSent] = useState(false);
  const [verified, setVerified] = useState(false);
  const [userEmail, setUserEmail] = useState<string | null>(null);

  // Get user email on mount
  useEffect(() => {
    const fetchUser = async () => {
      // Wait a bit for tokens to be stored after login redirect
      await new Promise(resolve => setTimeout(resolve, 300));
      
      const token = AuthService.getToken();
      if (!token) {
        router.push('/auth/signin');
        return;
      }
      
      const user = await AuthService.getCurrentUser();
      if (user) {
        setUserEmail(user.email);
        // Check if already verified
        if (user.email_status) {
          setVerified(true);
          setTimeout(() => {
            router.push('/wallets');
          }, 2000);
        }
      }
    };
    fetchUser();
  }, [router]);

  const handleCodeChange = (index: number, value: string) => {
    // Only allow digits
    if (value && !/^\d$/.test(value)) return;

    const newCode = [...code];
    newCode[index] = value;
    setCode(newCode);
    setError('');

    // Auto-focus next input
    if (value && index < 5) {
      const nextInput = document.getElementById(`code-${index + 1}`);
      nextInput?.focus();
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace' && !code[index] && index > 0) {
      const prevInput = document.getElementById(`code-${index - 1}`);
      prevInput?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').trim();
    if (/^\d{6}$/.test(pastedData)) {
      const digits = pastedData.split('');
      const newCode = [...code];
      digits.forEach((digit, index) => {
        if (index < 6) {
          newCode[index] = digit;
        }
      });
      setCode(newCode);
      setError('');
      // Focus last input
      document.getElementById('code-5')?.focus();
    }
  };

  async function handleRequestCode() {
    setRequesting(true);
    setError('');
    
    try {
      await AuthService.requestOTP('email-verification');
      setCodeSent(true);
      toast.success("Verification code sent!", {
        description: `Check your email${userEmail ? ` at ${userEmail}` : ''} for the 6-digit code`,
        duration: 6000
      });
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to send verification code');
      }
    } finally {
      setRequesting(false);
    }
  }

  async function handleVerify(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');

    const codeString = code.join('');
    if (codeString.length !== 6) {
      setError('Please enter the complete 6-digit code');
      setLoading(false);
      return;
    }

    try {
      await AuthService.verifyOTP('email-verification', codeString);
      
      // Refresh user data to get updated email_status
      const updatedUser = await AuthService.getCurrentUser();
      if (updatedUser) {
        setUserEmail(updatedUser.email);
      }
      
      setVerified(true);
      toast.success("Email verified successfully!", {
        description: "Your email has been verified",
        duration: 4000
      });
      
      // Redirect to wallets after 2 seconds
      setTimeout(() => {
        router.push('/wallets');
      }, 2000);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
        // Clear code on error
        setCode(['', '', '', '', '', '']);
        document.getElementById('code-0')?.focus();
      } else {
        setError('Failed to verify code');
      }
    } finally {
      setLoading(false);
    }
  }

  if (verified) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background/50 to-secondary/20 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <Card className="shadow-2xl border-0 bg-card/50 backdrop-blur-sm">
            <CardContent className="pt-6 pb-6">
              <div className="text-center space-y-4">
                <div className="flex justify-center">
                  <div className="rounded-full bg-green-100 dark:bg-green-900/20 p-3">
                    <CheckCircle2 className="h-12 w-12 text-green-600 dark:text-green-400" />
                  </div>
                </div>
                <div>
                  <h3 className="text-xl font-semibold">Email Verified!</h3>
                  <p className="text-sm text-muted-foreground mt-2">
                    Your email has been successfully verified. Redirecting...
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background/50 to-secondary/20 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="text-center mb-8">
          <Link 
            href="/" 
            className="text-3xl font-bold bg-gradient-to-r from-primary to-primary/80 bg-clip-text text-transparent hover:from-primary/80 hover:to-primary transition-all duration-300"
          >
            Bixor Engine
          </Link>
          <p className="mt-2 text-sm text-muted-foreground">
            Verify your email address to continue
          </p>
        </div>

        <Card className="shadow-2xl border-0 bg-card/50 backdrop-blur-sm">
          <CardHeader className="space-y-1 pb-4">
            <CardTitle className="text-2xl font-bold text-center flex items-center justify-center gap-2">
              <Shield className="h-6 w-6" />
              Verify Email
            </CardTitle>
            <CardDescription className="text-center">
              {codeSent 
                ? "Enter the 6-digit code sent to your email"
                : "Request a verification code to verify your email address"
              }
            </CardDescription>
          </CardHeader>
          
          <CardContent className="space-y-4">
            {userEmail && (
              <div className="p-3 rounded-md bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800">
                <p className="text-sm text-blue-700 dark:text-blue-400 flex items-center gap-2">
                  <Mail className="h-4 w-4" />
                  Code will be sent to: <strong>{userEmail}</strong>
                </p>
              </div>
            )}

            {error && (
              <div className="p-3 rounded-md bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800">
                <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
              </div>
            )}

            {!codeSent ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground text-center">
                  Click the button below to receive a verification code via email.
                </p>
                <Button
                  onClick={handleRequestCode}
                  disabled={requesting}
                  className="w-full h-11 bg-gradient-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary text-primary-foreground font-medium transition-all duration-300"
                >
                  {requesting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Sending code...
                    </>
                  ) : (
                    <>
                      <Mail className="mr-2 h-4 w-4" />
                      Send Verification Code
                    </>
                  )}
                </Button>
              </div>
            ) : (
              <form onSubmit={handleVerify} className="space-y-4">
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-center block">
                    Enter 6-digit code
                  </Label>
                  <div className="flex gap-2 justify-center" onPaste={handlePaste}>
                    {code.map((digit, index) => (
                      <Input
                        key={index}
                        id={`code-${index}`}
                        type="text"
                        inputMode="numeric"
                        maxLength={1}
                        value={digit}
                        onChange={(e) => handleCodeChange(index, e.target.value)}
                        onKeyDown={(e) => handleKeyDown(index, e)}
                        className="w-12 h-14 text-center text-lg font-semibold"
                        autoFocus={index === 0}
                      />
                    ))}
                  </div>
                  <p className="text-xs text-muted-foreground text-center">
                    Code expires in 10 minutes
                  </p>
                </div>

                <Button
                  type="submit"
                  disabled={loading || code.join('').length !== 6}
                  className="w-full h-11 bg-gradient-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary text-primary-foreground font-medium transition-all duration-300"
                >
                  {loading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Verifying...
                    </>
                  ) : (
                    <>
                      <CheckCircle2 className="mr-2 h-4 w-4" />
                      Verify Email
                    </>
                  )}
                </Button>

                <div className="text-center">
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={handleRequestCode}
                    disabled={requesting}
                    className="text-sm text-muted-foreground hover:text-foreground"
                  >
                    {requesting ? (
                      <>
                        <Loader2 className="mr-2 h-3 w-3 animate-spin" />
                        Sending...
                      </>
                    ) : (
                      'Resend code'
                    )}
                  </Button>
                </div>
              </form>
            )}

            <div className="relative my-6">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border" />
              </div>
            </div>

            <div className="text-center">
              <Link href="/auth/signin">
                <Button variant="outline" className="w-full h-11 font-medium">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Back to Sign In
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

