import { Link } from "react-router-dom";
import { useState, useEffect } from "react";
import type { Meeting } from "../types/meeting";
import { formatDateTime } from "../utils/common";
import { getStatusLabel, getStatusStyle } from "../utils/meetingStatus";
import ErrorMessage from "../components/ErrorMessage";


export default function MeetingListPage() {
  const [meetings, setMeetings] = useState<Meeting[]>([]);

  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMeetings = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/meetings")

        const data = await response.json();
        
        if (!response.ok) {
          setErrorMessage(data.error);
          return;
        }
        
        setMeetings(data);

      } catch (error) {
          console.error(error);
          setErrorMessage("会議一覧の取得に失敗しました")

      } finally {
        setLoading(false);
      }
    };

    fetchMeetings();
  }, [])

  if (loading) {
    return <div>Loading...</div>;
  }


  return (
    <div className="bg-slate-50 min-h-screen">
      <div className="max-w-7xl mx-auto p-8">
        <div className="flex justify-between items-center mb-8">

          {/* ヘッダー */}
          <div>
            <h1 className="text-3xl font-bold text-slate-900">
              MeetingPilot
            </h1>

            <p className="text-slate-500 mt-1">
              会議一覧
            </p>
          </div>

          <Link
            to="/meetings/new"
            className="bg-slate-900 text-white px-5 py-3 rounded-xl hover:bg-slate-800"
          >
            ＋ 新規会議
          </Link>

        </div>

        {/* エラーメッセージ */}
        {errorMessage && (
          <ErrorMessage message={errorMessage} />
        )}

        <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden">
          <table className="w-full">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr className="text-left">
                <th className="px-6 py-4 text-sm font-semibold text-slate-600">会議名</th>
                <th className="px-6 py-4 text-sm font-semibold text-slate-600">会議相手</th>
                <th className="px-6 py-4 text-sm font-semibold text-slate-600">開始日時</th>
                <th className="px-6 py-4 text-sm font-semibold text-slate-600">時間</th>
                <th className="px-6 py-4 text-sm font-semibold text-slate-600">状態</th>
              </tr>
            </thead>

            <tbody>
              {(meetings ?? []).map((meeting) => (
                <tr key={meeting.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="px-6 py-4 font-medium">
                    <Link
                      to={`/meetings/${meeting.id}`}
                      className="hover:text-blue-600"
                    >
                      {meeting.title}
                    </Link>
                  </td>

                  <td className="px-6 py-4 text-slate-600">{meeting.target_name}</td>

                  <td className="px-6 py-4 text-slate-600">
                    {formatDateTime(meeting.scheduled_start_at)}
                  </td>

                  <td className="px-6 py-4 text-slate-600">{meeting.planned_minutes}分</td>

                  <td className="px-6 py-4 text-slate-600">
                    <span className={`inline-block px-3 py-1 rounded-full text-sm ${getStatusStyle(meeting.status)}`}>
                      {getStatusLabel(meeting.status)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}