import { Cron } from "croner";
import { GetCommonDayText, GetListOfSubscribers, GetUserByID, RegisterNewUser, type User } from "./models/User";

new Cron("*/5 * * * * *", async () => {
  const date = new Date();

 /* const res = await GetUserByID("tg17400");
  if (res.result === "error") {
    console.log(res.error);
  } else {
    console.log(res.value);
  }*/

  const res1 = await GetListOfSubscribers();
  if (res1.result === "error") {
    console.log(res1.error);
  } else {
    console.log(JSON.stringify(res1.value.users, null, 2));
  }
  console.log(date);

  /*const newUser: User = {
      ID: { Table: "Users", ID: "tg1" },
      UserName: "test",
      State: {},
      Name: "",
      SurName: "",
      FullName: "",
      Subscsriber: false,
      LanguageCode: "",
      Email: "",
      Birthdate: "",
      Balance: 0
  }
  const res = await RegisterNewUser(newUser)
  if (res.result == "error") {
    console.log(res.error)
  } else{
    console.log(JSON.stringify(res.value, null, 2))
  }
  console.log(
    "--------------------------------------------------------------------------"
  );*/
  const res = await GetCommonDayText()
  if (res.result === "error") {
    console.log(res.error)
  } else {
    console.log(res.value.text)
  }
});
