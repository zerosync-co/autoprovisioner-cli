import { UI } from "../ui"
import * as prompts from "@clack/prompts"
import open from "open"

export const UpgradeCommand = {
  command: "upgrade",
  describe: "upgrade plan for higher rate limits and more usage with the best models",
  handler: async () => {
    const pricingUrl = "https://autoprovisioner.zerosync.co/pricing"

    prompts.note("This will open a browser window to complete the subscription upgrade")
    const shouldOpen = await prompts.confirm({
      message: "Open subscription link in browser?",
    })

    if (prompts.isCancel(shouldOpen)) throw new UI.CancelledError()

    if (!shouldOpen) {
      prompts.log.info(`You can manually visit: ${pricingUrl}`)
    } else {
      open(pricingUrl)
    }

    prompts.outro("Done")
  },
}
