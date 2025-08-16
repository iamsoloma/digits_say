require("dotenv").config();

const {
  Bot,
  GrammyError,
  HttpError,
  InlineKeyboard,
} = require("grammy");

//const { hydrate } = require("@grammyjs/hydrate");

const bot = new Bot(process.env.BOT_API_TOKEN);
//bot.use(hydrate());

bot.api.setMyCommands([
  { command: "start", description: "Запуск" },
  { command: "menu", description: "Получить меню" },
]);

bot.command("start", async (ctx) => {
  /*await ctx.reply("Привет!", {
    reply_parameters: {message_id: ctx.msg.message_id}
  });*/
  
  await ctx.react("👀");
});

const menuKeyboard = new InlineKeyboard()
  .text("Узнать статус заказа", "order-status")
  .text("Обратиться в поддержку", "support");

const backKeyboard = new InlineKeyboard().text("Назад в меню", "back");

bot.command("menu", async (ctx) => {
  await ctx.reply("Выбери пункт меню:", { reply_markup: menuKeyboard });
});

bot.callbackQuery("order-status", async (ctx) => {
  await ctx.api.editMessageText(
    ctx.chatId,
    ctx.callbackQuery.message.message_id,
    "Статус заказа: в пути",
    { reply_markup: backKeyboard }
  );
});

bot.callbackQuery("support", async (ctx) => {
  await ctx.api.editMessageText(
    ctx.chatId,
    ctx.callbackQuery.message.message_id,
    "Напишите ваш запрос:",
    { reply_markup: backKeyboard }
  );
});


bot.callbackQuery("back", async (ctx) => {
  await ctx.api.editMessageText(
    ctx.chatId,
    ctx.callbackQuery.message.message_id,
    "Выбери пункт меню:",
    { reply_markup: menuKeyboard }
  );
});

bot.catch((err) => {
  const ctx = err.ctx;
  console.error(
    `Error while handing update ${ctx.update.update_id}: `
  );
  const e = err.error;

  if (e instanceof GrammyError) {
    console.error(`Error in request: `, e.description);
  } else if (e instanceof HttpError) {
    console.error(`Could not connect to Telegram: `, e);
  } else {
    console.error("Unknown error: ", e);
  }
});

bot.start();
