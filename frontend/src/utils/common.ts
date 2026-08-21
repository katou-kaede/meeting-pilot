export function formatDateTime(
  value: string | null
): string {
  if (!value) {
    return "-";
  }

  return new Date(value).toLocaleString(
    "ja-JP",
    {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    }
  );
}