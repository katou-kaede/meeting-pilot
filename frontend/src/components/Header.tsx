import { useNavigate } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import {
  CalendarDays,
  ChevronDown,
  LogOut,
  UserRound,
  UserX,
} from "lucide-react";

import { useAuth } from "../contexts/AuthContext";

export default function Header() {
  const { user, logout, deactivateAccount } = useAuth();
  const navigate = useNavigate();
  const menuRef = useRef<HTMLDivElement>(null);

  const [menuOpen, setMenuOpen] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // メニュー外をクリックしたら閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current &&
        !menuRef.current.contains(
          event.target as Node
        )
      ) {
        setMenuOpen(false);
      }
    };

    document.addEventListener(
      "mousedown",
      handleClickOutside
    );

    return () => {
      document.removeEventListener(
        "mousedown",
        handleClickOutside
      );
    };
  }, []);

  const handleLogout = async () => {
    if (processing) return;

    setErrorMessage("");
    setProcessing(true);

    try {
      await logout();

      navigate("/login", {
        replace: true,
      });
    } catch (error) {
      setErrorMessage(
        error instanceof Error
          ? error.message
          : "ログアウトに失敗しました"
      );
    } finally {
      setProcessing(false);
    }
  };

  const handleDeactivateAccount = async () => {
    if (processing) return;

    const confirmed = window.confirm(
      "退会しますか？\n退会すると、このアカウントではログインできなくなります。"
    );

    if (!confirmed) return;

    setErrorMessage("");
    setProcessing(true);

    try {
      await deactivateAccount();

      navigate("/login", {
        replace: true,
        state: {
          message: "退会処理が完了しました",
        },
      });
    } catch (error) {
      setErrorMessage(
        error instanceof Error
          ? error.message
          : "退会処理に失敗しました"
      );
    } finally {
      setProcessing(false);
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
            Meeting Pilot
          </p>
        </button>

        {/* ユーザーメニュー */}
        <div
          ref={menuRef}
          className="relative"
        >
          <button
            type="button"
            onClick={() => {
              setMenuOpen((current) => !current);
              setErrorMessage("");
            }}
            aria-expanded={menuOpen}
            aria-haspopup="menu"
            className="flex cursor-pointer items-center gap-3 rounded-2xl bg-white/5 px-3 py-2 text-left ring-1 ring-white/10 transition hover:bg-white/10"
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-full bg-white/10 text-slate-300">
              <UserRound size={17} />
            </span>

            <div className="hidden max-w-52 sm:block">
              <p className="truncate text-sm font-semibold text-white">
                {user?.name}
              </p>

              <p className="truncate text-xs text-slate-400">
                {user?.email}
              </p>
            </div>

            <ChevronDown
              size={16}
              className={`text-slate-400 transition ${
                menuOpen ? "rotate-180" : ""
              }`}
            />
          </button>

          {menuOpen && (
            <div
              role="menu"
              className="absolute right-0 mt-2 w-72 overflow-hidden rounded-2xl border border-slate-200 bg-white p-2 text-slate-800 shadow-xl"
            >
              <div className="border-b border-slate-100 px-3 py-3 sm:hidden">
                <p className="truncate text-sm font-semibold">
                  {user?.name}
                </p>

                <p className="truncate text-xs text-slate-500">
                  {user?.email}
                </p>
              </div>

              {errorMessage && (
                <p className="m-2 rounded-lg bg-red-50 px-3 py-2 text-xs text-red-700">
                  {errorMessage}
                </p>
              )}

              <button
                type="button"
                role="menuitem"
                onClick={handleLogout}
                disabled={processing}
                className="flex w-full cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 text-sm transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <LogOut
                  size={17}
                  className="text-slate-500"
                />
                ログアウト
              </button>

              <button
                type="button"
                role="menuitem"
                onClick={handleDeactivateAccount}
                disabled={processing}
                className="mt-1 flex w-full cursor-pointer items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-red-600 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <UserX size={17} />
                {processing
                  ? "処理中..."
                  : "退会する"}
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}