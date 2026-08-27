import {
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  afterEach,
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from "vitest";

import MeetingMemberManager from "./MeetingMemberManager";
import { useAuth } from "../contexts/AuthContext";
import { MemoryRouter } from "react-router-dom";

vi.mock("../contexts/AuthContext", () => ({
  useAuth: vi.fn(),
}));

const mockUseAuth = vi.mocked(useAuth);

const members = [
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
    role: "editor",
    created_at: "2026-08-27T00:00:00Z",
  },
];

beforeEach(() => {
  mockUseAuth.mockReturnValue({
    user: {
      id: 1,
      name: "Employee1",
      email: "employee1@example.com",
      is_active: true,
      created_at: "2026-08-27T00:00:00Z",
      updated_at: "2026-08-27T00:00:00Z",
    },
    loading: false,
    login: vi.fn(),
    logout: vi.fn(),
    fetchCurrentUser: vi.fn(),
  });
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("MeetingMemberManager", () => {
  test("主催者には削除ボタンを表示しない", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(
      new Response(
        JSON.stringify(members),
        {
          status: 200,
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
    );

    renderMeetingMemberManager();

    expect(
      await screen.findByText("Employee1")
    ).toBeInTheDocument();

    expect(
      screen.queryByRole("button", {
        name: "Employee1を参加メンバーから削除",
      })
    ).not.toBeInTheDocument();

    expect(
      screen.getByRole("button", {
        name: "Employee2を参加メンバーから削除",
      })
    ).toBeInTheDocument();
  });

  
  test("編集者が設定済みなら新しい編集者を選択できない", async () => {
    const user = userEvent.setup();

    const fetchMock = vi.spyOn(
      globalThis,
      "fetch"
    );

    // 1回目：登録済みメンバー取得
    fetchMock.mockResolvedValueOnce(
      new Response(
        JSON.stringify(members),
        {
          status: 200,
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
    );

    // 2回目：追加候補検索
    fetchMock.mockResolvedValueOnce(
      new Response(
        JSON.stringify([
          {
            id: 3,
            name: "Employee3",
            email: "employee3@example.com",
            is_active: true,
            created_at: "2026-08-27T00:00:00Z",
            updated_at: "2026-08-27T00:00:00Z",
          },
        ]),
        {
          status: 200,
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
    );

    renderMeetingMemberManager();

    await screen.findByText("Employee1");

    await user.click(
      screen.getByRole("button", {
        name: "検索",
      })
    );

    await user.click(
      await screen.findByRole("button", {
        name: /Employee3/,
      })
    );

    const editorOption = screen.getByRole(
      "option",
      {
        name: "編集者（設定済み）",
      }
    );

    expect(editorOption).toBeDisabled();

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(2);
    });
  });
});

function renderMeetingMemberManager() {
  render(
    <MemoryRouter>
      <MeetingMemberManager meetingId={1} />
    </MemoryRouter>
  );
}
