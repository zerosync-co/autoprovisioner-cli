import z from "zod"
import { App } from "../app/app"
import { Config } from "../config/config"
import { mergeDeep, sortBy } from "remeda"
import { NoSuchModelError, type LanguageModel, type Provider as SDK } from "ai"
import { Log } from "../util/log"
import { BunProc } from "../bun"
import { BashTool } from "../tool/bash"
import { CredentialBashTool } from "../tool/credential-bash"
import { EditTool } from "../tool/edit"
import { WebFetchTool } from "../tool/webfetch"
import { GlobTool } from "../tool/glob"
import { GrepTool } from "../tool/grep"
import { ListTool } from "../tool/ls"
import { PatchTool } from "../tool/patch"
import { ReadTool } from "../tool/read"
import type { Tool } from "../tool/tool"
import { WriteTool } from "../tool/write"
import { TodoReadTool, TodoWriteTool } from "../tool/todo"
import { AuthAnthropic } from "../auth/anthropic"
import { AuthCopilot } from "../auth/copilot"
import { AuthZerosync } from "../auth/zerosync"
import { ModelsDev } from "./models"
import { NamedError } from "../util/error"
import { Auth } from "../auth"
import { TaskTool } from "../tool/task"

export namespace Provider {
  const log = Log.create({ service: "provider" })

  type CustomLoader = (
    provider: ModelsDev.Provider,
    api?: string,
  ) => Promise<{
    autoload: boolean
    getModel?: (sdk: any, modelID: string) => Promise<any>
    options?: Record<string, any>
  }>

  type Source = "env" | "config" | "custom" | "api"

