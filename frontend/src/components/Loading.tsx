type Props = {
  message?: string;
};

export default function Loading({
  message = "読み込み中...",
}: Props) {
  return (
    <div className="flex min-h-80 items-center justify-center">
      <div className="text-center">
        <div className="mx-auto h-8 w-8 animate-spin rounded-full border-4 border-slate-200 border-t-slate-700" />

        <p className="mt-4 text-sm text-slate-500">
          {message}
        </p>
      </div>
    </div>
  );
}