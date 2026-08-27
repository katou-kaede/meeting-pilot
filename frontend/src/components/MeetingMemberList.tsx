import { useEffect, useState } from "react";
import { UserRound } from "lucide-react";

import type { MeetingMember } from "../types/user";
import ErrorMessage from "./ErrorMessage";

type Props = {
  meetingId: number;
};

export default function MeetingMemberList({
  meetingId,
}: Props) {
  const [members, setMembers] = useState<MeetingMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const fetchMembers = async () => {
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
    };

    fetchMembers();
  }, [meetingId]);

  return (
    <section className="mt-6 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
      <div className="mb-5 flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">
            参加メンバー
          </h2>
        </div>

        <span className="rounded-lg bg-slate-100 px-2 py-1 text-xs text-slate-600">
          {members.length}人
        </span>
      </div>

      {loading ? (
        <div className="flex min-h-24 items-center justify-center">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-slate-200 border-t-slate-700" />
        </div>

      ) : errorMessage ? (
        <ErrorMessage message={errorMessage} />

      ) : members.length === 0 ? (
        <div className="rounded-2xl bg-slate-50 px-4 py-6 text-center text-sm text-slate-500">
          参加メンバーはまだ登録されていません
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
            </div>
          ))}
        </div>
      )}
    </section>
  );
}