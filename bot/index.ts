import { Bot, GrammyError, HttpError, InlineKeyboard } from "grammy";
import { Cron } from "croner";
import {
  GetListOfSubscribers,
  GetUserByID,
  RegisterNewUser,
  UpdateUser,
  type User,
} from "./models/User";
import { backToStartKeyboard, MakeStartMenu } from "./menus/start";
import { delay, formatDateToSurreal, parseDate } from "./utils";
import { GetCommonDayText, GetConscienceText } from "./models/Recomendations";
require("dotenv").config();

if (process.env.BOT_API_TOKEN === undefined) {
  throw new Error("BOT_API_TOKEN is empty");
}
if (process.env.APIEndpoint === undefined) {
  throw new Error("APIEndpoint is empty");
}

const bot = new Bot(process.env.BOT_API_TOKEN!);

bot.api.setMyCommands([
  { command: "start", description: "Запуск" },
  //{ command: "menu", description: "Получить меню" },
  {
    command: "conscience",
    description: "Получить рекомендацию на основе числа сознания",
  },
]);

bot.command("start", async (ctx) => {
  const userResp = await GetUserByID(String("tg" + ctx.message?.from.id));
  if (userResp.result === "error" && userResp.error !== "404") {
    console.log("Error getting user by Telegram ID: ", userResp.error);
    await ctx.reply(
      "Произошла ошибка при получении пользователя. Попробуй позже.",
      {
        reply_parameters: { message_id: ctx.message?.message_id! },
      }
    );
  } else if (userResp.result === "error" && userResp.error === "404") {
    const user: User = {
      id: { Table: "Users", ID: "tg" + ctx.message?.from.id },
      State: {},
      UserName: ctx.message?.from.username!,
      Name: ctx.message?.from.first_name!,
      Surname: ctx.message?.from.last_name!,
      LanguageCode: ctx.message?.from.language_code!,
      FullName: "",
      Subscriber: true,
      Email: "",
      Birthdate: "",
      Balance: 0,
    };

    const resp = await RegisterNewUser(user);
    if (resp.result === "error") {
      console.log("Error registering new user: " + resp.error);
      await ctx.reply("Произошла ошибка при регистрации. Попробуй позже.", {
        reply_parameters: { message_id: ctx.message?.message_id! },
      });
    } else {
      const menu = MakeStartMenu(user);
      await ctx.reply(menu[1], {
        reply_markup: menu[0],
        reply_parameters: { message_id: ctx.message?.message_id! },
      });
    }
  } else if (userResp.result === "success") {
    const user: User = userResp.value.user;
    const menu = MakeStartMenu(userResp.value.user);
    await ctx.reply(menu[1], {
      reply_markup: menu[0],
      reply_parameters: { message_id: ctx.message?.message_id! },
    });
  }
  /*await ctx.reply("Привет!", {
    reply_parameters: {message_id: ctx.msg.message_id}
  });*/
});

bot.command("conscience", async (ctx) => {
  const userResp = await GetUserByID(String("tg" + ctx.message?.from.id));
  if (userResp.result === "error" && userResp.error !== "404") {
    console.log("Error getting user by Telegram ID: " + userResp.error);
    ctx.reply("Произошла ошибка при получении пользователя. Попробуй позже.", {
      reply_parameters: { message_id: ctx.message?.message_id! },
    });
  } else if (userResp.result === "error" && userResp.error === "404") {
    ctx.reply(
      "Похоже, что ты ещё не зарегистрирован. Напиши /start, чтобы начать.",
      { reply_parameters: { message_id: ctx.message?.message_id! } }
    );
  } else if (userResp.result === "success") {
    const resp = await GetConscienceText(userResp.value.user.id.ID);
    if (resp.result === "error" && resp.error !== "404") {
      console.log("Error getting consciousness from storage: " + resp.error);
      ctx.reply(
        "Произошла ошибка при получении рекомендаций. Попробуй позже.",
        { reply_parameters: { message_id: ctx.message?.message_id! } }
      );
    } else if (resp.result === "error" && resp.error === "404") {
      console.log("Consciousness not found in storage ");
      ctx.reply("Произошла ошибка при поиске рекомендаций. Попробуй позже.", {
        reply_parameters: { message_id: ctx.message?.message_id! },
      });
    } else if (resp.result === "success") {
      ctx.reply(resp.value.text, {
        reply_parameters: { message_id: ctx.message?.message_id! },
        reply_markup: { remove_keyboard: true },
        parse_mode: "HTML",
      });
    }
  }
});

