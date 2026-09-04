import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  MemoryRouter,
  Route,
  Routes,
} from "react-router-dom";
import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from "vitest";

import LoginPage from "./LoginPage";
import { useAuth } from "../contexts/AuthContext";

vi.mock("../contexts/AuthContext", () => ({
  useAuth: vi.fn(),
}));

vi.mock("../components/Loading", () => ({
  default: () => <div>読み込み中</div>,
}));

const mockUseAuth = vi.mocked(useAuth);

function renderLoginPage() {
  render(
    <MemoryRouter initialEntries={["/login"]}>
      <Routes>
        <Route
          path="/login"
          element={<LoginPage />}
        />

        <Route
          path="/"
          element={<h1>会議一覧</h1>}
        />
      </Routes>
    </MemoryRouter>
  );
}

describe("LoginPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("入力内容でログイン処理を実行する", async () => {
    const user = userEvent.setup();
    const loginMock = vi.fn().mockResolvedValue(undefined);

    mockUseAuth.mockReturnValue({
      user: null,
      loading: false,
      login: loginMock,
      logout: vi.fn(),
      fetchCurrentUser: vi.fn(),
      deactivateAccount: vi.fn(),
    });

    renderLoginPage();

    await user.type(
      screen.getByLabelText("メールアドレス"),
      "employee1@example.com"
    );

    await user.type(
      screen.getByLabelText("パスワード"),
      "password123"
    );

    await user.click(
      screen.getByRole("button", {
        name: "ログイン",
      })
    );

    expect(loginMock).toHaveBeenCalledWith(
      "employee1@example.com",
      "password123"
    );
  });


    test("ログイン失敗時にエラーメッセージを表示する", async () => {
        const user = userEvent.setup();

        const loginMock = vi.fn().mockRejectedValue(
        new Error(
            "メールアドレスまたはパスワードが正しくありません"
        )
        );

        mockUseAuth.mockReturnValue({
          user: null,
          loading: false,
          login: loginMock,
          logout: vi.fn(),
          fetchCurrentUser: vi.fn(),
          deactivateAccount: vi.fn(),
        });

        renderLoginPage();

        await user.type(
        screen.getByLabelText("メールアドレス"),
        "unknown@example.com"
        );

        await user.type(
        screen.getByLabelText("パスワード"),
        "wrong-password"
        );

        await user.click(
        screen.getByRole("button", {
            name: "ログイン",
        })
        );

        expect(
        await screen.findByText(
            "メールアドレスまたはパスワードが正しくありません"
        )
        ).toBeInTheDocument();

        expect(
        screen.queryByRole("heading", {
            name: "会議一覧",
        })
        ).not.toBeInTheDocument();
    });


    test("ログイン成功後に会議一覧へ移動する", async () => {
        const user = userEvent.setup();
        const loginMock = vi.fn().mockResolvedValue(undefined);

        mockUseAuth.mockReturnValue({
          user: null,
          loading: false,
          login: loginMock,
          logout: vi.fn(),
          fetchCurrentUser: vi.fn(),
          deactivateAccount: vi.fn(),
        });

        renderLoginPage();

        await user.type(
        screen.getByLabelText("メールアドレス"),
        "employee1@example.com"
        );

        await user.type(
        screen.getByLabelText("パスワード"),
        "password123"
        );

        await user.click(
        screen.getByRole("button", {
            name: "ログイン",
        })
        );

        expect(
        await screen.findByRole("heading", {
            name: "会議一覧",
        })
        ).toBeInTheDocument();
    });

    
});