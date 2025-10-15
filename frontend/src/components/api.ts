import { postJSON } from "../lib/http";
import type { Item, Stats } from "./types";

export function saveAndCalculate(items: Item[]) {
  return postJSON<Stats>("/api/items/bulk", { items });
}
