import { InlineKeyboard } from "grammy";
import type { User } from "../models/User";

export const backToStartKeyboard = new InlineKeyboard().text(
  "Назад",
  "backToStart"
);

export function MakeStartMenu(user: User): [InlineKeyboard, string] {
  var start = new InlineKeyboard();

  var text: string = "";
  if (user.FullName === "" && user.Birthdate === "") {
    text = `Привет ${user.Name}, давай знакомиться!`;
  } else {
    const birthdate =
      user.Birthdate.slice(8, 10) +
      "." +
      user.Birthdate.slice(5, 7) +
      "." +
      user.Birthdate.slice(0, 4);
    var sub: string;
    if (user.Subscriber) {
      sub = "Да";
    } else {
      sub = "Нет";
    }

    text = `Полное имя: ${user.FullName}\nДата рождения: ${birthdate}\nБаланс: ${user.Balance}\nИмя для обращений: ${user.Name}\nФамилия: ${user.Surname}\nПодписчик: ${sub}`;
  }

  if (user.FullName === "") {
    start.text("Ввести полное имя", "State=Register.FullName");
  }
  if (user.Birthdate === "") {
    start.text("Ввести дату рождения", "State=Register.Birthdate");
  }
  start.row();


  start.text("Пополнить баланс", "State=Balance.Amount")
  start.row()
  var msg: string;
  if (user.Subscriber) {
    msg = "Отписаться от рассылки";
  } else {
    msg = "Подписаться на ежедневную рассылку";
  }
  start.text(msg, "ChangeSubscription");
  start.row();

  start.url("Связь с автором", "https://t.me/Nadya_Green");
  start.url("Связь с разрабом", "https://t.me/iamsoloma");

  return [start, text];
}
