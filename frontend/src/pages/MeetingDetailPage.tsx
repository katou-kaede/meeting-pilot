import { useEffect, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import type { MeetingDetail } from "../types/meeting";
import { formatDateTime } from "../utils/common";
import { getStatusLabel, getStatusStyle } from "../utils/meetingStatus";
import AgendaSidebar from "../components/AgendaSidebar";
import ErrorMessage from "../components/ErrorMessage";
import Loading from "../components/Loading";
import MeetingMemberList from "../components/MeetingMemberList";
import {
  Play,
  ArrowLeft,
  MessageCirclePlus,
} from "lucide-react";


export default function MeetingDetailPage() {
  const { id } = useParams();

  const [meeting, setMeeting] = useState<MeetingDetail | null>(null);

  const navigate = useNavigate();
  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState("");
  const [connectionError, setConnectionError] = useState("");

  // 通信状態
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [selectedAgendaIndex, setSelectedAgendaIndex] = useState(0);

  useEffect(() => {
    const fetchMeeting = async () => {
      try {
        setErrorMessage("");

        const response = await fetch(
          `http://localhost:8080/api/meetings/${id}`,
          {
            credentials: "include",
          }
        );

        const data = await response.json();

        if (!response.ok) {
          setErrorMessage(data.error);
          return;
        }

        setMeeting(data);

      } catch (error) {
        console.error(error);
        setErrorMessage("会議情報の取得に失敗しました");

      } finally {
        setLoading(false);
      }
    };

    fetchMeeting();
  }, [id]);

  // WebSocket接続用のuseEffect
  useEffect(() => {
    if (!id) return;

    const socket = new WebSocket(
      `ws://localhost:8080/ws/meetings/${id}`
    );

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);

        if (message.type === "meeting_started") {
          navigate(`/meetings/${id}/session`);
        }
      } catch (error) {
        console.error(
          "WebSocket message parse failed:",
          error
        );
      }
    };

    socket.onerror = (error) => {
      setConnectionError("リアルタイム接続に問題が発生しました");
      console.error("WebSocket error:", error);
    };

    socket.onclose = (event) => {
      console.log("WebSocket disconnected", event.code);

      if (event.code !== 1000) {
        setConnectionError(
          "リアルタイム接続が切断されました。画面を再読み込みしてください"
        );
      }
    };

    return () => {
      socket.close();
    };
  }, [id, navigate]);

  // 削除処理
  const handleDelete = async () => {
    if (!confirm("削除しますか？")) {
      return;
    }

    try {
      setErrorMessage("");
      setSaving(true);

      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}`,
        {
          credentials: "include",
          method: "DELETE",
        }
      );

      if (response.ok) {
        navigate("/");
        return
      }

      const error = await response.json();
      setErrorMessage(error.error);

    } catch (error) {
      setErrorMessage("会議の削除に失敗しました");
      console.error(error);

    } finally {
      setSaving(false);
    }
  };

  // 開始処理
  const handleStart = async () => {
    if (!meeting) { return }

    try {
      setErrorMessage("");
      setSaving(true);

      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/start`,
        {
          credentials: "include",
          method: "PATCH",
        }
      );

      if (!response.ok) {
        const error = await response.json();
        setErrorMessage(error.error);
        return;
      }

      setMeeting({
        ...meeting,
        status: "in_progress",
      });

      navigate(`/meetings/${id}/session`);

    } catch (error) {
      setErrorMessage("会議の開始に失敗しました");
      console.error(error);

    } finally {
      setSaving(false);
    }
  };

  const agendas = meeting?.agendas ?? [];

  const selectedAgenda = agendas[selectedAgendaIndex] ?? null;

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="bg mx-auto max-w-[1600px] p-6 lg:p-8">

      <Link
        to="/"
        className="inline-flex items-center text-slate-500 hover:text-slate-900"
      >
        <ArrowLeft size={18} /> 会議一覧へ戻る
      </Link>

      {/* エラーメッセージ */}
      {(errorMessage || connectionError) && (
        <ErrorMessage message={errorMessage || connectionError} />
      )}

      {meeting && (
        <>
          {/* タイトルエリア */}
          <div className="flex justify-between items-start mt-6">
            <div>
              <h1 className="text-3xl font-bold text-slate-900">
                {meeting.title}
              </h1>

              <p className="mt-2 text-slate-500">
                {meeting.description}
              </p>
            </div>

            <div className="flex gap-2">
              {meeting.status === "scheduled" && (
                <button
                  onClick={handleStart}
                  disabled={saving}
                  className="flex items-center rounded-xl bg-blue-600 px-4 py-2 text-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                >
                  {saving ? ("処理中...") : (
                    <><Play size={15} className="me-1" />会議開始</>
                  )}
                </button>
              )}
              {meeting.status === "in_progress" && (
                <button
                  onClick={() => navigate(`/meetings/${meeting.id}/session`)}
                  className="flex items-center rounded-xl bg-blue-600 px-4 py-2 text-white cursor-pointer"
                >
                  <MessageCirclePlus size={18} className="me-1" /> 会議画面に参加
                </button>
              )}

              <Link
                to={`/meetings/${meeting.id}/edit`}
                className="px-4 py-2 rounded-xl border border-slate-300"
              >
                編集
              </Link>

              {meeting.current_user_role === "owner" && (
              <button
                onClick={handleDelete}
                disabled={saving}
                className="px-4 py-2 rounded-xl bg-red-600 text-white disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                {saving ? "処理中..." : "削除"}
              </button>
              )}
            </div>

          </div>

          {/* 会議情報カード */}
          <div className="mt-6 bg-white/70 border border-white/70 rounded-2xl p-4 shadow-sm backdrop-blur-xl">

            <div className="grid grid-cols-2 gap-4 text-sm md:grid-cols-4">

              <div>
                <p className="text-slate-500">会議相手</p>
                <p>{meeting.target_name}</p>
              </div>

              <div>
                <p className="text-slate-500">開始日時</p>
                <p>
                  {formatDateTime(meeting.scheduled_start_at)}
                </p>
              </div>

              <div>
                <p className="text-slate-500">予定時間</p>
                <p>{meeting.planned_minutes}分</p>
              </div>

              <div>
                <p className="text-slate-500">状態</p>

                <span className={`inline-block px-2 py-1 rounded-full text-sm ${getStatusStyle(meeting.status)}`}>
                  {getStatusLabel(meeting.status)}
                </span>

              </div>

            </div>

          </div>

          {/* 会議後入力項目 */}
          <div className="mt-6 grid gap-6 md:grid-cols-2">
            <div className="rounded-2xl border border-blue-100 bg-blue-50/70 p-6 shadow-sm backdrop-blur-xl">
              <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-blue-600">
                Meeting Result
              </p>

              <h2 className="mb-4 text-lg font-semibold text-slate-900">
                決定事項
              </h2>

              <p className="whitespace-pre-wrap text-slate-700">
                {meeting.decisions || ""}
              </p>
            </div>

            <div className="rounded-2xl border border-amber-100 bg-amber-50/70 p-6 shadow-sm backdrop-blur-xl">
              <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-amber-600">
                Next Action
              </p>

              <h2 className="mb-4 text-lg font-semibold text-slate-900">
                TODO
              </h2>

              <p className="whitespace-pre-wrap text-slate-700">
                {meeting.todo || ""}
              </p>
            </div>
          </div>

          {/* アジェンダ */}
          <div className="mt-6">

            <div className="mt-6 grid gap-6 lg:grid-cols-[320px_minmax(0,1fr)]">

              {/* 左ペイン */}
              <AgendaSidebar
                agendas={agendas}
                selectedIndex={selectedAgendaIndex}
                onSelect={setSelectedAgendaIndex}
              />

              {/* 右ペイン */}
              <main className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
                {selectedAgenda ? (
                  <>
                    <div className="border-b border-slate-200 pb-4">
                      <p className="text-sm font-medium text-slate-500">
                        議題 {selectedAgendaIndex + 1} / {agendas.length}
                      </p>

                      <div className="mt-2 flex items-start justify-between gap-6">
                        <h2 className="text-3xl font-bold text-slate-900">
                          {selectedAgenda.title}
                        </h2>

                        <span className="shrink-0 rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700">
                          {selectedAgenda.planned_minutes}分
                        </span>
                      </div>
                    </div>

                    <h3 className="mt-4 font-semibold text-slate-500">
                      目的：{selectedAgenda.purpose || ""}
                    </h3>

                    <div className="mt-5 grid gap-5 lg:grid-cols-2 items-stretch">
                      <section className="h-full rounded-2xl border border-slate-200 bg-slate-50 p-5">
                        <h3 className="mb-2 font-semibold text-slate-900">
                          議論ポイント
                        </h3>

                        <p className="whitespace-pre-wrap text-slate-700">
                          {selectedAgenda.discussion_points || ""}
                        </p>
                      </section>

                      <section className="h-full rounded-2xl border border-slate-200 bg-slate-50 p-5">
                        <h3 className="mb-2 font-semibold text-slate-900">
                          質問事項
                        </h3>

                        <p className="whitespace-pre-wrap text-slate-700">
                          {selectedAgenda.questions || ""}
                        </p>
                      </section>
                    </div>

                    <section className="mt-5 rounded-2xl border border-amber-100 bg-amber-50/70 p-5">
                      <h3 className="mb-2 font-semibold text-amber-800">
                        メモ
                      </h3>

                      <p className="whitespace-pre-wrap text-slate-700">
                        {selectedAgenda.memo || ""}
                      </p>
                    </section>
                  </>
                ) : (
                  <div className="flex min-h-80 items-center justify-center text-slate-500">
                    アジェンダがありません
                  </div>
                )}
              </main>

            </div>

          </div>

          {/* 参加メンバー */}
          <MeetingMemberList meetingId={meeting.id} />
        </>
      )}

    </div>
  );
}