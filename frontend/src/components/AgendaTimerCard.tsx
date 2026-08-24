import CircularTimer from "../components/CircularTimer";

type Props = {
  title: string;
  plannedSeconds: number;
  elapsedSeconds: number;
  remainingSeconds: number;
  progress: number;
  currentIndex: number;
  agendaCount: number;
  formatTimer: (seconds: number) => string;
  onPrevious: () => void;
  onNext: () => void;
};

export default function AgendaTimerCard({
  title,
  plannedSeconds,
  elapsedSeconds,
  remainingSeconds,
  progress,
  currentIndex,
  agendaCount,
  formatTimer,
  onPrevious,
  onNext,
}: Props) {
  return (
    <div className="h-full rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex h-full items-center gap-6">
        <CircularTimer
          remainingSeconds={remainingSeconds}
          progress={progress}
          color="#8b5cf6"
          label="議題残り"
          overtimeLabel="議題超過"
          formatTimer={formatTimer}
        />

        <div className="min-w-0 flex-1">
          <p className="text-sm text-slate-500">
            現在の議題
          </p>

          <h2 className="mt-1 truncate text-2xl font-bold text-slate-900">
            {title || "議題なし"}
          </h2>

          <p className="mt-2 text-sm text-slate-500">
            経過 {formatTimer(elapsedSeconds)} /{" "}
            {formatTimer(plannedSeconds)}
          </p>

          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              onClick={onPrevious}
              disabled={currentIndex === 0}
              className="cursor-pointer rounded-xl border border-slate-300 px-4 py-2 text-slate-700 disabled:cursor-not-allowed disabled:opacity-40"
            >
              ← 議題を戻す
            </button>

            <button
              type="button"
              onClick={onNext}
              disabled={
                agendaCount === 0 ||
                currentIndex >= agendaCount - 1
              }
              className="cursor-pointer rounded-xl bg-slate-900 px-4 py-2 text-white disabled:cursor-not-allowed disabled:opacity-40"
            >
              議題を進める →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}