bot.callbackQuery(
  [
    "State=Register.FullName",
    "State=Register.Birthdate",
    "State=Balance.Amount",
  ],
  async (ctx) => {
    const userResp = await GetUserByID(
      "tg" + String(ctx.callbackQuery.from.id)
    );
    if (userResp.result === "error") {
      console.log("Error getting user by Telegram ID: " + userResp.error);
      ctx.reply(
        "Произошла ошибка при получении пользователя. Попробуй позже.",
        {
          reply_parameters: {
            message_id: ctx.callbackQuery.message?.message_id!,
          },
        }
      );
    } else {
      if (ctx.callbackQuery.data === "State=Register.FullName") {
        userResp.value.user.State["Register"] = "FullName";
      }
      if (ctx.callbackQuery.data === "State=Register.Birthdate") {
        userResp.value.user.State["Register"] = "Birthdate";
      }
      if (ctx.callbackQuery.data === "State=Balance.Amount") {
        userResp.value.user.State["Balance"] = "Amount";
      }

      const updateResp = await UpdateUser(userResp.value.user);
      if (updateResp.result === "error") {
        console.log("Error updating user state: " + updateResp.error);
        ctx.reply(
          "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
          {
            reply_parameters: {
              message_id: ctx.callbackQuery.message?.message_id!,
            },
          }
        );
      } else {
        var msg: string = "";
        if (ctx.callbackQuery.data === "State=Register.FullName") {
          msg =
            "Введи свое полное имя латинскими буквами как в загране или банковской карте(Nadezda, Vitaliy):";
        }
        if (ctx.callbackQuery.data === "State=Register.Birthdate") {
          msg = "Введи свою дату рождения в формате ДД.ММ.ГГГГ";
        }
        if (ctx.callbackQuery.data === "State=Balance.Amount") {
          msg = "Введи количество звёзд для пополнения(1 🌟 ~ 1.5 рубля):";
        }
        await ctx.api.editMessageText(
          ctx.chatId!,
          ctx.callbackQuery.message?.message_id!,
          msg,
          { reply_markup: backToStartKeyboard }
        );
      }
    }
  }
);

bot.callbackQuery(["backToStart"], async (ctx) => {
  const userResp = await GetUserByID("tg" + String(ctx.callbackQuery.from.id));
  if (userResp.result === "error") {
    console.log("Error getting user by Telegram ID: " + userResp.error);
    ctx.reply("Произошла ошибка при получении пользователя. Попробуй позже.", {
      reply_parameters: {
        message_id: ctx.callbackQuery.message?.message_id!,
      },
    });
  } else {
    if (ctx.callbackQuery.data === "backToStart") {
      userResp.value.user.State["Register"] = "";
      userResp.value.user.State["Balance"] = "";
    }

    const updateResp = await UpdateUser(userResp.value.user);
    if (updateResp.result === "error") {
      console.log("Error updating user state: " + updateResp.error);
      ctx.reply(
        "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
        {
          reply_parameters: {
            message_id: ctx.callbackQuery.message?.message_id!,
          },
        }
      );
    } else {
      const menu = MakeStartMenu(userResp.value.user);
      await ctx.api.editMessageText(
        ctx.chatId!,
        ctx.callbackQuery.message?.message_id!,
        menu[1],
        { reply_markup: menu[0] }
      );
    }
  }
});

bot.callbackQuery(["ChangeSubscription"], async (ctx) => {
  const userResp = await GetUserByID("tg" + String(ctx.callbackQuery.from.id));
  if (userResp.result === "error") {
    console.log("Error getting user by Telegram ID: " + userResp.error);
    ctx.reply("Произошла ошибка при получении пользователя. Попробуй позже.", {
      reply_parameters: {
        message_id: ctx.callbackQuery.message?.message_id!,
      },
    });
  } else {
    if (ctx.callbackQuery.data === "ChangeSubscription") {
      userResp.value.user.State["Register"] = "";
      if (userResp.value.user.Subscriber === true) {
        userResp.value.user.Subscriber = false;
      } else {
        userResp.value.user.Subscriber = true;
      }
    }

    const updateResp = await UpdateUser(userResp.value.user);
    if (updateResp.result === "error") {
      console.log("Error updating user state: " + updateResp.error);
      ctx.reply(
        "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
        {
          reply_parameters: {
            message_id: ctx.callbackQuery.message?.message_id!,
          },
        }
      );
    } else {
      const menu = MakeStartMenu(userResp.value.user);
      console.log(menu[1]);
      await ctx.api.editMessageText(
        ctx.chatId!,
        ctx.callbackQuery.message?.message_id!,
        menu[1],
        { reply_markup: menu[0] }
      );
    }
  }
});

