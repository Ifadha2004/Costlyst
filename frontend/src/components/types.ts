export type Item = { name: string; price: number; quantity: number };

export type Stats = {
  lineItemCount: number;
  totalQuantity: number;
  totalCost: number;
  avgUnitPrice: number;
  avgLineCost: number;
};
