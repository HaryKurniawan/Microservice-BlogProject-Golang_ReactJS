import { useNavigate } from 'react-router-dom';

export default function LoginPrompt() {
  const navigate = useNavigate();

  return (
    <div className="
      bg-[var(--color-bg-card)] border border-[var(--color-border)]
      rounded-[var(--radius-card)] p-6 text-center sticky top-6
    ">
      <div className="text-3xl mb-3">🔑</div>
      <h3 className="text-sm font-semibold text-[var(--color-text-base)] mb-1">
        Tulis Post Baru
      </h3>
      <p className="text-xs text-[var(--color-text-muted)] mb-4 leading-relaxed">
        Anda perlu login terlebih dahulu untuk dapat memposting artikel baru ke platform ini.
      </p>
      <button
        onClick={() => navigate('/auth')}
        className="
          w-full py-2 px-4 rounded-[var(--radius-btn)]
          text-xs font-medium text-white cursor-pointer
          bg-[var(--color-accent)]
          hover:bg-[var(--color-accent-hover)]
          transition-all duration-150
        "
      >
        Login Sekarang
      </button>
    </div>
  );
}
