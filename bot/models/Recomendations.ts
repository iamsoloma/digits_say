import type { Result } from "../utils";

export async function GetCommonDayText(): Promise<
  Result<{ text: string }, string>
> {
  const response = await fetch(`${process.env.APIEndpoint}/commonday`);
  if (!response.ok) {
    return {
      result: "error",
      error: `API response not ok: ${response.statusText}}`,
    };
  } else {
    const result = await response.body?.text()!;
    return { result: "success", value: { text: result } };
  }
}

export async function GetConscienceText(
  id: string
): Promise<Result<{ text: string }, string>> {
  const response = await fetch(
    `${process.env.APIEndpoint}/conscience?id=${id}`
  );
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
    const result = await response.body?.text()!;
    return { result: "success", value: { text: result } };
  }
}
