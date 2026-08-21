type Props = {
  remainingSeconds: number;
  progress: number;
  color: string;
  label: string;
  overtimeLabel: string;
  formatTimer: (seconds: number) => string;
};

export default function CircularTimer({
  remainingSeconds,
  progress,
  color,
  label,
  overtimeLabel,
  formatTimer,
}: Props) {
  const isOvertime = remainingSeconds < 0;

  return (
    <div
      className="relative flex h-28 w-28 shrink-0 items-center justify-center rounded-full"
      style={{
        background: isOvertime
          ? "conic-gradient(#ef4444 100%, #e2e8f0 0)"
          : `conic-gradient(
              ${color} ${progress}%,
              #e2e8f0 ${progress}%
            )`,
      }}
    >
      <div className="absolute inset-2 flex items-center justify-center rounded-full bg-white">
        <div className="text-center">
          <div
            className={`text-xs ${
              isOvertime ? "text-red-600" : "text-slate-500"
            }`}
          >
            {isOvertime ? overtimeLabel : label}
          </div>

          <div
            className={`text-2xl font-bold ${
              isOvertime ? "text-red-600" : "text-slate-900"
            }`}
          >
            {formatTimer(remainingSeconds)}
          </div>
        </div>
      </div>
    </div>
  );
}