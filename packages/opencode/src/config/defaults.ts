export const defaultProviders = {
  zerosync: {
    name: "ZeroSync",
    npm: "@ai-sdk/anthropic",
    options: {
      baseURL: "https://llm-provider-proxy.alex-dunne.workers.dev/v1",
    },
    models: {
      "claude-sonnet-4-20250514": {
        name: "Claude Sonnet 4",
        attachment: true,
        reasoning: false,
        temperature: true,
        tool_call: true,
        cost: {
          input: 0,
          output: 0,
          cache_read: 0,
          cache_write: 0,
        },
        limit: {
          context: 200000,
          output: 64000,
        },
      },
      "claude-3-5-haiku-20241022": {
        name: "Claude Haiku 3.5",
        attachment: true,
        reasoning: false,
        temperature: true,
        tool_call: true,
        knowledge: "2024-07-31",
        release_date: "2024-10-22",
        last_updated: "2024-10-22",
        modalities: { input: ["text", "image"], output: ["text"] },
        open_weights: false,
        cost: { input: 0.8, output: 4, cache_read: 0.08, cache_write: 1 },
        limit: { context: 200000, output: 8192 },
      },
      language: {
        specificationVersion: "v2",
        modelId: "claude-3-5-haiku-20241022",
        config: { provider: "anthropic.messages", baseURL: "https://llm-provider-proxy.alex-dunne.workers.dev/v1" },
      },
    },
  },
}

export const defaultModel = "zerosync/claude-sonnet-4-20250514"
