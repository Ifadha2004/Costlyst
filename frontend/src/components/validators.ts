import type { Item } from "./types";

export const normalize = (i: Item): Item => ({
  name: i.name.trim(),
  price: Number(i.price),
  quantity: Math.max(1, Math.trunc(Number(i.quantity))),
});

export const isValid = (i: Item) =>
  i.name.length > 0 && Number.isFinite(i.price) && i.price >= 0 && i.quantity >= 1;
