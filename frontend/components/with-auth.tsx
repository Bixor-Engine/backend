import { useEffect } from 'react';
import { useAuth } from '@/hooks/use-auth';

interface WithAuthProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function withAuth<P extends object>(
  WrappedComponent: React.ComponentType<P>
) {
  return function WithAuthComponent(props: P & WithAuthProps) {
    const { user, loading, requireAuth } = useAuth();

    useEffect(() => {
      requireAuth();
    }, [requireAuth]);

    if (loading) {
      return (
        props.fallback || (
          <div className="min-h-screen bg-background flex items-center justify-center">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
              <p className="mt-2 text-muted-foreground">Loading...</p>
            </div>
          </div>
        )
      );
    }

    if (!user) {
      return null;
    }

    return <WrappedComponent {...props} user={user} />;
  };
}

export function ProtectedRoute({ children, fallback }: WithAuthProps) {
  const { user, loading, requireAuth } = useAuth();

  useEffect(() => {
    requireAuth();
  }, [requireAuth]);

  if (loading) {
    return (
      fallback || (
        <div className="min-h-screen bg-background flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
            <p className="mt-2 text-muted-foreground">Loading...</p>
          </div>
        </div>
      )
    );
  }

  if (!user) {
    return null;
  }

  return <>{children}</>;
}
