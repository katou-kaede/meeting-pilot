import { render, screen } from "@testing-library/react";
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

import ProtectedRoute from "./ProtectedRoute";
import { useAuth } from "../contexts/AuthContext";

vi.mock("../contexts/AuthContext", () => ({
  useAuth: vi.fn(),
}));

vi.mock("./Loading", () => ({
  default: () => <div>読み込み中</div>,
}));

const mockUseAuth = vi.mocked(useAuth);

function renderProtectedRoute() {
  render(
    <MemoryRouter initialEntries={["/private"]}>
      <Routes>
        <Route
          path="/private"
          element={
            <ProtectedRoute>
              <h1>会議一覧</h1>
            </ProtectedRoute>
          }
        />

        <Route
          path="/login"
          element={<h1>ログイン画面</h1>}
        />
      </Routes>
    </MemoryRouter>
  );
}

describe("ProtectedRoute", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("認証状態の確認中はLoadingを表示する", () => {
    mockUseAuth.mockReturnValue({
      user: null,
      loading: true,
      login: vi.fn(),
      logout: vi.fn(),
      fetchCurrentUser: vi.fn(),
      deactivateAccount: vi.fn(),
    });

    renderProtectedRoute();

    expect(
      screen.getByText("読み込み中")
    ).toBeInTheDocument();

    expect(
      screen.queryByText("会議一覧")
    ).not.toBeInTheDocument();
  });

  test("未ログインならログイン画面へ移動する", () => {
    mockUseAuth.mockReturnValue({
      user: null,
      loading: false,
      login: vi.fn(),
      logout: vi.fn(),
      fetchCurrentUser: vi.fn(),
      deactivateAccount: vi.fn(),
    });

    renderProtectedRoute();

    expect(
      screen.getByRole("heading", {
        name: "ログイン画面",
      })
    ).toBeInTheDocument();

    expect(
      screen.queryByRole("heading", {
        name: "会議一覧",
      })
    ).not.toBeInTheDocument();
  });

  test("ログイン済みなら子画面を表示する", () => {
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
      deactivateAccount: vi.fn(),
    });

    renderProtectedRoute();

    expect(
      screen.getByRole("heading", {
        name: "会議一覧",
      })
    ).toBeInTheDocument();

    expect(
      screen.queryByRole("heading", {
        name: "ログイン画面",
      })
    ).not.toBeInTheDocument();
  });
});