  const CUSTOM_LOADERS: Record<string, CustomLoader> = {
    async zerosync(provider) {
      const access = await AuthZerosync.access()
      if (!access) return { autoload: false }
      for (const model of Object.values(provider.models)) {
        model.cost = {
          input: 0,
          output: 0,
        }
      }
      return {
        autoload: true,
        options: {
          apiKey: "",
          async fetch(input: any, init: any) {
            const access = await AuthZerosync.access()

            delete init.headers["authorization"]
            delete init.headers["Authorization"]

            const headers = {
              ...init.headers,
              authorization: `Bearer ${access}`,
            }
            delete headers["x-api-key"]
            return fetch(input, {
              ...init,
              headers,
            })
          },
        },
      }
    },
    async anthropic(provider) {
      const access = await AuthAnthropic.access()
      if (!access) return { autoload: false }
      for (const model of Object.values(provider.models)) {
        model.cost = {
          input: 0,
          output: 0,
        }
      }
      return {
        autoload: true,
        options: {
          apiKey: "",
          async fetch(input: any, init: any) {
            const access = await AuthAnthropic.access()
            const headers = {
              ...init.headers,
              authorization: `Bearer ${access}`,
              "anthropic-beta": "oauth-2025-04-20",
            }
            delete headers["x-api-key"]
            return fetch(input, {
              ...init,
              headers,
            })
          },
        },
      }
    },
    "github-copilot": async (provider) => {
      const copilot = await AuthCopilot()
      if (!copilot) return { autoload: false }
      let info = await Auth.get("github-copilot")
      if (!info || info.type !== "oauth") return { autoload: false }

      if (provider && provider.models) {
        for (const model of Object.values(provider.models)) {
          model.cost = {
            input: 0,
            output: 0,
          }
        }
      }

      return {
        autoload: true,
        options: {
          apiKey: "",
          async fetch(input: any, init: any) {
            const info = await Auth.get("github-copilot")
            if (!info || info.type !== "oauth") return
            if (!info.access || info.expires < Date.now()) {
              const tokens = await copilot.access(info.refresh)
              if (!tokens) throw new Error("GitHub Copilot authentication expired")
              await Auth.set("github-copilot", {
                type: "oauth",
                ...tokens,
              })
              info.access = tokens.access
            }
            let isAgentCall = false
            let isVisionRequest = false
            try {
              const body = typeof init.body === "string" ? JSON.parse(init.body) : init.body
              if (body?.messages) {
                isAgentCall = body.messages.some((msg: any) => msg.role && ["tool", "assistant"].includes(msg.role))
                isVisionRequest = body.messages.some(
                  (msg: any) =>
                    Array.isArray(msg.content) && msg.content.some((part: any) => part.type === "image_url"),
                )
              }
            } catch {}
            const headers: Record<string, string> = {
              ...init.headers,
              ...copilot.HEADERS,
              Authorization: `Bearer ${info.access}`,
              "Openai-Intent": "conversation-edits",
              "X-Initiator": isAgentCall ? "agent" : "user",
            }
            if (isVisionRequest) {
              headers["Copilot-Vision-Request"] = "true"
            }
            delete headers["x-api-key"]
            return fetch(input, {
              ...init,
              headers,
            })
          },
        },
      }
    },
    openai: async () => {
      return {
        autoload: false,
        async getModel(sdk: any, modelID: string) {
          return sdk.responses(modelID)
        },
        options: {},
      }
    },
    "amazon-bedrock": async () => {
      if (!process.env["AWS_PROFILE"] && !process.env["AWS_ACCESS_KEY_ID"] && !process.env["AWS_BEARER_TOKEN_BEDROCK"])
        return { autoload: false }

      const region = process.env["AWS_REGION"] ?? "us-east-1"

      const { fromNodeProviderChain } = await import(await BunProc.install("@aws-sdk/credential-providers"))
      return {
        autoload: true,
        options: {
          region,
          credentialProvider: fromNodeProviderChain(),
        },
        async getModel(sdk: any, modelID: string) {
          let regionPrefix = region.split("-")[0]

          switch (regionPrefix) {
            case "us": {
              const modelRequiresPrefix = ["claude", "deepseek"].some((m) => modelID.includes(m))
              if (modelRequiresPrefix) {
                modelID = `${regionPrefix}.${modelID}`
              }
              break
            }
            case "eu": {
              const regionRequiresPrefix = [
                "eu-west-1",
                "eu-west-3",
                "eu-north-1",
                "eu-central-1",
                "eu-south-1",
                "eu-south-2",
              ].some((r) => region.includes(r))
              const modelRequiresPrefix = ["claude", "nova-lite", "nova-micro", "llama3", "pixtral"].some((m) =>
                modelID.includes(m),
              )
              if (regionRequiresPrefix && modelRequiresPrefix) {
                modelID = `${regionPrefix}.${modelID}`
              }
              break
            }
            case "ap": {
              const modelRequiresPrefix = ["claude", "nova-lite", "nova-micro", "nova-pro"].some((m) =>
                modelID.includes(m),
              )
              if (modelRequiresPrefix) {
                regionPrefix = "apac"
                modelID = `${regionPrefix}.${modelID}`
              }
              break
            }
          }

          return sdk.languageModel(modelID)
        },
      }
    },
    openrouter: async () => {
      return {
        autoload: false,
        options: {
          headers: {
            "HTTP-Referer": "https://opencode.ai/",
            "X-Title": "opencode",
          },
        },
      }
    },
  }

