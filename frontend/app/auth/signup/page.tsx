import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function SignUp() {
  const router = useRouter();
  const [form, setForm] = useState({
    username: '',
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const res = await fetch('http://localhost:8080/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Registration failed');
      // Optionally auto-login or redirect
      router.push('/auth/signin');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={{ padding: 32 }}>
      <h2>Sign Up</h2>
      <form onSubmit={handleSubmit} style={{ maxWidth: 320 }}>
        <label>Username<br />
          <input type="text" value={form.username} onChange={e => setForm(f => ({ ...f, username: e.target.value }))} required />
        </label><br /><br />
        <label>Email<br />
          <input type="email" value={form.email} onChange={e => setForm(f => ({ ...f, email: e.target.value }))} required />
        </label><br /><br />
        <label>Password<br />
          <input type="password" value={form.password} onChange={e => setForm(f => ({ ...f, password: e.target.value }))} required />
        </label><br /><br />
        <button type="submit" disabled={loading}>{loading ? 'Signing up...' : 'Sign Up'}</button>
        {error && <p style={{ color: 'red' }}>{error}</p>}
      </form>
    </main>
  );
}
