const apiBaseUrl =
  import.meta.env.VITE_API_BASE_URL;

const wsBaseUrl =
  import.meta.env.VITE_WS_BASE_URL;

if (!apiBaseUrl) {
  throw new Error(
    "VITE_API_BASE_URLが設定されていません"
  );
}

if (!wsBaseUrl) {
  throw new Error(
    "VITE_WS_BASE_URLが設定されていません"
  );
}

export const API_BASE_URL = apiBaseUrl;
export const WS_BASE_URL = wsBaseUrl;