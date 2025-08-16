export type Result<T, E> =
  | { result: "success"; value: T }
  | { result: "error"; error: E };


export function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

export function parseDate(input: string): Date | null {
  // Expected input format: "dd.MM.yyyy"
  const parts = input.split(".");
  if (parts.length !== 3) return null;

  const day = parseInt(parts[0]!, 10);
  const month = parseInt(parts[1]!, 10) - 1; // Month is zero-based in JS Date
  const year = parseInt(parts[2]!, 10);

  if (
    isNaN(day) || isNaN(month) || isNaN(year) ||
    day < 1 || day > 31 ||
    month < 0 || month > 11 ||
    year < 1000
  ) {
    return null; // Invalid date parts
  }

  return new Date(year, month, day);
}

export function formatDateToSurreal(date: Date): string {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0'); // месяцы с 0
  const day = date.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
}
