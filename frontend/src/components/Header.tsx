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
    <header className="sticky top-0 z-50 border-b border-slate-200/80 bg-white/85 backdrop-blur-xl">
      <div className="mx-auto flex h-18 max-w-[1600px] items-center justify-between px-6 lg:px-8">
        {/* ロゴ */}
        <button
          type="button"
          onClick={() => navigate("/")}
          className="group inline-flex cursor-pointer items-center gap-3"
        >
          <span className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-900 text-white shadow-sm transition group-hover:bg-slate-800">
            <CalendarDays size={20} />
          </span>

          <div className="text-left">
            <p className="text-xl font-bold tracking-tight text-slate-900">
              MeetingPilot
            </p>
          </div>
        </button>

        {/* ユーザー情報 */}
        <div className="flex items-center gap-3">
          <div className="hidden items-center gap-3 rounded-2xl bg-slate-50 px-4 py-2 sm:flex">
            <span className="flex h-9 w-9 items-center justify-center rounded-full bg-white text-slate-500 shadow-sm ring-1 ring-slate-200">
              <UserRound size={17} />
            </span>

            <div className="max-w-52">
              <p className="truncate text-sm font-semibold text-slate-800">
                {user?.name}
              </p>

              <p className="truncate text-xs text-slate-500">
                {user?.email}
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={handleLogout}
            className="inline-flex cursor-pointer items-center gap-2 rounded-xl px-3 py-2 text-sm font-medium text-slate-500 transition hover:bg-red-50 hover:text-red-600"
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