bot.on("pre_checkout_query", (ctx) => {
  return ctx.answerPreCheckoutQuery(true).catch(() => {
    console.error("answerPreCheckoutQuery failed");
  });
});

bot.on("message:successful_payment", async (ctx) => {
  if (!ctx.message || !ctx.message.successful_payment || !ctx.from) {
    return;
  }
  const userResp = await GetUserByID(String("tg" + ctx.message.from.id));
  if (userResp.result === "error" && userResp.error !== "404") {
    console.log("Error getting user by Telegram ID: " + userResp.error);
    ctx.reply("Произошла ошибка при получении пользователя. Попробуй позже.", {
      reply_parameters: { message_id: ctx.message.message_id },
    });
  } else if (userResp.result === "error" && userResp.error === "404") {
    ctx.reply(
      "Похоже, что ты ещё не зарегистрирован. Напиши /start, чтобы начать.",
      { reply_parameters: { message_id: ctx.message.message_id } }
    );
  } else if (userResp.result === "success") {
    userResp.value.user.State["Balance"] = "";
    userResp.value.user.Balance +=
      ctx.message.successful_payment.total_amount * 1.5;
    const updateResp = await UpdateUser(userResp.value.user);
    if (updateResp.result === "error") {
      console.log("Error updating user state: " + updateResp.error);
      ctx.reply(
        "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
        {
          reply_parameters: {
            message_id: ctx.message.message_id,
          },
        }
      );
    } else {
      const msg = await ctx.reply(
        "Баланс успешно пополнен! ID транзакции: " +
          ctx.message.successful_payment.telegram_payment_charge_id,
        {
          reply_parameters: { message_id: ctx.message.message_id },
        }
      );
    }
  }
});

bot.command("refund", async (ctx) => {
  const transactionID = ctx.match;

  bot.api
    .refundStarPayment(ctx.from?.id!, transactionID)
    .then(async () => {
      const userResp = await GetUserByID(String("tg" + ctx.message?.from.id!));
      if (userResp.result === "error" && userResp.error !== "404") {
        console.log("Error getting user by Telegram ID: " + userResp.error);
        ctx.reply(
          "Произошла ошибка при получении пользователя. Попробуй позже.",
          {
            reply_parameters: { message_id: ctx.message?.message_id! },
          }
        );
      } else if (userResp.result === "error" && userResp.error === "404") {
        ctx.reply(
          "Похоже, что ты ещё не зарегистрирован. Напиши /start, чтобы начать.",
          { reply_parameters: { message_id: ctx.message?.message_id! } }
        );
      } else if (userResp.result === "success") {
        const trList = await ctx.api.getStarTransactions();
        var amount:number = 0.0
        for (let t of trList.transactions) {
          if (transactionID === t.id) {
            amount = t.amount
          }
        }
        userResp.value.user.Balance -= amount * 1.5;
        const updateResp = await UpdateUser(userResp.value.user);
        if (updateResp.result === "error") {
          console.log("Error updating user state: " + updateResp.error);
          ctx.reply(
            "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
            {
              reply_parameters: {
                message_id: ctx.message?.message_id!,
              },
            }
          );
        }
      }
      return ctx.reply("Успешный возврат");
    })
    .catch(() => ctx.reply("Refund failed"));
});

