import { createBrowserRouter } from "react-router-dom";
import MeetingListPage from "../pages/MeetingListPage";
import MeetingCreatePage from "../pages/MeetingCreatePage";
import MeetingDetailPage from "../pages/MeetingDetailPage";
import MeetingEditPage from "../pages/MeetingEditPage";
import MeetingSessionPage from "../pages/MeetingSessionPage";
import LoginPage from "../pages/LoginPage";
import UserCreatePage from "../pages/UserCreatePage";
import ProtectedRoute from "../components/ProtectedRoute";
import AppLayout from "../layouts/AppLayout";

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/users/create",
    element: <UserCreatePage />,
  },
  {
    element: (
      <ProtectedRoute>
        <AppLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        path: "/",
        element: <MeetingListPage />,
      },
      {
        path: "/meetings/new",
        element: <MeetingCreatePage />,
      },
      {
        path: "/meetings/:id",
        element: <MeetingDetailPage />,
      },
      {
        path: "/meetings/:id/edit",
        element: <MeetingEditPage />
      },
      {
        path: "/meetings/:id/session",
        element: <MeetingSessionPage />
      },
    ]
  }
]);
