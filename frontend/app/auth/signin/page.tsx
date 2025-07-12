import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function SignIn() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const res = await fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Login failed');
      // Store tokens as needed (localStorage/session/cookie)
      // Redirect to home or dashboard
      router.push('/');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={{ padding: 32 }}>
      <h2>Sign In</h2>
      <form onSubmit={handleSubmit} style={{ maxWidth: 320 }}>
        <label>Email<br />
          <input type="email" value={email} onChange={e => setEmail(e.target.value)} required />
        </label><br /><br />
        <label>Password<br />
          <input type="password" value={password} onChange={e => setPassword(e.target.value)} required />
        </label><br /><br />
        <button type="submit" disabled={loading}>{loading ? 'Signing in...' : 'Sign In'}</button>
        {error && <p style={{ color: 'red' }}>{error}</p>}
      </form>
    </main>
  );
}
