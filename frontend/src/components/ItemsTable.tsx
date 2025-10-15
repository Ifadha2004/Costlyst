import React from "react";
import type { Item } from "./types";

type Props = {
  items: Item[];
  onChange: (index: number, key: keyof Item, value: Item[keyof Item]) => void;
  onAdd: () => void;
  onDelete: (index: number) => void;
};

export default function ItemsTable({ items, onChange, onAdd, onDelete }: Props) {
  return (
    <>
      <table style={{ width: "100%", borderCollapse: "collapse", marginTop: 8 }}>
        <thead>
          <tr>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Name</th>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Unit Price (LKR)</th>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Quantity</th>
            <th style={{ width: 80 }} />
          </tr>
        </thead>
        <tbody>
          {items.map((it, i) => (
            <tr key={i}>
              <td style={{ padding: 8 }}>
                <input
                  value={it.name}
                  placeholder="e.g., Notebook"
                  onChange={(e) => onChange(i, "name", e.target.value)}
                />
              </td>
              <td style={{ padding: 8 }}>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={it.price}
                  onChange={(e) => onChange(i, "price", Number(e.target.value))}
                />
              </td>
              <td style={{ padding: 8 }}>
                <input
                  type="number"
                  step="1"
                  min="1"
                  value={it.quantity}
                  onChange={(e) =>
                    onChange(i, "quantity", Math.max(1, Math.trunc(Number(e.target.value))))
                  }
                />
              </td>
              <td style={{ padding: 8 }}>
                {items.length > 1 && <button onClick={() => onDelete(i)}>Delete</button>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div style={{ display: "flex", gap: 12, marginTop: 12 }}>
        <button onClick={onAdd}>Add Row</button>
      </div>
    </>
  );
}
