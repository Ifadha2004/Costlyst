import { postJSON } from "../lib/http";
import type { Item, Stats } from "./types";

type BulkSaveResponse = {
  message: string;
  data: {
    batch: Stats;
    global: Stats;
  };
};

export async function saveAndCalculate(items: Item[]): Promise<Stats> {
  const response = await postJSON<BulkSaveResponse>("/api/items/bulk", { items });
  return response.data.global; // ✅ wrapped under data
}
