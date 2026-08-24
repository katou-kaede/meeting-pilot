import CircularTimer from "../components/CircularTimer";
import {
  Play,
  Pause,
} from "lucide-react";

type Props = {
  plannedMinutes: number;
  plannedSeconds: number;
  elapsedSeconds: number;
  remainingSeconds: number;
  progress: number;
  saving: boolean;
  paused: boolean;
  formatTimer: (seconds: number) => string;
  onSave: () => void;
  onComplete: () => void;
  onPauseResume: () => void;
};

export default function MeetingTimerCard({
  plannedMinutes,
  plannedSeconds,
  elapsedSeconds,
  remainingSeconds,
  progress,
  saving,
  paused,
  formatTimer,
  onSave,
  onComplete,
  onPauseResume,
}: Props) {
  return (
    <div className="h-full rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex h-full items-center gap-6">
        <CircularTimer
          remainingSeconds={remainingSeconds}
          progress={progress}
          color="#3b82f6"
          label="全体残り"
          overtimeLabel="全体超過"
          formatTimer={formatTimer}
        />

        <div className="min-w-0 flex-1">
          <p className="text-sm text-slate-500">
            {plannedMinutes}分の全体タイマー
          </p>

          <p className="mt-1 font-medium text-slate-700">
            経過 {formatTimer(elapsedSeconds)} /{" "}
            {formatTimer(plannedSeconds)}
          </p>

          <div className="mt-4 flex flex-wrap gap-2">
            <button
              type="button"
              onClick={onSave}
              disabled={saving}
              className="cursor-pointer rounded-xl border border-slate-300 bg-white px-5 py-2 font-medium text-slate-700 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {saving ? "処理中..." : "一時保存"}
            </button>

            <button
              type="button"
              onClick={onComplete}
              disabled={saving}
              className="cursor-pointer rounded-xl bg-emerald-600 px-6 py-2 font-medium text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {saving ? "処理中..." : "会議終了"}
            </button>

            <button
              type="button"
              onClick={onPauseResume}
              disabled={saving}
              className="flex items-center cursor-pointer rounded-xl bg-blue-600 px-5 py-2 font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {paused ? (
                <><Play size={15} className="me-1"/>会議を再開</>
              ) : (
                <><Pause size={15} className="me-1"/>一時停止</>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}