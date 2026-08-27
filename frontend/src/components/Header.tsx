import UserStatus from './UserStatus';

interface HeaderProps {
  user: any | null;
  postsCount: number;
  onLogout: () => void;
}

export default function Header({ user, postsCount, onLogout }: HeaderProps) {
  return (
    <header className="flex items-center justify-between flex-wrap gap-4 mb-10 pb-6 border-b border-[var(--color-border)]">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-[var(--color-text-base)]">
          Post Manager
        </h1>
        <p className="text-sm text-[var(--color-text-subtle)] mt-0.5">
          Go REST API Microservice · React · PostgreSQL
        </p>
      </div>

      <div className="flex items-center gap-3">
        <UserStatus user={user} onLogout={onLogout} />

        <span className="
          text-xs font-medium text-[var(--color-accent)]
          bg-[var(--color-accent-muted)]
          px-3 py-1 rounded-full
        ">
          {postsCount} post{postsCount !== 1 ? 's' : ''}
        </span>
      </div>
    </header>
  );
}