bot.on("message:text", async (ctx) => {
  const userResp = await GetUserByID(String("tg" + ctx.message?.from.id!));
  if (userResp.result === "error" && userResp.error !== "404") {
    console.log("Error getting user by Telegram ID: " + userResp.error);
    ctx.reply("Произошла ошибка при получении пользователя. Попробуй позже.", {
      reply_parameters: { message_id: ctx.message.message_id },
    });
  } else if (userResp.result === "error" && userResp.error === "404") {
    console.log("Message");
    ctx.reply(
      "Похоже, что ты ещё не зарегистрирован. Напиши /start, чтобы начать.",
      { reply_parameters: { message_id: ctx.message.message_id } }
    );
  } else if (userResp.result === "success") {
    if (userResp.value.user.State["Register"] === "FullName") {
      if (!/^[a-zA-Z]+$/.test(ctx.message.text!)) {
        console.log("Error parsing FullName in " + ctx.message.text!);
        ctx.reply(
          "Неверный формат имени. Пожалуйста, введи ТОЛЬКО латинским буквами",
          { reply_parameters: { message_id: ctx.message.message_id! } }
        );
      } else {
        userResp.value.user.FullName = ctx.message.text!;
        userResp.value.user.State["Register"] = "";
        const updateResp = await UpdateUser(userResp.value.user);
        if (updateResp.result === "error") {
          console.log("Error updating user state: " + updateResp.error);
          ctx.reply(
            "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
            {
              reply_parameters: {
                message_id: ctx.message.message_id,
              },
            }
          );
        } else {
          const msg = await ctx.reply("Твоё полное имя успешно сохранено.", {
            reply_parameters: { message_id: ctx.message.message_id },
          });
          await delay(2500); //2.5 in seconds
          ctx.deleteMessage();
          bot.api.deleteMessage(ctx.message.from.id, msg.message_id);
        }
      }
    } else if (userResp.value.user.State["Register"] === "Birthdate") {
      const date = parseDate(ctx.message.text!);
      if (!date) {
        console.log("Error parsing birthdate in " + ctx.message.text!);
        ctx.reply(
          "Неверный формат даты. Пожалуйста, введи дату в формате ДД.ММ.ГГГГ",
          { reply_parameters: { message_id: ctx.message.message_id! } }
        );
      } else {
        userResp.value.user.Birthdate = formatDateToSurreal(date);
        userResp.value.user.State["Register"] = "";
        const updateResp = await UpdateUser(userResp.value.user);
        if (updateResp.result === "error") {
          console.log("Error updating user state: " + updateResp.error);
          ctx.reply(
            "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
            {
              reply_parameters: {
                message_id: ctx.message.message_id,
              },
            }
          );
        } else {
          const msg = await ctx.reply("Твоя дата рождения успешно сохранена.", {
            reply_parameters: { message_id: ctx.message.message_id },
          });
          await delay(2500); //2.5 in seconds
          ctx.deleteMessage();
          bot.api.deleteMessage(ctx.message.from.id, msg.message_id);
        }
      }
    } else if ((userResp.value.user.State["Balance"] = "Amount")) {
      const Amount = Number(ctx.message.text!);
      if (!Amount || Amount < 1) {
        console.log("Error parsing amount in " + ctx.message.text!);
        ctx.reply("Неверно указано целое число звёзд.", {
          reply_parameters: { message_id: ctx.message.message_id! },
        });
      } else {
        userResp.value.user.State["Balance"] = "Processing";
        const updateResp = await UpdateUser(userResp.value.user);
        if (updateResp.result === "error") {
          console.log("Error updating user state: " + updateResp.error);
          ctx.reply(
            "Произошла ошибка при обновлении состояния пользователя. Попробуй позже.",
            {
              reply_parameters: {
                message_id: ctx.message.message_id,
              },
            }
          );
        } else {
          const msg = await ctx.replyWithInvoice(
            "Пополнение баланса",
            "Пополнение балнса на Digits Say",
            "{}",
            "XTR",
            [{ amount: Amount, label: "Пополнение баланса" }]
          );
        }
      }
    }
  }
});

const DailyMail = new Cron("0 0 5 * * *", sendDailyMessage);

async function sendDailyMessage() {
  var users: User[] = [];
  var text: string = "";

  var reqtext = await GetCommonDayText();
  if (reqtext.result === "error") {
    console.log(reqtext.error);
    //throw Error(reqtext.error);
  } else {
    text = reqtext.value.text;
  }

  const res = await GetListOfSubscribers();
  if (res.result === "error") {
    console.log(res.error);
  } else {
    users = res.value.users;
  }

  //console.log(JSON.stringify(users, null, 2))

  for (var user of users) {
    try {
      await bot.api.sendMessage(user.id.ID.replace("tg", ""), text, {
        parse_mode: "HTML",
      });
    } catch (err) {
      if (err instanceof HttpError) {
        console.error(`Could not connect to Telegram: `, err.message);
      } else if (err instanceof GrammyError) {
        console.log(
          `Can not send daily mail to user with id=${user.id.ID}: ${err.message}`
        );
      } else {
        console.error("Unknown error: ", err);
      }
    }
    //bot.api.sendMessage(user.id.ID.replace("tg", ""), text);
  }
}

bot.catch((err) => {
  const ctx = err.ctx;
  console.error(`Error while handing update ${ctx.update.update_id}: `);
  const e = err.error;

  if (e instanceof GrammyError) {
    console.error(`Error in request: `, e.description);
  } else if (e instanceof HttpError) {
    console.error(`Could not connect to Telegram: `, e);
  } else {
    console.error("Unknown error: ", e);
  }
});

await bot.init();
console.log("Online as a @" + bot.botInfo.username);
bot.start();
