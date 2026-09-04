import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";

import type { User } from "../types/user";
import { API_BASE_URL } from "../config/env";

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (
    email: string,
    password: string
  ) => Promise<void>;
  logout: () => Promise<void>;
  deactivateAccount: () => Promise<void>;
  fetchCurrentUser: () => Promise<void>;
};

const AuthContext = createContext<
  AuthContextType | undefined
>(undefined);

type Props = {
  children: ReactNode;
};

export function AuthProvider({ children }: Props) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchCurrentUser = async () => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/me`,
        {
          credentials: "include",
        }
      );

      if (!response.ok) {
        setUser(null);
        return;
      }

      const data: User = await response.json();

      setUser(data);
    } catch (error) {
      console.error(
        "認証ユーザーの取得に失敗しました",
        error
      );

      setUser(null);
    }
  };

  const login = async (
    email: string,
    password: string
  ) => {
    const response = await fetch(
      `${API_BASE_URL}/api/login`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          email,
          password,
        }),
      }
    );

    const data = await response.json();

    if (!response.ok) {
      throw new Error(
        data.error || "ログインに失敗しました"
      );
    }

    setUser(data as User);
  };

  const logout = async () => {
    const response = await fetch(
        `${API_BASE_URL}/api/logout`,
        {
        method: "POST",
        credentials: "include",
        }
    );

    if (!response.ok) {
        throw new Error("ログアウトに失敗しました");
    }

    setUser(null);
  };

  const deactivateAccount = async () => {
    const response = await fetch(
      `${API_BASE_URL}/api/me`,
      {
        method: "DELETE",
        credentials: "include",
      }
    );

    if (!response.ok) {
      const data = await response.json();

      throw new Error(
        data.error || "退会処理に失敗しました"
      );
    }

    setUser(null);
  };

  useEffect(() => {
    const initializeAuth = async () => {
      setLoading(true);

      await fetchCurrentUser();

      setLoading(false);
    };

    initializeAuth();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        deactivateAccount,
        fetchCurrentUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error(
      "useAuthはAuthProvider内で使用してください"
    );
  }

  return context;
}