import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import type { AgendaForm } from "../types/meeting";
import { API_BASE_URL } from "../config/env";
import ErrorMessage from "../components/ErrorMessage";
import MeetingForm from "../components/MeetingForm";
import { ArrowLeft } from "lucide-react";

export default function MeetingCreatePage() {
  // エラーメッセージ
  const [errorMessage, setErrorMessage] = useState("");

  const [saving, setSaving] = useState(false);

  // Meetings
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [targetName, setTargetName] = useState("");
  const [scheduledStartAt, setScheduledStartAt] = useState("");
  const [plannedMinutes, setPlannedMinutes] = useState("60");

  // Agendas
  const [agendas, setAgendas] =
    useState<AgendaForm[]>([
      {
        title: "",
        purpose: "",
        discussion_points: "",
        questions: "",
        memo: "",
        planned_minutes: "10",
      },
    ]);

  const navigate = useNavigate();


  // 登録処理
  const handleSubmit = async () => {
    // console.log(
    //   JSON.stringify(
    //     {
    //       title,
    //       target_name: targetName,
    //       description,
    //       scheduled_start_at: scheduledStartAt,
    //       planned_minutes: plannedMinutes,
    //       agendas,
    //     },
    //     null,
    //     2
    //   )
    // );

    try {
      setErrorMessage("");
      setSaving(true);

      const response = await fetch(`${API_BASE_URL}/api/meetings`, {
        method: "POST",
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
        navigate("/");
        return
      }

      const error = await response.json();
      setErrorMessage(error.error);

    } catch (error) {
      console.error(error);
      setErrorMessage("会議情報の登録に失敗しました");

    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="mx-auto max-w-[1600px] p-6 lg:p-8">
      <Link
        to="/"
        className="inline-flex items-center text-slate-500 hover:text-slate-900"
      >
        <ArrowLeft size={18} /> 会議一覧へ戻る
      </Link>

      <div className="flex justify-between items-start my-6">
        <h1 className="text-3xl font-bold text-slate-900">
          会議新規作成
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
        agendas={agendas}
        setAgendas={setAgendas}
      />
    </div>
  );
}