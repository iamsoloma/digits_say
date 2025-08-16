import type { Result } from "../utils";

export type User = {
  id: { Table: string; ID: string };
  State: Record<string, any>;
  UserName: string;
  Name: string;
  Surname: string;
  FullName: string;
  Subscsriber: boolean;
  LanguageCode: string;
  Email: string;
  Birthdate: string;
  Balance: number;
};

export async function GetUserByID(
  id: string
): Promise<Result<{ user: User }, string>> {
  const response = await fetch(`${process.env.APIEndpoint}/user/?id=${id}`);
  if (response.status === 404) {
    return {
      result: "error",
      error: "404",
    };
  } else if (!response.ok) {
    return {
      result: "error",
      error: `API response not ok: ${response.statusText}`,
    };
  } else {
    const result = (await response.json()) as User;
    return { result: "success", value: { user: result } };
  }
}

export async function RegisterNewUser(user: User): Promise<Result<{}, string>> {
  const response = await fetch(`${process.env.APIEndpoint}/user`, {
    method: "POST",
    body: JSON.stringify(user),
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (!response.ok) {
    return {
      result: "error",
      error: `API response not ok: ${
        response.statusText
      }: ${await response.body?.text()}`,
    };
  } else {
    return { result: "success", value: {} };
  }
}

export async function UpdateUser(user: User): Promise<Result<string, string>> {
  const response = await fetch(`${process.env.APIEndpoint}/user`, {
    method: "PATCH",
    body: JSON.stringify(user),
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (!response.ok) {
    return {
      result: "error",
      error: response.statusText + ": " + response.body?.text,
    };
  } else {
    return {
      result: "success",
      value: "Ok",
    };
  }
}

export async function GetListOfSubscribers(): Promise<
  Result<{ users: Array<User> }, string>
> {
  const response = await fetch(`${process.env.APIEndpoint}/subscribers`);
  if (!response.ok) {
    return {
      result: "error",
      error: `API response not ok: ${response.statusText}`,
    };
  } else {
    const result = (await response.json()) as Array<User>;
    return { result: "success", value: { users: result } };
  }
}