  const state = App.state("provider", async () => {
    const config = await Config.get()
    const database = await ModelsDev.get()

    const providers: {
      [providerID: string]: {
        source: Source
        info: ModelsDev.Provider
        getModel?: (sdk: any, modelID: string) => Promise<any>
        options: Record<string, any>
      }
    } = {}
    const models = new Map<string, { info: ModelsDev.Model; language: LanguageModel }>()
    const sdk = new Map<string, SDK>()

    log.info("init")

    function mergeProvider(
      id: string,
      options: Record<string, any>,
      source: Source,
      getModel?: (sdk: any, modelID: string) => Promise<any>,
    ) {
      const provider = providers[id]
      if (!provider) {
        const info = database[id]
        if (!info) return
        if (info.api && !options["baseURL"]) options["baseURL"] = info.api
        providers[id] = {
          source,
          info,
          options,
          getModel,
        }
        return
      }
      provider.options = mergeDeep(provider.options, options)
      provider.source = source
      provider.getModel = getModel ?? provider.getModel
    }

    const configProviders = Object.entries(config.provider ?? {})

    for (const [providerID, provider] of configProviders) {
      const existing = database[providerID]
      const parsed: ModelsDev.Provider = {
        id: providerID,
        npm: provider.npm ?? existing?.npm,
        name: provider.name ?? existing?.name ?? providerID,
        env: provider.env ?? existing?.env ?? [],
        api: provider.api ?? existing?.api,
        models: existing?.models ?? {},
      }

      for (const [modelID, model] of Object.entries(provider.models ?? {})) {
        const existing = parsed.models[modelID]
        const parsedModel: ModelsDev.Model = {
          id: modelID,
          name: model.name ?? existing?.name ?? modelID,
          release_date: model.release_date ?? existing?.release_date,
          attachment: model.attachment ?? existing?.attachment ?? false,
          reasoning: model.reasoning ?? existing?.reasoning ?? false,
          temperature: model.temperature ?? existing?.temperature ?? false,
          tool_call: model.tool_call ?? existing?.tool_call ?? true,
          cost:
            !model.cost && !existing?.cost
              ? {
                  input: 0,
                  output: 0,
                  cache_read: 0,
                  cache_write: 0,
                }
              : {
                  cache_read: 0,
                  cache_write: 0,
                  ...existing?.cost,
                  ...model.cost,
                },
          options: {
            ...existing?.options,
            ...model.options,
          },
          limit: model.limit ??
            existing?.limit ?? {
              context: 0,
              output: 0,
            },
        }
        parsed.models[modelID] = parsedModel
      }
      database[providerID] = parsed
    }

    const disabled = await Config.get().then((cfg) => new Set(cfg.disabled_providers ?? []))
    // load env
    for (const [providerID, provider] of Object.entries(database)) {
      if (disabled.has(providerID)) continue
      const apiKey = provider.env.map((item) => process.env[item]).at(0)
      if (!apiKey) continue
      mergeProvider(
        providerID,
        // only include apiKey if there's only one potential option
        provider.env.length === 1 ? { apiKey } : {},
        "env",
      )
    }

    // load apikeys
    for (const [providerID, provider] of Object.entries(await Auth.all())) {
      if (disabled.has(providerID)) continue
      if (provider.type === "api") {
        mergeProvider(providerID, { apiKey: provider.key }, "api")
      }
    }

    // load custom
    for (const [providerID, fn] of Object.entries(CUSTOM_LOADERS)) {
      if (disabled.has(providerID)) continue
      const result = await fn(database[providerID])
      if (result && (result.autoload || providers[providerID])) {
        mergeProvider(providerID, result.options ?? {}, "custom", result.getModel)
      }
    }

    // load config
    for (const [providerID, provider] of configProviders) {
      mergeProvider(providerID, provider.options ?? {}, "config")
    }

    for (const [providerID, provider] of Object.entries(providers)) {
      if (Object.keys(provider.info.models).length === 0) {
        delete providers[providerID]
        continue
      }
      log.info("found", { providerID })
    }

    return {
      models,
      providers,
      sdk,
    }
  })

  export async function list() {
    return state().then((state) => state.providers)
  }

  async function getSDK(provider: ModelsDev.Provider) {
    return (async () => {
      using _ = log.time("getSDK", {
        providerID: provider.id,
      })
      const s = await state()
      const existing = s.sdk.get(provider.id)
      if (existing) return existing
      const pkg = provider.npm ?? provider.id
      const mod = await import(await BunProc.install(pkg, "beta"))
      const fn = mod[Object.keys(mod).find((key) => key.startsWith("create"))!]
      const loaded = fn({
        name: provider.id,
        ...s.providers[provider.id]?.options,
      })
      s.sdk.set(provider.id, loaded)
      return loaded as SDK
    })().catch((e) => {
      throw new InitError({ providerID: provider.id }, { cause: e })
    })
  }

