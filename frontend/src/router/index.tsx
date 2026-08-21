import { createBrowserRouter } from "react-router-dom";
import MeetingListPage from "../pages/MeetingListPage";
import MeetingCreatePage from "../pages/MeetingCreatePage";
import MeetingDetailPage from "../pages/MeetingDetailPage";
import MeetingEditPage from "../pages/MeetingEditPage";
import MeetingSessionPage from "../pages/MeetingSessionPage";

export const router = createBrowserRouter([
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
  }
]);
