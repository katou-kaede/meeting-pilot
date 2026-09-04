import { useState, type SubmitEvent } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";

import { useAuth } from "../contexts/AuthContext";
import ErrorMessage from "../components/ErrorMessage";
import Loading from "../components/Loading";
import { API_BASE_URL } from "../config/env";

export default function UserCreatePage() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordConfirmation, setPasswordConfirmation] =
    useState("");

  const [errorMessage, setErrorMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (
    event: SubmitEvent<HTMLFormElement>
  ) => {
    event.preventDefault();

    if (submitting) return;

    setErrorMessage("");

    if (password !== passwordConfirmation) {
      setErrorMessage("確認用パスワードが一致しません");
      return;
    }

    setSubmitting(true);

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/users`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name,
            email,
            password,
          }),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setErrorMessage(
          data.error || "アカウントの登録に失敗しました"
        );
        return;
      }

      navigate("/login", {
        replace: true,
        state: {
          message:
            "アカウントを登録しました。ログインしてください。",
        },
      });
    } catch (error) {
      console.error(error);
      setErrorMessage(
        "サーバーとの通信に失敗しました"
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <Loading />;
  }

  // ログイン済みなら会議一覧へ移動
  if (user) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-6">
      <div className="w-full max-w-md rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-slate-900">
            アカウント登録
          </h1>

          <p className="mt-2 text-sm text-slate-500">
            MeetGuideを利用するユーザーを登録します
          </p>
        </div>

        {errorMessage && (
          <ErrorMessage message={errorMessage} />
        )}

        <form
          onSubmit={handleSubmit}
          className="space-y-5"
        >
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              氏名
            </label>

            <input
              type="text"
              value={name}
              onChange={(event) =>
                setName(event.target.value)
              }
              autoComplete="name"
              required
              maxLength={100}
              className="w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-slate-500"
            />
          </div>

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
              maxLength={255}
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
              autoComplete="new-password"
              required
              minLength={8}
              className="w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-slate-500"
            />

            <p className="mt-1 text-xs text-slate-500">
              8文字以上で入力してください
            </p>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">
              パスワード（確認）
            </label>

            <input
              type="password"
              value={passwordConfirmation}
              onChange={(event) =>
                setPasswordConfirmation(event.target.value)
              }
              autoComplete="new-password"
              required
              minLength={8}
              className="w-full rounded-xl border border-slate-300 px-3 py-2 outline-none focus:border-slate-500"
            />
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="w-full cursor-pointer rounded-xl bg-slate-900 px-4 py-3 font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {submitting
              ? "登録中..."
              : "アカウントを登録"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          既にアカウントをお持ちですか？
          <Link
            to="/login"
            className="ml-1 font-medium text-slate-900 hover:underline"
          >
            ログイン
          </Link>
        </p>
      </div>
    </div>
  );
}