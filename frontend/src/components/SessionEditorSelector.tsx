import type { MeetingMember } from "../types/user";

type Props = {
  currentUserRole: "owner" | "editor" | "viewer";
  editorUserId: number | null;
  members: MeetingMember[];
  onChangeEditor: (userId: number | null) => void;
};

export default function SessionEditorSelector({
  currentUserRole,
  editorUserId,
  members,
  onChangeEditor,
}: Props) {
  if (currentUserRole !== "owner") {
    return null;
  }

  const owner = members.find(
    (member) => member.role === "owner"
  );

  return (
    <div>
      <label
        htmlFor="session-editor"
        className="mb-1 block text-sm font-medium text-slate-700"
      >
        編集者
      </label>

      <select
        id="session-editor"
        value={editorUserId ?? ""}
        onChange={(event) => {
          const userId = event.target.value
            ? Number(event.target.value)
            : null;

          onChangeEditor(userId);
        }}
        className="rounded-xl border border-slate-300 bg-white px-3 py-2"
      >
        <option value="">
          {owner
            ? `${owner.name}（主催者）`
            : "主催者"}
        </option>

        {members
          .filter((member) => member.role !== "owner")
          .map((member) => (
            <option
              key={member.user_id}
              value={member.user_id}
            >
              {member.name}
            </option>
          ))}
      </select>
    </div>
  );
}