import { App } from "../../app/app"
import { Config } from "../../config/config"
import { cmd } from "./cmd"

export const ModelsCommand = cmd({
  command: "models",
  describe: "list available models from zerosync config",
  handler: async () => {
    await App.provide({ cwd: process.cwd() }, async () => {
      const config = await Config.get()

      const zerosyncProvider = config.provider?.["zerosync"]
      if (zerosyncProvider?.models) {
        for (const modelID of Object.keys(zerosyncProvider.models)) {
          console.log(`zerosync/${modelID}`)
        }
      } else {
        console.log("No zerosync models found in config")
      }
    })
  },
})
