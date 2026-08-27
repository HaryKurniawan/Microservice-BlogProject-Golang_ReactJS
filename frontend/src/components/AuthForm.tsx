import { useState } from 'react';
import { notification } from 'antd';
import * as postApi from '../api/postApi';

interface AuthFormProps {
  onAuthSuccess: (token: string, user: any) => void;
}

export default function AuthForm({ onAuthSuccess }: AuthFormProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim() || !password.trim() || (!isLogin && !name.trim())) return;

    setLoading(true);
    try {
      if (isLogin) {
        // Proses Login
        const data = await postApi.loginUser({ email, password });
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user));
        notification.success({
          title: 'Berhasil Login',
          description: `Selamat datang kembali, ${data.user.name}!`,
          placement: 'topRight',
        });
        onAuthSuccess(data.token, data.user);
      } else {
        // Proses Register
        await postApi.registerUser({ name, email, password });
        notification.success({
          title: 'Registrasi Berhasil',
          description: 'Akun Anda berhasil dibuat. Silakan login sekarang!',
          placement: 'topRight',
        });
        setIsLogin(true); // Alihkan ke form login
        setName('');
        setPassword('');
      }
    } catch (err: any) {
      notification.error({
        title: 'Gagal',
        description: err.response?.data?.error || 'Terjadi kesalahan pada sistem. Silakan coba lagi.',
        placement: 'topRight',
      });
    } finally {
      setLoading(false);
    }
  };

  const inputCls = `
    w-full bg-[var(--color-bg-input)] border border-[var(--color-border)]
    rounded-[var(--radius-input)] px-3.5 py-2
    text-sm text-[var(--color-text-base)] placeholder:text-[var(--color-text-subtle)]
    outline-none transition-all duration-150 font-[var(--font-sans)]
    focus:border-[var(--color-border-focus)] focus:bg-white
  `;

  return (
    <div className="
      bg-[var(--color-bg-card)] border border-[var(--color-border)]
      rounded-[var(--radius-card)] p-6 max-w-sm w-full mx-auto mb-8
    ">
      <h2 className="text-base font-semibold text-[var(--color-text-base)] mb-5 text-center">
        {isLogin ? 'Masuk ke Akun Anda' : 'Daftar Akun Baru'}
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {!isLogin && (
          <div>
            <label className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">
              Nama
            </label>
            <input
              type="text"
              className={inputCls}
              placeholder="Nama Lengkap"
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
          </div>
        )}

        <div>
          <label className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">
            Email
          </label>
          <input
            type="email"
            className={inputCls}
            placeholder="nama@email.com"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>

        <div>
          <label className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">
            Password
          </label>
          <input
            type="password"
            className={inputCls}
            placeholder="••••••••"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="
            w-full py-2 px-4 rounded-[var(--radius-btn)]
            text-sm font-medium text-white cursor-pointer
            bg-[var(--color-accent)]
            hover:bg-[var(--color-accent-hover)]
            active:scale-[0.98] disabled:opacity-50
            transition-all duration-150 mt-2
          "
        >
          {loading ? 'Memproses...' : isLogin ? 'Masuk' : 'Daftar'}
        </button>
      </form>

      <div className="mt-4 text-center">
        <button
          type="button"
          onClick={() => setIsLogin(!isLogin)}
          className="text-xs text-[var(--color-accent)] hover:underline cursor-pointer"
        >
          {isLogin ? 'Belum punya akun? Daftar di sini' : 'Sudah punya akun? Masuk di sini'}
        </button>
      </div>
    </div>
  );
}
