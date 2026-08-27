import { useNavigate } from 'react-router-dom';

interface UserStatusProps {
  user: any | null;
  onLogout: () => void;
}

export default function UserStatus({ user, onLogout }: UserStatusProps) {
  const navigate = useNavigate();

  if (user) {
    return (
      <div className="flex items-center gap-3">
        <span className="text-sm text-[var(--color-text-muted)]">
          Halo, <strong className="text-[var(--color-text-base)]">{user.name}</strong>
        </span>
        <button
          onClick={onLogout}
          className="
            text-xs font-medium px-3 py-1.5 rounded-[var(--radius-btn)]
            bg-[var(--color-delete-bg)] text-[var(--color-error)]
            border border-[var(--color-delete-border)]
            hover:bg-[var(--color-error)] hover:text-white
            transition-all duration-150 cursor-pointer
          "
        >
          Logout
        </button>
      </div>
    );
  }

  return (
    <button
      onClick={() => navigate('/auth')}
      className="
        text-xs font-medium px-4 py-2 rounded-[var(--radius-btn)]
        bg-[var(--color-accent)] text-white
        hover:bg-[var(--color-accent-hover)]
        transition-all duration-150 cursor-pointer
      "
    >
      Login / Register
    </button>
  );
}
