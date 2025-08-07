declare module "bun" {
    interface Env {
        TGAPIToken: string
        APIEndpoint: string
    }
}

if (process.env.TGAPIToken == "" ||process.env.TGAPIToken == undefined) {
    throw console.error("Empty TGAPIToken .env");
}
if (process.env.APIEndpoint == "" ||process.env.APIEndpoint == undefined) {
    throw console.error("Empty APIEndpoint in .env");
}

import {Bot} from "grammy";

const bot = new Bot(process.env.TGAPIToken)

bot.command("start", ctx => ctx.reply("HEllo!"))

bot.start()