// ステータス表示用
export const getStatusLabel = (status: string) => {
  switch (status) {
    case "scheduled":
      return "予定";
    case "completed":
      return "完了";
    case "in_progress":
      return "進行中";
    default:
      return status;
  }
};

export const getStatusStyle = (status: string) => {
  switch (status) {
    case "scheduled":
      return "bg-slate-100 text-slate-700";

    case "completed":
      return "bg-emerald-100 text-emerald-700";

    case "in_progress":
      return "bg-blue-100 text-blue-700";

    default:
      return "bg-slate-100 text-slate-700";
  }
};