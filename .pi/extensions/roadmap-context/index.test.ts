import { beforeAll, describe, expect, mock, test } from "bun:test";

mock.module("typebox", () => ({
	Type: {
		Object: (schema: unknown) => ({ type: "object", properties: schema }),
		Optional: (schema: unknown) => schema,
		String: (options?: unknown) => ({ type: "string", ...((options as object | undefined) ?? {}) }),
	},
}));

let registerRoadmapContextExtension: (pi: any) => void;

beforeAll(async () => {
	registerRoadmapContextExtension = (await import("./index")).default;
});

function registerToolForTest() {
	let registeredTool: any;
	const sentMessages: Array<{ content: string; options: unknown }> = [];
	const pi = {
		registerTool(definition: any) {
			registeredTool = definition;
		},
		sendUserMessage(content: string, options?: unknown) {
			sentMessages.push({ content, options });
		},
	};

	registerRoadmapContextExtension(pi as any);
	if (!registeredTool) {
		throw new Error("extension did not register compact_roadmap_context");
	}
	return { registeredTool, sentMessages };
}

function executeToolAndCaptureCompactOptions() {
	const { registeredTool, sentMessages } = registerToolForTest();
	let compactOptions: any;
	const ctx = {
		hasUI: false,
		compact(options: any) {
			compactOptions = options;
		},
	};

	const result = registeredTool.execute(
		"tool-call-1",
		{ task_path: "docs/roadmap/T001-example.md", next_work: "Run roadmapctl next." },
		undefined,
		undefined,
		ctx,
	);

	expect(result.content[0].text).toBe("Roadmap context compaction queued.");
	expect(compactOptions).toBeDefined();
	return { compactOptions, sentMessages };
}

describe("compact_roadmap_context", () => {
	test("queues a follow-up roadmap continuation after compaction completes", () => {
		const { compactOptions, sentMessages } = executeToolAndCaptureCompactOptions();

		compactOptions.onComplete({ summary: "compacted" });

		expect(sentMessages).toHaveLength(1);
		expect(sentMessages[0].content).toContain("Continue the /roadmap loop");
		expect(sentMessages[0].content).toContain("roadmapctl context");
		expect(sentMessages[0].options).toEqual({ deliverAs: "followUp" });
	});

	test("queues the same continuation when compaction fails after aborting the active turn", () => {
		const { compactOptions, sentMessages } = executeToolAndCaptureCompactOptions();

		compactOptions.onError(new Error("Nothing to compact (session too small)"));

		expect(sentMessages).toHaveLength(1);
		expect(sentMessages[0].content).toContain("Compaction failed after the durable task");
		expect(sentMessages[0].content).toContain("Continue the /roadmap loop");
	});

	test("does not auto-continue when compaction is explicitly cancelled", () => {
		const { compactOptions, sentMessages } = executeToolAndCaptureCompactOptions();

		compactOptions.onError(new Error("Compaction cancelled"));

		expect(sentMessages).toHaveLength(0);
	});
});
