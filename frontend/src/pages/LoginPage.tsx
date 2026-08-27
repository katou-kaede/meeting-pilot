import { useState, useEffect } from "react";
import { Navigate, useNavigate, Link, useLocation } from "react-router-dom";

import { useAuth } from "../contexts/AuthContext";
import ErrorMessage from "../components/ErrorMessage";
import Loading from "../components/Loading";

export default function LoginPage() {
  const { user, loading, login } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [errorMessage, setErrorMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const location = useLocation();

  const [successMessage, setSuccessMessage] = useState(
    (location.state as { message?: string } | null)?.message ?? ""
  );
  
  const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (submitting) return;

    setErrorMessage("");
    setSubmitting(true);

    try {
      await login(email, password);

      navigate("/", {
        replace: true,
      });

    } catch (error) {
      if (error instanceof Error) {
        setErrorMessage(error.message);
      } else {
        setErrorMessage("ログインに失敗しました");
      }

    } finally {
      setSubmitting(false);
    }
  };

  // 成功メッセージを3秒後に消す
  useEffect(() => {
    if (!successMessage) return;

    // 「一時保存しました」を3秒後に消すためのタイマー
    const timerId = window.setTimeout(() => {
      setSuccessMessage("");

      navigate("/login", {
        replace: true,
        state: null
      })
    }, 3000);

    return () => {
      window.clearTimeout(timerId);
    };
  }, [successMessage, navigate]);

  if (loading) {
    return <Loading />;
  }

  // ログイン済みの場合は会議一覧へ移動
  if (user) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-6">
      <div className="w-full max-w-md rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-slate-900">
            MeetGide
          </h1>

          <p className="mt-2 text-sm text-slate-500">
            アカウントへログインしてください
          </p>
        </div>

        {/* エラーメッセージ */}
        {errorMessage && (
          <ErrorMessage message={errorMessage} />
        )}

        {/* 成功メッセージ */}
        {successMessage && (
          <div className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-emerald-700">
            {successMessage}
          </div>
        )}

        <form
          onSubmit={handleSubmit}
          className="space-y-5"
        >
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              メールアドレス
            </label>

            <input
              type="email"
              value={email}
              onChange={(event) =>
                setEmail(event.target.value)
              }
              autoComplete="email"
              required
              className="w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-slate-500"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              パスワード
            </label>

            <input
              type="password"
              value={password}
              onChange={(event) =>
                setPassword(event.target.value)
              }
              autoComplete="current-password"
              required
              className="w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-slate-500"
            />
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="w-full cursor-pointer rounded-xl bg-slate-900 px-4 py-3 font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {submitting ? "ログイン中..." : "ログイン"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          アカウントをお持ちでない方は
          <Link
            to="/users/create"
            className="ml-1 font-medium text-slate-900 hover:underline"
          >
            新規登録
          </Link>
        </p>
      </div>
    </div>
  );
}