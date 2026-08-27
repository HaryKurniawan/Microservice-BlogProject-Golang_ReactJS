import { useNavigate } from 'react-router-dom';
import AuthForm from '../components/AuthForm';

export default function AuthPage() {
  const navigate = useNavigate();

  const handleAuthSuccess = () => {
    // Setelah sukses login, arahkan user kembali ke halaman utama
    navigate('/');
  };

  return (
    <div className="min-h-[80vh] flex flex-col justify-center items-center px-6 py-12">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-[var(--color-text-base)]">
            Post Manager
          </h1>
          <p className="text-xs text-[var(--color-text-muted)] mt-1">
            Silakan masuk atau daftar terlebih dahulu untuk berkontribusi membuat post.
          </p>
        </div>

        <AuthForm onAuthSuccess={handleAuthSuccess} />
      </div>
    </div>
  );
}
