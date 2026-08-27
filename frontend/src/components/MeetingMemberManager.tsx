import { useCallback, useEffect, useState } from "react";
import {
  Search,
  UserPlus,
  UserRound,
  X,
} from "lucide-react";

import type {
  MeetingMember,
  MemberCandidate,
} from "../types/user";
import ErrorMessage from "./ErrorMessage";
import { useAuth } from "./../contexts/AuthContext";
import { useNavigate } from "react-router-dom";

type Props = {
  meetingId: number;
};

export default function MeetingMemberManager({
  meetingId,
}: Props) {
  const [members, setMembers] = useState<MeetingMember[]>([]);
  const [candidates, setCandidates] = useState<MemberCandidate[]>([]);

  const [keyword, setKeyword] = useState("");
  const [selectedUserId, setSelectedUserId] =
    useState<number | null>(null);
  const [selectedRole, setSelectedRole] =
    useState<"editor" | "viewer">("viewer");

  const [loading, setLoading] = useState(true);
  const [searching, setSearching] = useState(false);
  const [saving, setSaving] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const hasEditor = members.some(
    (member) => member.role === "editor"
  );

  const { user } = useAuth();
  const navigate = useNavigate();

  // ============================================
  // 登録済みメンバー取得
  // ============================================
  const fetchMembers = useCallback(async () => {
    setErrorMessage("");

    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${meetingId}/members`,
        {
          credentials: "include",
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setErrorMessage(
          data.error || "参加メンバーの取得に失敗しました"
        );
        return;
      }

      setMembers(data as MeetingMember[]);
    } catch (error) {
      console.error(error);
      setErrorMessage("参加メンバーの取得に失敗しました");
    } finally {
      setLoading(false);
    }
  }, [meetingId]);

  useEffect(() => {
    fetchMembers();
  }, [fetchMembers]);

  // ============================================
  // ユーザー候補検索
  // ============================================
  const handleSearch = async () => {
    setErrorMessage("");
    setSearching(true);

    try {
      const params = new URLSearchParams({
        keyword: keyword.trim(),
      });

      const response = await fetch(
        `http://localhost:8080/api/meetings/${meetingId}/member-candidates?${params}`,
        {
          credentials: "include",
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setErrorMessage(
          data.error || "ユーザーの検索に失敗しました"
        );
        return;
      }

      setCandidates(data as MemberCandidate[]);
      setSelectedUserId(null);
    } catch (error) {
      console.error(error);
      setErrorMessage("ユーザーの検索に失敗しました");
    } finally {
      setSearching(false);
    }
  };

  // ============================================
  // メンバー追加
  // ============================================
  const handleAddMember = async () => {
    if (!selectedUserId || saving) {
      return;
    }

    setErrorMessage("");
    setSaving(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${meetingId}/members`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            user_id: selectedUserId,
            role: selectedRole,
          }),
        }
      );

      if (!response.ok) {
        const data = await response.json();

        setErrorMessage(
          data.error || "メンバーの追加に失敗しました"
        );
        return;
      }

      setSelectedUserId(null);

      // 追加したユーザーを候補一覧から除外
      setCandidates((prev) =>
        prev.filter(
          (candidate) => candidate.id !== selectedUserId
        )
      );

      await fetchMembers();
    } catch (error) {
      console.error(error);
      setErrorMessage("メンバーの追加に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  // ============================================
  // メンバー削除
  // ============================================
  const handleDeleteMember = async (
    member: MeetingMember
  ) => {

    const isSelf = member.user_id === user?.id;
    if (isSelf) {
      const confirmed = confirm(
        "この会議から脱退しますか？\n脱退すると、この会議を閲覧・編集できなくなります。"
      );
      if (!confirmed) return;
    }
    
    setErrorMessage("");
    setSaving(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${meetingId}/members/${member.user_id}`,
        {
          method: "DELETE",
          credentials: "include",
        }
      );

      if (!response.ok) {
        const data = await response.json();

        setErrorMessage(
          data.error || "メンバーの削除に失敗しました"
        );
        return;
      }

      if (isSelf) {
        navigate("/", {
          replace: true
        })
      }

      await fetchMembers();
    } catch (error) {
      console.error(error);
      setErrorMessage("メンバーの削除に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <section className="mt-6 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-slate-900">
          参加メンバー管理
        </h2>
      </div>

      {errorMessage && (
        <ErrorMessage message={errorMessage} />
      )}

      {/* ユーザー検索 */}
      <div className="rounded-2xl bg-slate-50 p-5">
        <label className="mb-2 block text-sm font-medium text-slate-700">
          ユーザー検索
        </label>

        <div className="flex gap-2">
          <div className="relative flex-1">
            <Search
              size={17}
              className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
            />

            <input
              type="search"
              value={keyword}
              onChange={(event) =>
                setKeyword(event.target.value)
              }
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  event.preventDefault();
                  handleSearch();
                }
              }}
              placeholder="氏名またはメールアドレス"
              className="w-full rounded-xl border border-slate-300 bg-white py-2 pl-10 pr-3"
            />
          </div>

          <button
            type="button"
            onClick={handleSearch}
            disabled={searching}
            className="cursor-pointer rounded-xl bg-slate-900 px-5 py-2 font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {searching ? "検索中..." : "検索"}
          </button>
        </div>

        {/* 検索結果 */}
        {candidates.length > 0 && (
          <div className="mt-4 space-y-2">
            {candidates.map((candidate) => {
              const isSelected =
                selectedUserId === candidate.id;

              return (
                <button
                  key={candidate.id}
                  type="button"
                  onClick={() =>
                    setSelectedUserId(candidate.id)
                  }
                  className={`flex w-full cursor-pointer items-center gap-3 rounded-xl border p-3 text-left transition ${
                    isSelected
                      ? "border-slate-900 bg-white ring-1 ring-slate-900"
                      : "border-slate-200 bg-white hover:border-slate-300"
                  }`}
                >
                  <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-slate-100 text-slate-500">
                    <UserRound size={17} />
                  </span>

                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium text-slate-900">
                      {candidate.name}
                    </p>

                    <p className="truncate text-xs text-slate-500">
                      {candidate.email}
                    </p>
                  </div>
                </button>
              );
            })}
          </div>
        )}

        {/* ロール・追加 */}
        {selectedUserId && (
          <div className="mt-4 border-t border-slate-200 pt-4">
            <div className="mb-1 flex items-center gap-3">
              <label className="block text-sm font-medium text-slate-700">
                ロール
              </label>

              <p className="text-xs font-medium text-red-600">
                ※ 編集者は会議ごとに1人のみ設定できます
              </p>
            </div>

            <div className="mb-1 flex items-center gap-2">
              <select
                value={selectedRole}
                onChange={(event) =>
                  setSelectedRole(
                    event.target.value as
                      | "editor"
                      | "viewer"
                  )
                }
                className="rounded-xl border border-slate-300 bg-white px-3 py-2"
              >
                <option value="viewer">参加者</option>
                <option
                  value="editor"
                  disabled={hasEditor}
                >
                  編集者
                  {hasEditor ? "（設定済み）" : ""}
                </option>
              </select>

              <button
                type="button"
                onClick={handleAddMember}
                disabled={saving}
                className="inline-flex cursor-pointer items-center gap-2 rounded-xl bg-slate-900 px-4 py-2 font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <UserPlus size={17} />
                {saving ? "処理中..." : "メンバー追加"}
              </button>
            </div>
          </div>
        )}
      </div>

      {/* 登録済みメンバー */}
      <div className="mt-6">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="font-semibold text-slate-900">
            登録済みメンバー
          </h3>

          <span className="rounded-lg bg-slate-100 px-2 py-1 text-xs text-slate-600">
            {members.length}人
          </span>
        </div>

        {loading ? (
          <div className="flex min-h-24 items-center justify-center">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-slate-200 border-t-slate-700" />
          </div>
        ) : (
          <div className="grid gap-3 md:grid-cols-2">
            {members.map((member) => (
              <div
                key={member.user_id}
                className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-slate-50 p-4"
              >
                <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-white text-slate-500 shadow-sm ring-1 ring-slate-200">
                  <UserRound size={18} />
                </span>

                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-slate-900">
                    {member.name}
                  </p>

                  <p className="truncate text-xs text-slate-500">
                    {member.email}
                  </p>
                </div>

                <span
                  className={`shrink-0 rounded-full px-2.5 py-1 text-xs font-medium ${
                    member.role === "owner"
                      ? "bg-blue-100 text-blue-700"
                      : member.role === "editor"
                        ? "bg-violet-100 text-violet-700"
                        : "bg-slate-200 text-slate-600"
                  }`}
                >
                  {member.role === "owner"
                    ? "主催者"
                    : member.role === "editor"
                      ? "編集者"
                      : "参加者"}
                </span>

                {/* 主催者は参加メンバーから除外できない */}
                {member.role !== "owner" ? (
                  <button
                    type="button"
                    onClick={() => handleDeleteMember(member)}
                    disabled={saving}
                    title={`${member.name}を削除`}
                    aria-label={`${member.name}を参加メンバーから削除`}
                    className="flex h-7 w-7 cursor-pointer items-center justify-center rounded-full text-slate-400 hover:bg-red-50 hover:text-red-600 disabled:cursor-not-allowed disabled:opacity-40"
                  >
                    <X size={14} />
                  </button>
                ) : (
                  <div className="w-7"></div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}