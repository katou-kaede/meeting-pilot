import { Link, useNavigate, useParams } from "react-router-dom";
import { useCallback, useEffect, useState } from "react";
import type { MeetingSessionDetail } from "../types/meeting";
import type { MeetingMember } from "../types/user";
import ErrorMessage from "../components/ErrorMessage";
import Loading from "../components/Loading";
import AgendaSidebar from "../components/AgendaSidebar";
import AgendaTimerCard from "../components/AgendaTimerCard";
import MeetingTimerCard from "../components/MeetingTimerCard";
import SessionResultForm from "../components/SessionResultForm";
import SessionEditorSelector from "../components/SessionEditorSelector";
import { ArrowLeft } from "lucide-react";
import { useMeetingSessionWebSocket } from "../hooks/useMeetingSessionWebSocket";

export default function MeetingSessionPage() {
  // ============================================
  // State・画面データ
  // ============================================
  const { id } = useParams();

  const [meeting, setMeeting] = useState<MeetingSessionDetail | null>(null);

  const [selectedAgendaIndex, setSelectedAgendaIndex] = useState(0);
  const agendas = meeting?.agendas ?? [];

  const [currentAgendaIndex, setCurrentAgendaIndex] = useState(0);

  const selectedAgenda = agendas[selectedAgendaIndex] ?? null;
  const currentAgenda = agendas[currentAgendaIndex] ?? null;

  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState("");
  const [connectionError, setConnectionError] = useState("");

  // 通信状態
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState("");

  const navigate = useNavigate();

  // タイマー
  const [agendaRemainingSeconds, setAgendaRemainingSeconds] = useState<number[]>([]);
  const [currentTime, setCurrentTime] = useState(() => Date.now());

  // 参加メンバー
  const [members, setMembers] = useState<MeetingMember[]>([]);

  // ============================================
  // 会議情報の取得
  // ============================================
  const fetchMeeting = useCallback(async (syncSelectedAgenda = false) => {
    setErrorMessage("");
    try {

      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/session`, 
        {
          credentials: "include",
        }
      );

      const data = await response.json();

      if (!response.ok) {
        setErrorMessage(data.error);
        return;
      }

      const meetingData = data as MeetingSessionDetail;

      if (meetingData.status !== "in_progress") {
        navigate(`/meetings/${id}`, {
          replace: true,  // 進行中でない場合は、セッションページから詳細ページへリダイレクト
        });

        return;
      }

      setMeeting(meetingData);

      // DBの現在議題からIndexを復元
      const currentIndex = meetingData.agendas.findIndex(
        (agenda) => agenda.id === meetingData.current_agenda_id
      );

      const safeCurrentIndex = currentIndex >= 0 ? currentIndex : 0;

      setCurrentAgendaIndex(safeCurrentIndex);

      // 初回表示時だけ、選択中議題も進行中議題へ合わせる
      if (syncSelectedAgenda) {
        setSelectedAgendaIndex(safeCurrentIndex);
      }

      // 議題ごとの残り時間を初期化
      const now = Date.now();

      const initialRemainingSeconds = meetingData.agendas.map((agenda) => {
        let elapsedSeconds = agenda.elapsed_seconds;

        // 現在計測中の場合だけ、開始後の時間を加算
        if (
          agenda.id === meetingData.current_agenda_id &&
          agenda.actual_start_at &&
          agenda.actual_end_at === null &&
          meetingData.paused_at === null
        ) {
          elapsedSeconds += Math.floor(
            (now - new Date(agenda.actual_start_at).getTime()) / 1000
          );
        }

        return agenda.planned_minutes * 60 - elapsedSeconds;
      });

      setAgendaRemainingSeconds(initialRemainingSeconds);

    } catch (error) {
      console.error(error);
      setErrorMessage("会議情報の取得に失敗しました");

    } finally {
      setLoading(false);
    }
  }, [id, navigate]);

  // 初回データ取得
  useEffect(() => {
    fetchMeeting(true);
  }, [fetchMeeting]);

  // ============================================
  // 参加メンバー取得
  // ============================================
  const fetchMembers = useCallback(async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/members`,
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
    }
  }, [id]);

  useEffect(() => {
    fetchMembers();
  }, [fetchMembers]);

  // ============================================
  // WebSocket接続用
  // ============================================
  const handleSessionUpdated = useCallback(() => {
    fetchMeeting(false);
  }, [fetchMeeting]);

  const handleMeetingCompleted = useCallback(() => {
    navigate(`/meetings/${id}`);
  }, [navigate, id]);

  const handleConnectionError = useCallback(
    (message: string) => {
      setConnectionError(message);
    },
    []
  );

  useMeetingSessionWebSocket({
    meetingId: id,
    onSessionUpdated: handleSessionUpdated,
    onMeetingCompleted: handleMeetingCompleted,
    onConnectionError: handleConnectionError,
  });

  // ============================================
  // 「一時保存しました」を3秒後に消す
  // ============================================
  useEffect(() => {
    if (!saveMessage) return;

    // 「一時保存しました」を3秒後に消すためのタイマー
    const timerId = window.setTimeout(() => {
      setSaveMessage("");
    }, 3000);

    return () => {
      window.clearTimeout(timerId);
    };
  }, [saveMessage]);

  // ============================================
  // 現在進行中の議題の残り時間を減らす
  // ============================================
  useEffect(() => {
    if (!meeting || meeting.paused_at) return;

    const timerId = window.setInterval(() => {
      setAgendaRemainingSeconds((prev) =>
        prev.map((seconds, index) =>
          index === currentAgendaIndex
            ? seconds - 1
            : seconds
        )
      );
    }, 1000);

    return () => window.clearInterval(timerId);
  }, [currentAgendaIndex, meeting?.paused_at]);

  // ============================================
  // 1秒ごとに現在時刻を更新
  // ============================================
  useEffect(() => {
    const timerId = window.setInterval(() => {
      setCurrentTime(Date.now());
    }, 1000);

    return () => {
      window.clearInterval(timerId);
    };
  }, []);


  // ============================================
  // 会議終了処理
  // ============================================
  const handleComplete = async () => {
    if (!meeting) return;

    if (!confirm("会議を終了しますか？")) {
      return;
    }

    try {
      setErrorMessage("");
      setSaving(true);

      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/complete`,
        {
          method: "PATCH",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            decisions: meeting.decisions,
            todo: meeting.todo,
            agendas: (meeting.agendas ?? []).map((agenda) => ({
              id: agenda.id,
              memo: agenda.memo,
            })),
          }),
        }
      );

      if (!response.ok) {
        const error = await response.json();
        setErrorMessage(error.error);
        return;
      }

      navigate(`/meetings/${id}`);

    } catch (error) {
      setErrorMessage("会議の終了に失敗しました");
      console.error(error);

    } finally {
      setSaving(false);
    }
  };

  // ============================================
  // 一時保存処理
  // ============================================
  const handleSaveSession = async () => {
    if (!meeting || saving) return;

    setErrorMessage("");
    setSaveMessage("");
    setSaving(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/session`,
        {
          method: "PATCH",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            decisions: meeting.decisions,
            todo: meeting.todo,
            agendas: meeting.agendas.map((agenda) => ({
              id: agenda.id,
              memo: agenda.memo,
            })),
          }),
        }
      );

      if (!response.ok) {
        const error = await response.json();
        setErrorMessage(error.error || "会議内容の一時保存に失敗しました");
        return;
      }

      setSaveMessage("一時保存しました");

    } catch (error) {
      console.error(error);
      setErrorMessage("一時保存に失敗しました");

    } finally {
      setSaving(false);
    }
  };

  // ============================================
  // 議題の切り替え処理
  // ============================================
  const changeCurrentAgenda = async (targetIndex: number) => {
    const targetAgenda = agendas[targetIndex];

    if (!targetAgenda) return;

    try {
      setErrorMessage("");

      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/current-agenda`,
        {
          method: "PATCH",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            agenda_id: targetAgenda.id,
          }),
        }
      );

      if (!response.ok) {
        const error = await response.json();
        setErrorMessage(error.error);
        return;
      }

      setCurrentAgendaIndex(targetIndex);
      setSelectedAgendaIndex(targetIndex);

      // フロント側の開始時刻も更新
      setMeeting((prev) => {
        if (!prev) return prev;

        return {
          ...prev,
          current_agenda_id: targetAgenda.id,
          agendas: prev.agendas.map((agenda) =>
            agenda.id === targetAgenda.id
              ? {
                ...agenda,
                actual_start_at: new Date().toISOString(),
                actual_end_at: null,
              }
              : agenda
          ),
        };
      });
    } catch (error) {
      console.error(error);
      setErrorMessage("議題の切り替えに失敗しました");
    }
  };

  // ============================================
  // 一時停止・再開処理
  // ============================================
  const handlePauseResume = async () => {
    if (!meeting || saving) return;

    const isPaused = meeting.paused_at !== null;
    const action = isPaused ? "resume" : "pause";

    setErrorMessage("");
    setSaving(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/meetings/${id}/${action}`,
        {
          method: "PATCH",
          credentials: "include",
        }
      );

      if (!response.ok) {
        const error = await response.json();

        setErrorMessage(
          error.error ||
          (isPaused
            ? "会議の再開に失敗しました"
            : "会議の一時停止に失敗しました")
        );
        return;
      }

      // 最新の停止状態・時間を取得
      await fetchMeeting(false);

    } catch (error) {
      console.error(error);
      setErrorMessage(
        isPaused
          ? "会議の再開に失敗しました"
          : "会議の一時停止に失敗しました"
      );

    } finally {
      setSaving(false);
    }
  };

  // ============================================
  // 編集者の変更
  // ============================================
  const handleChangeEditor = async (
    userId: number | null
  ) => {
    const response = await fetch(
      `http://localhost:8080/api/meetings/${id}/editor`,
      {
        method: "PATCH",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user_id: userId,
        }),
      }
    );

    if (!response.ok) {
      const data = await response.json();
      setErrorMessage(data.error);
      return;
    }

    await fetchMeeting(false);
  };

  // ============================================
  // 表示用ユーティリティ
  // ============================================
  // タイマーの表示形式
  const formatTimer = (seconds: number) => {
    const absoluteSeconds = Math.abs(seconds);
    const minutes = Math.floor(absoluteSeconds / 60);
    const remainingSeconds = absoluteSeconds % 60;

    return `${String(minutes).padStart(2, "0")}:${String(
      remainingSeconds
    ).padStart(2, "0")}`;
  };

  // 議題を進める処理
  const handleNextAgenda = async () => {
    if (currentAgendaIndex >= agendas.length - 1) {
      return;
    }

    await changeCurrentAgenda(currentAgendaIndex + 1);
  };

  // 議題を戻す処理
  const handlePreviousAgenda = async () => {
    if (currentAgendaIndex <= 0) {
      return;
    }

    await changeCurrentAgenda(currentAgendaIndex - 1);
  };

  // 会議全体の予定時間
  const meetingPlannedSeconds = (meeting?.planned_minutes ?? 0) * 60;

  // 会議全体の経過時間計算用の基準時刻(paused_atがあればそこ、なければ現在時刻)
  const meetingTimerEndTime = meeting?.paused_at
    ? new Date(meeting.paused_at).getTime()
    : currentTime;

  // 会議全体の経過時間計算
  const meetingElapsedSeconds = meeting?.actual_start_at
    ? Math.max(
      0,
      Math.floor(
        (
          meetingTimerEndTime -
          new Date(meeting.actual_start_at).getTime()
        ) / 1000
      ) - meeting.total_paused_seconds
    )
    : 0;

  // 会議全体の残り時間計算
  const meetingRemainingSeconds = meetingPlannedSeconds - meetingElapsedSeconds;

  // 現在議題の残り時間
  const currentAgendaRemainingSeconds = agendaRemainingSeconds[currentAgendaIndex] ?? 0;

  // 現在議題の経過時間（秒）
  const currentAgendaElapsedSeconds =
    (currentAgenda?.planned_minutes ?? 0) * 60 -
    currentAgendaRemainingSeconds;

  // 現在議題の予定時間
  const currentAgendaPlannedSeconds = (currentAgenda?.planned_minutes ?? 0) * 60;

  // 会議全体時間の進捗率（タイマー表示用）
  const meetingProgress =
    meetingPlannedSeconds > 0
      ? Math.min(
        100,
        Math.max(
          0,
          (meetingElapsedSeconds / meetingPlannedSeconds) * 100
        )
      )
      : 0;

  // 現在議題時間の進捗率（タイマー表示用）
  const currentAgendaProgress =
    currentAgendaPlannedSeconds > 0
      ? Math.min(
        100,
        Math.max(
          0,
          (currentAgendaElapsedSeconds /
            currentAgendaPlannedSeconds) *
          100
        )
      )
      : 0;

  if (loading) {
    return <Loading />;
  }


  return (
    <div className="mx-auto max-w-[1600px] p-6 lg:p-8">

      <Link
        to={`/meetings/${id}`}
        className="inline-flex items-center text-slate-500 hover:text-slate-900"
      >
        <ArrowLeft size={18} /> 会議詳細へ戻る
      </Link>

      {/* エラーメッセージ */}
      {(errorMessage || connectionError) && (
        <ErrorMessage message={errorMessage || connectionError} />
      )}

      {/* 保存メッセージ */}
      {saveMessage && (
        <div className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-emerald-700">
          {saveMessage}
        </div>
      )}

      {meeting && (
        <>
          {/* タイトルエリア */}
          <div className="mt-6 flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-slate-900">
                {meeting.title}
              </h1>

              <p className="mt-2 text-slate-500">
                {meeting.description}
              </p>
            </div>

            {/* 編集者の変更プルダウン */}
            <SessionEditorSelector
              currentUserRole={meeting.current_user_role}
              editorUserId={meeting.editor_user_id}
              members={members}
              onChangeEditor={handleChangeEditor}
            />
          </div>

          {/* タイマー */}
          <div className="mt-6 grid gap-6 lg:grid-cols-2">

            {/* 全体タイマー */}
            <MeetingTimerCard
              plannedMinutes={meeting.planned_minutes}
              plannedSeconds={meetingPlannedSeconds}
              elapsedSeconds={meetingElapsedSeconds}
              remainingSeconds={meetingRemainingSeconds}
              progress={meetingProgress}
              saving={saving}
              paused={meeting.paused_at !== null}
              canEditSession={meeting.can_edit_session}
              formatTimer={formatTimer}
              onSave={handleSaveSession}
              onComplete={handleComplete}
              onPauseResume={handlePauseResume}
            />

            {/* 議題タイマー */}
            <AgendaTimerCard
              title={currentAgenda?.title ?? ""}
              plannedSeconds={currentAgendaPlannedSeconds}
              elapsedSeconds={currentAgendaElapsedSeconds}
              remainingSeconds={currentAgendaRemainingSeconds}
              progress={currentAgendaProgress}
              currentIndex={currentAgendaIndex}
              agendaCount={agendas.length}
              formatTimer={formatTimer}
              onPrevious={handlePreviousAgenda}
              onNext={handleNextAgenda}
            />
          </div>

          {/* アジェンダ */}
          <div className="mt-6 grid gap-6 lg:grid-cols-[320px_minmax(0,1fr)]">

            <AgendaSidebar
              agendas={agendas}
              selectedIndex={selectedAgendaIndex}
              currentIndex={currentAgendaIndex}
              onSelect={setSelectedAgendaIndex}
            />

            <main className="rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
              {selectedAgenda ? (
                <>
                  <div className="border-b border-slate-200 pb-4">
                    <p className="text-sm font-medium text-slate-500">
                      議題 {selectedAgendaIndex + 1} / {agendas.length}
                    </p>

                    <div className="mt-2 flex items-start justify-between gap-6">
                      <h2 className="text-2xl font-bold text-slate-900">
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
                        トピック
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

                    {meeting.can_edit_session ? (
                    <textarea
                      value={selectedAgenda.memo || ""}
                      rows={6}
                      onChange={(e) => {
                        const updated = [...agendas];
                        updated[selectedAgendaIndex] = {
                          ...updated[selectedAgendaIndex],
                          memo: e.target.value,
                        };

                        setMeeting({ ...meeting, agendas: updated });
                      }}
                      className="w-full rounded-xl border border-slate-300 px-3 py-2"
                    />
                    ) : (
                      <p className="whitespace-pre-wrap text-slate-700">
                        {selectedAgenda.memo || ""}
                      </p>
                    )}
                  </section>
                </>
              ) : (
                <div className="flex min-h-80 items-center justify-center text-slate-500">
                  アジェンダがありません
                </div>
              )}
            </main>
          </div>

          {/* 会議後入力項目 */}
          <SessionResultForm
            decisions={meeting.decisions || ""}
            todo={meeting.todo || ""}
            canEditSession={meeting.can_edit_session}
            onDecisionsChange={(value) => {
              setMeeting({...meeting, decisions: value});
            }}
            onTodoChange={(value) => {
              setMeeting({...meeting, todo: value});
            }}
          />

        </>
      )}

    </div>
  );
}