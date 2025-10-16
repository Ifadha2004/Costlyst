// src/components/api.ts
import { postJSON, getJSON } from "../lib/http";
import type { Item, Stats } from "./types";

export type BulkResp = {
  message: string;
  batch: Stats;
  global: Stats;
};

// Be lenient about shapes: {message,batch,global} OR {message,data:{batch,global}} OR {globalStats:...}
export async function saveAndCalculate(items: Item[]): Promise<BulkResp> {
  const raw = await postJSON<any>("/api/items/bulk", { items });

  const batch =
    raw?.batch ??
    raw?.data?.batch ??
    raw?.stats?.batch;

  const global =
    raw?.global ??
    raw?.globalStats ??
    raw?.data?.global ??
    raw?.stats?.global;

  const message = raw?.message ?? "saved";

  if (!batch || !global) {
    // Helpful during dev if backend contract doesn't match
    // eslint-disable-next-line no-console
    console.warn("Unexpected /api/items/bulk response shape:", raw);
    throw new Error("Unexpected server response. Check bulk payload keys.");
  }

  return { message, batch, global };
}

// (optional helpers if you use them elsewhere)
export async function listItems(limit = 100) {
  return getJSON(`/api/items?limit=${limit}`);
}