  export async function getModel(providerID: string, modelID: string) {
    const key = `${providerID}/${modelID}`
    const s = await state()
    if (s.models.has(key)) return s.models.get(key)!

    log.info("getModel", {
      providerID,
      modelID,
    })

    const provider = s.providers[providerID]
    if (!provider) throw new ModelNotFoundError({ providerID, modelID })
    const info = provider.info.models[modelID]
    if (!info) throw new ModelNotFoundError({ providerID, modelID })
    const sdk = await getSDK(provider.info)

    try {
      const language = provider.getModel ? await provider.getModel(sdk, modelID) : sdk.languageModel(modelID)
      log.info("found", { providerID, modelID })
      s.models.set(key, {
        info,
        language,
      })
      return {
        info,
        language,
      }
    } catch (e) {
      if (e instanceof NoSuchModelError)
        throw new ModelNotFoundError(
          {
            modelID: modelID,
            providerID,
          },
          { cause: e },
        )
      throw e
    }
  }

  export async function getSmallModel(providerID: string) {
    const cfg = await Config.get()

    if (cfg.small_model) {
      const parsed = parseModel(cfg.small_model)
      return getModel(parsed.providerID, parsed.modelID)
    }

    const provider = await state().then((state) => state.providers[providerID])
    if (!provider) return
    const priority = ["3-5-haiku", "3.5-haiku", "gemini-2.5-flash"]
    for (const item of priority) {
      for (const model of Object.keys(provider.info.models)) {
        if (model.includes(item)) return getModel(providerID, model)
      }
    }
  }

  const priority = ["gemini-2.5-pro-preview", "codex-mini", "claude-sonnet-4"]
  export function sort(models: ModelsDev.Model[]) {
    return sortBy(
      models,
      [(model) => priority.findIndex((filter) => model.id.includes(filter)), "desc"],
      [(model) => (model.id.includes("latest") ? 0 : 1), "asc"],
      [(model) => model.id, "desc"],
    )
  }

  export async function defaultModel() {
    const cfg = await Config.get()
    if (cfg.model) return parseModel(cfg.model)
    const provider = await list()
      .then((val) => Object.values(val))
      .then((x) => x.find((p) => !cfg.provider || Object.keys(cfg.provider).includes(p.info.id)))
    if (!provider) throw new Error("no providers found")
    const [model] = sort(Object.values(provider.info.models))
    if (!model) throw new Error("no models found")
    return {
      providerID: provider.info.id,
      modelID: model.id,
    }
  }

  export function parseModel(model: string) {
    const [providerID, ...rest] = model.split("/")
    return {
      providerID: providerID,
      modelID: rest.join("/"),
    }
  }

  const TOOLS = [
    BashTool,
    CredentialBashTool,
    EditTool,
    WebFetchTool,
    GlobTool,
    GrepTool,
    ListTool,
    // LspDiagnosticTool,
    // LspHoverTool,
    PatchTool,
    ReadTool,
    // MultiEditTool,
    WriteTool,
    TodoWriteTool,
    TodoReadTool,
    TaskTool,
  ]

  const TOOL_MAPPING: Record<string, Tool.Info[]> = {
    anthropic: TOOLS.filter((t) => t.id !== "patch"),
    openai: TOOLS.map((t) => ({
      ...t,
      parameters: optionalToNullable(t.parameters),
    })),
    azure: TOOLS.map((t) => ({
      ...t,
      parameters: optionalToNullable(t.parameters),
    })),
    google: TOOLS.map((t) => ({
      ...t,
      parameters: sanitizeGeminiParameters(t.parameters),
    })),
  }

