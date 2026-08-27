export type User = {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

// ============================================
// 会議参加メンバー
// ============================================
export type MeetingMember = {
  meeting_id: number;
  user_id: number;
  name: string;
  email: string;
  role: "owner" | "editor" | "viewer";
  created_at: string;
};

export type MemberCandidate = {
  id: number;
  name: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};