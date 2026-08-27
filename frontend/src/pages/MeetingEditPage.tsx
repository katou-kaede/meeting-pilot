import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import type { AgendaForm, MeetingDetail } from "../types/meeting";
import ErrorMessage from "../components/ErrorMessage";
import MeetingForm from "../components/MeetingForm";
import Loading from "../components/Loading";
import MeetingMemberManager from "../components/MeetingMemberManager";
import { ArrowLeft } from "lucide-react";

export default function MeetingEditPage() {
  const { id } = useParams();
  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState("");
  const [loading, setLoading] = useState(true);

  const [saving, setSaving] = useState(false);

  // Meetings
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [targetName, setTargetName] = useState("");
  const [scheduledStartAt, setScheduledStartAt] = useState("");
  const [plannedMinutes, setPlannedMinutes] = useState("60");
  const [decisions, setDecisions] = useState("");
  const [todo, setTodo] = useState("");

  // Agendas
  const [agendas, setAgendas] = useState<AgendaForm[]>([]);

  const navigate = useNavigate();

  // データ取得
  useEffect(() => {
    const fetchMeeting = async () => {
      try {
        setErrorMessage("");

        const response = await fetch(
          `http://localhost:8080/api/meetings/${id}`,
          {
            credentials: "include",
          }
        );

        const data = await response.json();

        if (!response.ok) {
          setErrorMessage(data.error);
          return;
        }

        const meetingData = data as MeetingDetail;

        setTitle(meetingData.title);
        setDescription(meetingData.description);
        setTargetName(meetingData.target_name);
        // datetime-localに合う形式に変換
        setScheduledStartAt(
          meetingData.scheduled_start_at
            ? meetingData.scheduled_start_at.slice(0, 16)
            : ""
        );
        setPlannedMinutes(String(meetingData.planned_minutes));
        setDecisions(meetingData.decisions);
        setTodo(meetingData.todo);

        setAgendas(
          (meetingData.agendas ?? []).map((agenda) => ({
            id: agenda.id,
            title: agenda.title,
            purpose: agenda.purpose,
            discussion_points: agenda.discussion_points,
            questions: agenda.questions,
            memo: agenda.memo,
            planned_minutes: String(
              agenda.planned_minutes
            ),
          }))
        );

      } catch (error) {
        console.error(error);
        setErrorMessage("会議情報の取得に失敗しました");

      } finally {
        setLoading(false);
      }
    };

    fetchMeeting();
  }, [id]);

  // 更新処理
  const handleSubmit = async () => {
    // console.log(
    //   JSON.stringify(
    //     {
    //       title,
    //       target_name: targetName,
    //       description,
    //       scheduled_start_at: scheduledStartAt,
    //       planned_minutes: plannedMinutes,
    //       decisions,
    //       todo,
    //       agendas,
    //     },
    //     null,
    //     2
    //   )
    // );

    try {
      setErrorMessage("");
      setSaving(true);

      const response = await fetch(`http://localhost:8080/api/meetings/${id}`, {
        method: "PUT",
        credentials: "include",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          title,
          description,
          target_name: targetName,
          scheduled_start_at: scheduledStartAt,
          planned_minutes: Number(plannedMinutes),
          decisions,
          todo,

          agendas: agendas.map((agenda) => ({
            id: agenda.id,
            title: agenda.title,
            purpose: agenda.purpose,
            discussion_points: agenda.discussion_points,
            questions: agenda.questions,
            memo: agenda.memo,
            planned_minutes: Number(agenda.planned_minutes),
          }))
        })
      })

      if (response.ok) {
        navigate(`/meetings/${id}`);
        return
      }

      const error = await response.json();
      setErrorMessage(error.error);

    } catch (error) {
      setErrorMessage("会議の更新に失敗しました");
      console.error(error);

    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="mx-auto max-w-[1600px] p-6 lg:p-8">
      <Link
        to={`/meetings/${id}`}
        className="inline-flex items-center text-slate-500 hover:text-slate-900"
      >
        <ArrowLeft size={18} /> 会議詳細へ戻る
      </Link>

      <div className="flex justify-between items-start my-6">
        <h1 className="text-3xl font-bold text-slate-900">
          会議情報編集
        </h1>

        <button
          onClick={handleSubmit}
          disabled={saving}
          className="rounded-xl bg-blue-600 px-8 py-3 font-medium text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        >
          {saving ? "保存中..." : "保存"}
        </button>
      </div>

      {/* エラーメッセージ */}
      {errorMessage && (
        <ErrorMessage message={errorMessage} />
      )}

      <MeetingForm
        title={title}
        setTitle={setTitle}
        description={description}
        setDescription={setDescription}
        targetName={targetName}
        setTargetName={setTargetName}
        scheduledStartAt={scheduledStartAt}
        setScheduledStartAt={setScheduledStartAt}
        plannedMinutes={plannedMinutes}
        setPlannedMinutes={setPlannedMinutes}
        decisions={decisions}
        setDecisions={setDecisions}
        todo={todo}
        setTodo={setTodo}
        agendas={agendas}
        setAgendas={setAgendas}
      />

      <MeetingMemberManager meetingId={Number(id)} />
    </div>
  );
}