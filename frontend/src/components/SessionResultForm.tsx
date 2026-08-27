type Props = {
  decisions: string;
  todo: string;
  canEditSession: boolean;
  onDecisionsChange: (value: string) => void;
  onTodoChange: (value: string) => void;
};

export default function SessionResultForm({
  decisions,
  todo,
  canEditSession,
  onDecisionsChange,
  onTodoChange,
}: Props) {
  return (
    <div className="mt-6 grid gap-6 md:grid-cols-2">
      <div className="rounded-2xl border border-blue-100 bg-blue-50/70 p-6 shadow-sm backdrop-blur-xl">
        <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-blue-600">
          Meeting Result
        </p>

        <h2 className="mb-4 text-lg font-semibold text-slate-900">
          決定事項
        </h2>

        {canEditSession ? (
          <textarea
            aria-label="決定事項"
            value={decisions}
            rows={6}
            onChange={(event) =>
              onDecisionsChange(event.target.value)
            }
            className="w-full rounded-xl border border-slate-300 px-3 py-2"
          />
        ) : (
          <p className="whitespace-pre-wrap text-slate-700">
            {decisions || "未入力"}
          </p>
        )}
      </div>

      <div className="rounded-2xl border border-amber-100 bg-amber-50/70 p-6 shadow-sm backdrop-blur-xl">
        <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-amber-600">
          Next Action
        </p>

        <h2 className="mb-4 text-lg font-semibold text-slate-900">
          TODO
        </h2>

        {canEditSession ? (
          <textarea
            aria-label="TODO"
            value={todo}
            rows={6}
            onChange={(event) =>
              onTodoChange(event.target.value)
            }
            className="w-full rounded-xl border border-slate-300 px-3 py-2"
          />
        ) : (
          <p className="whitespace-pre-wrap text-slate-700">
            {todo || "未入力"}
          </p>
        )}
      </div>
    </div>
  );
}