  export async function tools(providerID: string) {
    /*
    const cfg = await Config.get()
    if (cfg.tool?.provider?.[providerID])
      return cfg.tool.provider[providerID].map(
        (id) => TOOLS.find((t) => t.id === id)!,
      )
        */
    return TOOL_MAPPING[providerID] ?? TOOLS
  }

  function sanitizeGeminiParameters(schema: z.ZodTypeAny, visited = new Set()): z.ZodTypeAny {
    if (!schema || visited.has(schema)) {
      return schema
    }
    visited.add(schema)

    if (schema instanceof z.ZodDefault) {
      const innerSchema = schema.removeDefault()
      // Handle Gemini's incompatibility with `default` on `anyOf` (unions).
      if (innerSchema instanceof z.ZodUnion) {
        // The schema was `z.union(...).default(...)`, which is not allowed.
        // We strip the default and return the sanitized union.
        return sanitizeGeminiParameters(innerSchema, visited)
      }
      // Otherwise, the default is on a regular type, which is allowed.
      // We recurse on the inner type and then re-apply the default.
      return sanitizeGeminiParameters(innerSchema, visited).default(schema._def.defaultValue())
    }

    if (schema instanceof z.ZodOptional) {
      return z.optional(sanitizeGeminiParameters(schema.unwrap(), visited))
    }

    if (schema instanceof z.ZodObject) {
      const newShape: Record<string, z.ZodTypeAny> = {}
      for (const [key, value] of Object.entries(schema.shape)) {
        newShape[key] = sanitizeGeminiParameters(value as z.ZodTypeAny, visited)
      }
      return z.object(newShape)
    }

    if (schema instanceof z.ZodArray) {
      return z.array(sanitizeGeminiParameters(schema.element, visited))
    }

    if (schema instanceof z.ZodUnion) {
      // This schema corresponds to `anyOf` in JSON Schema.
      // We recursively sanitize each option in the union.
      const sanitizedOptions = schema.options.map((option: z.ZodTypeAny) => sanitizeGeminiParameters(option, visited))
      return z.union(sanitizedOptions as [z.ZodTypeAny, z.ZodTypeAny, ...z.ZodTypeAny[]])
    }

    if (schema instanceof z.ZodString) {
      const newSchema = z.string({ description: schema.description })
      const safeChecks = ["min", "max", "length", "regex", "startsWith", "endsWith", "includes", "trim"]
      // rome-ignore lint/suspicious/noExplicitAny: <explanation>
      ;(newSchema._def as any).checks = (schema._def as z.ZodStringDef).checks.filter((check) =>
        safeChecks.includes(check.kind),
      )
      return newSchema
    }

    return schema
  }
  function optionalToNullable(schema: z.ZodTypeAny): z.ZodTypeAny {
    if (schema instanceof z.ZodObject) {
      const shape = schema.shape
      const newShape: Record<string, z.ZodTypeAny> = {}

      for (const [key, value] of Object.entries(shape)) {
        const zodValue = value as z.ZodTypeAny
        if (zodValue instanceof z.ZodOptional) {
          newShape[key] = zodValue.unwrap().nullable()
        } else {
          newShape[key] = optionalToNullable(zodValue)
        }
      }

      return z.object(newShape)
    }

    if (schema instanceof z.ZodArray) {
      return z.array(optionalToNullable(schema.element))
    }

    if (schema instanceof z.ZodUnion) {
      return z.union(
        schema.options.map((option: z.ZodTypeAny) => optionalToNullable(option)) as [
          z.ZodTypeAny,
          z.ZodTypeAny,
          ...z.ZodTypeAny[],
        ],
      )
    }

    return schema
  }

  export const ModelNotFoundError = NamedError.create(
    "ProviderModelNotFoundError",
    z.object({
      providerID: z.string(),
      modelID: z.string(),
    }),
  )

  export const InitError = NamedError.create(
    "ProviderInitError",
    z.object({
      providerID: z.string(),
    }),
  )
}
