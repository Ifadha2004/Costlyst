// src/lib/http.ts
export const API_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080";

/**
 * Generic POST JSON helper
 */
export async function postJSON<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error((await res.text()) || "Request failed");
  return res.json() as Promise<T>;
}

/**
 * Generic GET JSON helper
 */
export async function getJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`);
  if (!res.ok) throw new Error((await res.text()) || "Request failed");
  return res.json() as Promise<T>;
}

/**
 * Generic DELETE JSON helper
 */
export async function delJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, { method: "DELETE" });
  if (!res.ok) throw new Error((await res.text()) || "Request failed");
  return res.json() as Promise<T>;
}
