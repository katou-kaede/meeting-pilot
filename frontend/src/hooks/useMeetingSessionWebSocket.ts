import { useEffect } from "react";

type Props = {
  meetingId: string | undefined;
  onSessionUpdated: () => void;
  onMeetingCompleted: () => void;
  onConnectionError: (message: string) => void;
};

export function useMeetingSessionWebSocket({
  meetingId,
  onSessionUpdated,
  onMeetingCompleted,
  onConnectionError,
}: Props) {
  useEffect(() => {
    if (!meetingId) return;

    const socket = new WebSocket(
      `ws://localhost:8080/ws/meetings/${meetingId}`
    );

    socket.onopen = () => {
      onConnectionError("");
      console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);

        switch (message.type) {
          case "current_agenda_changed":
          case "meeting_paused":
          case "meeting_resumed":
          case "meeting_session_saved":
          case "meeting_editor_changed":
            onSessionUpdated();
            break;

          case "meeting_completed":
            onMeetingCompleted();
            break;
        }
      } catch (error) {
        console.error(
          "WebSocket message parse failed:",
          error
        );
      }
    };

    socket.onerror = (error) => {
      onConnectionError(
        "リアルタイム接続に問題が発生しました"
      );

      console.error("WebSocket error:", error);
    };

    socket.onclose = (event) => {
      console.log(
        "WebSocket disconnected",
        event.code
      );

      if (event.code !== 1000) {
        onConnectionError(
          "リアルタイム接続が切断されました。画面を再読み込みしてください"
        );
      }
    };

    return () => {
      socket.close();
    };
  }, [
    meetingId,
    onSessionUpdated,
    onMeetingCompleted,
    onConnectionError,
  ]);
}