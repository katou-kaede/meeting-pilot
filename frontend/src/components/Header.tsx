import { useNavigate } from "react-router-dom";
import {
  CalendarDays,
  LogOut,
  UserRound,
} from "lucide-react";

import { useAuth } from "../contexts/AuthContext";

export default function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await logout();

      navigate("/login", {
        replace: true,
      });
    } catch (error) {
      console.error("ログアウトに失敗しました", error);
    }
  };

  return (
    <header className="sticky top-0 z-50 border-b border-slate-700/60 bg-slate-900/95 text-white shadow-sm backdrop-blur-xl">
      <div className="mx-auto flex h-18 max-w-[1600px] items-center justify-between px-6 lg:px-8">
        {/* ロゴ */}
        <button
          type="button"
          onClick={() => navigate("/")}
          className="group inline-flex cursor-pointer items-center gap-3"
        >
          <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-white/10 text-white ring-1 ring-white/15 transition group-hover:bg-white/15">
            <CalendarDays size={20} />
          </span>

          <p className="text-xl font-bold tracking-tight text-white">
            MeetingPilot
          </p>
        </button>

        {/* ユーザー情報 */}
        <div className="flex items-center gap-3">
          <div className="hidden items-center gap-3 rounded-2xl bg-white/5 px-4 py-2 ring-1 ring-white/10 sm:flex">
            <span className="flex h-9 w-9 items-center justify-center rounded-full bg-white/10 text-slate-300">
              <UserRound size={17} />
            </span>

            <div className="max-w-52">
              <p className="truncate text-sm font-semibold text-white">
                {user?.name}
              </p>

              <p className="truncate text-xs text-slate-400">
                {user?.email}
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={handleLogout}
            className="inline-flex cursor-pointer items-center gap-2 rounded-xl px-3 py-2 text-sm font-medium text-slate-300 transition hover:bg-red-500/10 hover:text-red-300"
            title="ログアウト"
          >
            <LogOut size={17} />

            <span className="hidden md:inline">
              ログアウト
            </span>
          </button>
        </div>
      </div>
    </header>
  );
}