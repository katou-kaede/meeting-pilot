import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  describe,
  expect,
  test,
  vi,
} from "vitest";

import SessionEditorSelector from "./SessionEditorSelector";
import type { MeetingMember } from "../types/user";

const members: MeetingMember[] = [
  {
    meeting_id: 1,
    user_id: 1,
    name: "Employee1",
    email: "employee1@example.com",
    role: "owner",
    created_at: "2026-08-27T00:00:00Z",
  },
  {
    meeting_id: 1,
    user_id: 2,
    name: "Employee2",
    email: "employee2@example.com",
    role: "viewer",
    created_at: "2026-08-27T00:00:00Z",
  },
  {
    meeting_id: 1,
    user_id: 3,
    name: "Employee3",
    email: "employee3@example.com",
    role: "editor",
    created_at: "2026-08-27T00:00:00Z",
  },
];

describe("SessionEditorSelector", () => {
  test("ownerには編集者プルダウンを表示する", () => {
    render(
      <SessionEditorSelector
        currentUserRole="owner"
        editorUserId={null}
        members={members}
        onChangeEditor={vi.fn()}
      />
    );

    expect(
      screen.getByRole("combobox", {
        name: "編集者",
      })
    ).toBeInTheDocument();

    expect(
      screen.getByRole("option", {
        name: "Employee1（主催者）",
      })
    ).toBeInTheDocument();
  });

  test.each(["editor", "viewer"] as const)(
    "%sには編集者プルダウンを表示しない",
    (role) => {
      render(
        <SessionEditorSelector
          currentUserRole={role}
          editorUserId={null}
          members={members}
          onChangeEditor={vi.fn()}
        />
      );

      expect(
        screen.queryByRole("combobox", {
          name: "編集者",
        })
      ).not.toBeInTheDocument();
    }
  );

  test("参加者を選択するとユーザーIDを渡す", async () => {
    const user = userEvent.setup();
    const onChangeEditor = vi.fn();

    render(
      <SessionEditorSelector
        currentUserRole="owner"
        editorUserId={null}
        members={members}
        onChangeEditor={onChangeEditor}
      />
    );

    await user.selectOptions(
      screen.getByRole("combobox", {
        name: "編集者",
      }),
      "2"
    );

    expect(onChangeEditor).toHaveBeenCalledWith(2);
  });

  test("主催者を選択するとnullを渡す", async () => {
    const user = userEvent.setup();
    const onChangeEditor = vi.fn();

    render(
      <SessionEditorSelector
        currentUserRole="owner"
        editorUserId={3}
        members={members}
        onChangeEditor={onChangeEditor}
      />
    );

    await user.selectOptions(
      screen.getByRole("combobox", {
        name: "編集者",
      }),
      ""
    );

    expect(onChangeEditor).toHaveBeenCalledWith(null);
  });
});