import { useState } from "react";
import type { AgendaForm } from "../types/meeting";
import AgendaSidebar from "../components/AgendaSidebar";

type MeetingFormProps = {
	title: string;
	setTitle: React.Dispatch<React.SetStateAction<string>>;

	description: string;
	setDescription: React.Dispatch<React.SetStateAction<string>>;

	targetName: string;
	setTargetName: React.Dispatch<React.SetStateAction<string>>;

	scheduledStartAt: string;
	setScheduledStartAt: React.Dispatch<React.SetStateAction<string>>;

	plannedMinutes: string;
	setPlannedMinutes: React.Dispatch<React.SetStateAction<string>>;

	// Editのときだけ受け取る
	decisions?: string;
	setDecisions?: React.Dispatch<React.SetStateAction<string>>;

	// Editのときだけ受け取る
	todo?: string;
	setTodo?: React.Dispatch<React.SetStateAction<string>>;

	agendas: AgendaForm[];
	setAgendas: React.Dispatch<React.SetStateAction<AgendaForm[]>>;
};

export default function MeetingForm({
	title,
	setTitle,
	description,
	setDescription,
	targetName,
	setTargetName,
	scheduledStartAt,
	setScheduledStartAt,
	plannedMinutes,
	setPlannedMinutes,
	decisions,
	setDecisions,
	todo,
	setTodo,
	agendas,
	setAgendas,
}: MeetingFormProps) {

	const [selectedAgendaIndex, setSelectedAgendaIndex] = useState(0);
	const selectedAgenda = agendas[selectedAgendaIndex] ?? null;

	const updateAgenda = (
		field: keyof AgendaForm,
		value: string
	) => {
		const updated = [...agendas];

		updated[selectedAgendaIndex] = {
			...updated[selectedAgendaIndex],
			[field]: value,
		};

		setAgendas(updated);
	};

  // agendaの並び替え
  const moveAgenda = (
    index: number,
    direction: "up" | "down"
  ) => {
    const targetIndex =
      direction === "up" ? index - 1 : index + 1;

    if (targetIndex < 0 || targetIndex >= agendas.length) {
      return;
    }

    const updated = [...agendas];

    [updated[index], updated[targetIndex]] = [
      updated[targetIndex],
      updated[index],
    ];

    setAgendas(updated);
    setSelectedAgendaIndex(targetIndex);
  };

	return (
		<>
			{/* 会議作成フォーム */}
			<div className="mt-6 bg-white/70 border border-white/70 rounded-2xl p-4 shadow-sm backdrop-blur-xl">
				<h2 className="mb-4 text-lg font-semibold">会議情報</h2>

				<div className="space-y-4">
					<div>
						<label className="mb-1 block text-sm font-medium text-slate-700">会議名</label>
						<input
							type="text"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
							className="w-full border rounded-xl border-slate-300 px-3 py-2"
						/>
					</div>

					<div>
						<label className="mb-1 block text-sm font-medium text-slate-700">会議詳細</label>
						<textarea
							value={description}
							onChange={(e) => setDescription(e.target.value)}
							rows={4}
							className="w-full border rounded-xl border-slate-300 px-3 py-2"
						/>
					</div>

					<div>
						<label className="mb-1 block text-sm font-medium text-slate-700">会議相手</label>
						<input
							type="text"
							value={targetName}
							onChange={(e) => setTargetName(e.target.value)}
							className="w-full border rounded-xl border-slate-300 px-3 py-2"
						/>
					</div>

					<div className="flex gap-4">
						<div>
							<label className="mb-1 block text-sm font-medium text-slate-700">開始日時</label>
							<input
								type="datetime-local"
								value={scheduledStartAt}
								onChange={(e) => setScheduledStartAt(e.target.value)}
								className="w-full border rounded-xl border-slate-300 px-3 py-2"
							/>
						</div>

						<div>
							<label className="mb-1 block text-sm font-medium text-slate-700">予定時間(分)</label>
							<input
								type="number"
								value={plannedMinutes}
								onChange={(e) =>
									setPlannedMinutes(e.target.value)
								}
								className="w-full border border-slate-300 rounded-xl px-3 py-2"
							/>
						</div>
					</div>
				</div>
			</div>

			{/* 会議後入力項目 */}
			{(setDecisions && setTodo) && (
				<div className="grid gap-6 md:grid-cols-2">
					<div className="mt-6 bg-white/70 border border-white/70 rounded-2xl p-4 shadow-sm backdrop-blur-xl">
							<label className="mb-1 block text-sm font-medium text-slate-700">決定事項</label>
							<textarea
								value={decisions}
								onChange={(e) => setDecisions(e.target.value)}
								rows={4}
								className="w-full border rounded-xl border-slate-300 px-3 py-2"
							/>
					</div>
					<div className="mt-6 bg-white/70 border border-white/70 rounded-2xl p-4 shadow-sm backdrop-blur-xl">
							<label className="mb-1 block text-sm font-medium text-slate-700">TODO</label>
							<textarea
								value={todo}
								onChange={(e) => setTodo(e.target.value)}
								rows={4}
								className="w-full border rounded-xl border-slate-300 px-3 py-2"
							/>
					</div>
				</div>
			)}

			{/* アジェンダ作成フォーム */}
			<div className="my-6">

				<div className="mt-8 grid gap-6 lg:grid-cols-[320px_minmax(0,1fr)]">

					{/* 左ペイン */}
					<div>
						<AgendaSidebar
							agendas={agendas}
							selectedIndex={selectedAgendaIndex}
							meetingPlannedMinutes={plannedMinutes}
							onSelect={setSelectedAgendaIndex}
              onMove={moveAgenda}
						/>

						<button
							type="button"
							onClick={() => {
								setAgendas([
									...agendas,
									{
										title: "",
										purpose: "",
										discussion_points: "",
										questions: "",
										memo: "",
										planned_minutes: "10",
									},
								]);

								setSelectedAgendaIndex(agendas.length);
							}}
							className="mt-4 w-full rounded-xl bg-slate-900 px-4 py-2 text-white"
						>
							+ 議題追加
						</button>
					</div>

					{/* 右ペイン */}
					<main className="self-start rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
						{/* ヘッダー */}
						<div className="border-b border-slate-200 pb-4">
							<p className="text-sm font-medium text-slate-500">
								議題 {selectedAgendaIndex + 1} / {agendas.length}
							</p>

							<div className="mt-2 flex items-center justify-between">
								<h2 className="text-2xl font-bold text-slate-900">
									{selectedAgenda?.title || "新しい議題"}
								</h2>

								{agendas.length > 1 && (
									<button
										type="button"
										onClick={() => {
											const next = agendas.filter(
												(_, i) => i !== selectedAgendaIndex
											);

											setAgendas(next);
											setSelectedAgendaIndex(
												Math.max(0, selectedAgendaIndex - 1)
											);
										}}
										className="text-red-600 hover:text-red-800 cursor-pointer"
									>
										削除
									</button>
								)}
							</div>
						</div>

						{/* フォーム部分 */}
						<div className="mt-5">
							<label className="mb-1 block text-sm font-medium">
								議題名
							</label>

							<input
								value={selectedAgenda?.title ?? ""}
								onChange={(e) => {
									updateAgenda("title", e.target.value)
								}}
								className="w-full rounded-xl border border-slate-300 px-3 py-2 mb-2"
							/>

							<label className="mb-1 block text-sm font-medium text-slate-700">
								予定時間（分）
							</label>

							<input
								type="number"
								value={selectedAgenda?.planned_minutes ?? ""}
								onChange={(e) => {
									updateAgenda("planned_minutes", e.target.value)
								}}
								className="w-32 border border-slate-300 rounded-xl px-3 py-2 mb-2"
							/>

							<label className="mb-1 block text-sm font-medium text-slate-700">
								目的
							</label>

							<textarea
								value={selectedAgenda?.purpose ?? ""}
								onChange={(e) => {
									updateAgenda("purpose", e.target.value)
								}}
								rows={2}
								className="w-full border border-slate-300 rounded-xl px-3 py-2 mb-2"
							/>

							<div className="grid gap-5 lg:grid-cols-2 items-stretch">
								<section className="h-full rounded-2xl">
									<label className="mb-1 block text-sm font-medium text-slate-700">
										議論ポイント
									</label>

									<textarea
										value={selectedAgenda?.discussion_points ?? ""}
										onChange={(e) => {
											updateAgenda("discussion_points", e.target.value)
										}}
										rows={3}
										className="w-full border border-slate-300 rounded-xl px-3 py-2"
									/>
								</section>

								<section className="h-full rounded-2xl">
									<label className="block mb-1 text-sm font-medium text-slate-700">
										質問事項
									</label>

									<textarea
										value={selectedAgenda?.questions ?? ""}
										onChange={(e) => {
											updateAgenda("questions", e.target.value)
										}}
										rows={3}
										className="w-full rounded-xl border border-slate-300 px-3 py-2"
									/>
								</section>
							</div>

							<label className="block mb-1 text-sm font-medium text-slate-700">
								メモ
							</label>

							<textarea
								value={selectedAgenda?.memo ?? ""}
								onChange={(e) => {
									updateAgenda("memo", e.target.value)
								}}
								rows={3}
								className="w-full rounded-xl border border-slate-300 px-3 py-2"
							/>
						</div>
					</main>

				</div>

			</div>
		</>
	);
}