import { InlineKeyboard } from "grammy";
import type { User } from "../models/User";

export const backToStartKeyboard = new InlineKeyboard().text("Назад", "backToStart");

export function MakeStartMenu(user: User): [InlineKeyboard, string] {
  var start = new InlineKeyboard();

  var text: string = "";
  if (user.FullName === "" && user.Birthdate === "") {
    text = `Привет ${user.Name}, давай знакомиться!`;
  } else {
    text = `Полное имя: ${user.FullName}\nДата рождения: ${user.Birthdate}\nБаланс: ${user.Balance}\nИмя для обращений: ${user.Name}\nФамилия: ${user.Surname}`;
  }

  if (user.FullName === "") {
    start.text("Ввести полное имя", "State=Register.FullName");
  }
  if (user.Birthdate === "") {
    start.text("Ввести дату рождения", "State=Register.Birthdate");
  }
  start.row()
  

  start.url("Связь с автором", "https://t.me/Nadya_Green");
  start.url("Связь с разрабом", "https://t.me/iamsoloma");


  return [start, text];
}
