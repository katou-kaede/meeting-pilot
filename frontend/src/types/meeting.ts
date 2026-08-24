export type Meeting = {
  id: number;
  title: string;
  target_name: string;
  scheduled_start_at: string | null;
  planned_minutes: number;
  decisions: string;
  todo: string;
  status: string;
};

export type Agenda = {
  id: number;
  title: string;
  purpose: string;
  discussion_points: string;
  questions: string;
  memo: string;
  planned_minutes: number;
  sort_order: number;
};

// ============================================
// フォーム用
// ============================================
export type AgendaForm = {
  id?: number;
  title: string;
  purpose: string;
  discussion_points: string;
  questions: string;
  memo: string;
  planned_minutes: string;
};

// ============================================
// ミーティング詳細
// ============================================
export type MeetingDetail = {
  id: number;
  title: string;
  description: string;
  target_name: string;
  scheduled_start_at: string | null;
  actual_start_at: string | null;
  planned_minutes: number;
  decisions: string;
  todo: string;
  status: string;
  agendas: Agenda[];
};

// ============================================
// 会議中
// ============================================
export type AgendaSessionDetail = {
  id: number;
  title: string;
  purpose: string;
  discussion_points: string;
  questions: string;
  memo: string;
  planned_minutes: number;
  sort_order: number;
  actual_start_at: string | null;
  actual_end_at: string | null;
  elapsed_seconds: number;
};

export type MeetingSessionDetail = {
  id: number;
  title: string;
  description: string;
  planned_minutes: number;
  status: string;
  actual_start_at: string | null;
  paused_at: string | null;
  total_paused_seconds: number;
  current_agenda_id: number | null;
  decisions: string;
  todo: string;
  agendas: AgendaSessionDetail[];
};