import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { Type } from "typebox";

function valueOrPlaceholder(value: string | undefined, placeholder: string): string {
	const trimmed = value?.trim();
	return trimmed && trimmed.length > 0 ? trimmed : placeholder;
}

function buildInstructions(params: {
	task_path?: string;
	commit_hash?: string;
	validation_summary?: string;
	next_work?: string;
	config_summary?: string;
}): string {
	const lines = [
		"# Roadmap context continuation",
		"",
		"Preserve the current /roadmap loop state so another turn can continue safely.",
		"",
		"## Current roadmap goal",
		"Continue the active roadmap outcome/task execution in this repository using roadmapctl as the source of truth.",
		"",
		"## Completed task",
		`- Task path: ${valueOrPlaceholder(params.task_path, "not provided")}`,
		`- Commit hash: ${valueOrPlaceholder(params.commit_hash, "not provided")}`,
		"",
		"## Validation results",
		valueOrPlaceholder(params.validation_summary, "No validation summary was provided."),
		"",
		"## Next task or wave state",
		valueOrPlaceholder(params.next_work, "No next work summary was provided. Re-run roadmapctl next before continuing."),
		"",
		"## Unresolved blockers or conflicts",
		"Preserve any blockers, failed pushes, dirty working tree conflicts, or unsafe parallelization constraints mentioned in validation results or next work.",
		"",
		"## Relevant config values",
		valueOrPlaceholder(params.config_summary, "No config summary was provided. Re-run roadmapctl context before continuing."),
		"",
		"## Continuation checklist",
		"1. Re-run roadmapctl context, doctor, check, and next before mutating roadmap state.",
		"2. Do not rely on legacy .claude/roadmap.local.md as a lasting config source.",
		"3. Resume only tasks reported ready by roadmapctl next.",
		"4. Preserve exact user/worktree changes and avoid staging unrelated files.",
		"5. Verify acceptance criteria before transition complete, commit, push, or PR bookkeeping.",
	];
	return lines.join("\n");
}

export default function (pi: ExtensionAPI) {
	pi.registerTool({
		name: "compact_roadmap_context",
		label: "Compact Roadmap Context",
		description: "Queue Pi compaction with roadmap-specific continuation instructions.",
		promptSnippet: "Queue roadmap-specific context compaction after a completed roadmap task is durable.",
		promptGuidelines: [
			"Use compact_roadmap_context after a roadmap task commit/push or PR bookkeeping completes when repo config requests roadmap context compaction.",
		],
		parameters: Type.Object({
			task_path: Type.Optional(Type.String({ description: "Completed roadmap task path." })),
			commit_hash: Type.Optional(Type.String({ description: "Commit hash for the completed task." })),
			validation_summary: Type.Optional(Type.String({ description: "Acceptance checks and roadmapctl validation summary." })),
			next_work: Type.Optional(Type.String({ description: "Next task/wave state, blockers, or conflicts to preserve." })),
			config_summary: Type.Optional(Type.String({ description: "Relevant roadmap execution config values." })),
		}),
		execute(_toolCallId, params, _signal, _onUpdate, ctx) {
			const customInstructions = buildInstructions(params);
			try {
				ctx.compact({
					customInstructions,
					onComplete: () => {
						if (ctx.hasUI) {
							ctx.ui.notify("Roadmap context compaction completed", "info");
						}
					},
					onError: (error) => {
						if (ctx.hasUI) {
							ctx.ui.notify(`Roadmap context compaction failed: ${error.message}`, "error");
						}
					},
				});
			} catch (error) {
				const message = error instanceof Error ? error.message : String(error);
				return {
					content: [{ type: "text" as const, text: `Roadmap context compaction failed to queue: ${message}` }],
					details: { queued: false, error: message, customInstructions },
				};
			}

			return {
				content: [{ type: "text" as const, text: "Roadmap context compaction queued." }],
				details: { queued: true, customInstructions },
			};
		},
	});
}
