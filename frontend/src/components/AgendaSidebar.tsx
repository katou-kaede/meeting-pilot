type Agenda = {
  id?: number;
  title: string;
  planned_minutes: number | string;
};

type Props = {
  agendas: Agenda[];
  selectedIndex: number;
  currentIndex?: number;
  meetingPlannedMinutes?: number | string;
  onSelect: (index: number) => void;
  onMove? : (
    index: number,
    direction: "up" | "down"
  ) => void;
};

export default function AgendaSidebar({
  agendas,
  selectedIndex,
  currentIndex,
  meetingPlannedMinutes,
  onSelect,
  onMove,
}: Props) {
  // Agendaの合計時間の計算
  const totalMinutes = agendas.reduce(
    (sum, agenda) => sum + Number(agenda.planned_minutes),
    0
  );

  const isAgendaOver =
    meetingPlannedMinutes !== undefined &&
    totalMinutes > Number(meetingPlannedMinutes);

  return (
    <aside className="self-start rounded-3xl border border-white/70 bg-white/75 p-5 shadow-sm backdrop-blur-xl">
      <div className="mb-5 flex items-center justify-between">
        <h2 className="font-bold text-slate-900">
          アジェンダ
        </h2>

        <span className="rounded-lg bg-slate-100 px-2 py-1 text-xs text-slate-600">
          {agendas.length}件
        </span>
      </div>

      <div className="space-y-2">
        {agendas.map((agenda, index) => {
          const isSelected = selectedIndex === index;
          const isCompleted =
            currentIndex !== undefined && index < currentIndex;
          const isCurrent =
            currentIndex !== undefined && index === currentIndex;
        
          return(
          <div
            key={agenda.id ?? index}
            className="group relative"
          >
            <button
              type="button"
              onClick={() => onSelect(index)}
              className={`w-full rounded-2xl p-4 pr-12 text-left transition ${
                isSelected
                  ? "bg-slate-900 text-white shadow-md"
                  : isCompleted
                    ? "bg-emerald-50 text-emerald-800"
                    : isCurrent
                      ? "bg-blue-50 text-blue-800 ring-1 ring-blue-200"
                      : "text-slate-700 hover:bg-white hover:shadow-sm"
              }`}
            >
              <div className="flex gap-3">
                <span
                  className={`mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs ${
                    isSelected
                      ? "bg-white/15 text-white"
                      : isCompleted
                        ? "bg-emerald-100 text-emerald-700"
                        : isCurrent
                          ? "bg-blue-100 text-blue-700"
                          : "bg-slate-100 text-slate-500"
                  }`}
                >
                  {isCompleted ? "✓" : index + 1}
                </span>

                <div className="min-w-0">
                  <div className="truncate font-medium">
                    {agenda.title || `議題${index + 1}`}
                  </div>

                  <div
                    className={`mt-1 text-xs ${
                      isSelected
                        ? "text-slate-300"
                        : "text-slate-500"
                    }`}
                  >
                    予定 {agenda.planned_minutes}分
                  </div>
                </div>
              </div>
            </button>

            {onMove && (
              <div
                className={`absolute right-2 top-1/2 flex -translate-y-1/2 flex-col text-xs transition-opacity ${
                  isSelected
                    ? "opacity-100"
                    : "opacity-0 group-hover:opacity-100"
                }`}
              >
                <button
                  type="button"
                  disabled={index === 0}
                  onClick={() => onMove(index, "up")}
                  className={`rounded px-1.5 py-0.5 disabled:cursor-not-allowed disabled:opacity-20 ${
                    isSelected
                      ? "text-slate-300 hover:bg-white/10 hover:text-white"
                      : "text-slate-400 hover:bg-slate-100 hover:text-slate-700"
                  }`}
                  aria-label="議題を上へ移動"
                  title="上へ移動"
                >
                  ▲
                </button>

                <button
                  type="button"
                  disabled={index === agendas.length - 1}
                  onClick={() => onMove(index, "down")}
                  className={`rounded px-1.5 py-0.5 disabled:cursor-not-allowed disabled:opacity-20 ${
                    isSelected
                      ? "text-slate-300 hover:bg-white/10 hover:text-white"
                      : "text-slate-400 hover:bg-slate-100 hover:text-slate-700"
                  }`}
                  aria-label="議題を下へ移動"
                  title="下へ移動"
                >
                  ▼
                </button>
              </div>
            )}
          </div>
          );
        })}
        
      </div>

      <div className="mt-6 flex justify-between border-t border-slate-200 pt-4 text-sm font-medium">
        <span className="text-slate-500">
          合計時間
        </span>

        <span>{totalMinutes}分</span>
      </div>

      {isAgendaOver && (
        <div className="mt-2 text-sm text-red-600">
          会議の予定時間を超えています
        </div>
      )}
    </aside>
  );
}