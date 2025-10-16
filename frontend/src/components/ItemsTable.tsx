// import React from "react";
// import type { Item } from "./types";

// type Props = {
//   items: Item[];
//   onChange: (index: number, key: keyof Item, value: Item[keyof Item]) => void;
//   onAdd: () => void;
//   onDelete: (index: number) => void;
// };

// export default function ItemsTable({ items, onChange, onAdd, onDelete }: Props) {
//   return (
//     <>
//       <table style={{ width: "100%", borderCollapse: "collapse", marginTop: 8 }}>
//         <thead>
//           <tr>
//             <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Name</th>
//             <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Unit Price (LKR)</th>
//             <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #ddd" }}>Quantity</th>
//             <th style={{ width: 80 }} />
//           </tr>
//         </thead>
//         <tbody>
//           {items.map((it, i) => (
//             <tr key={i}>
//               <td style={{ padding: 8 }}>
//                 <input
//                   value={it.name}
//                   placeholder="e.g., Notebook"
//                   onChange={(e) => onChange(i, "name", e.target.value)}
//                 />
//               </td>
//               <td style={{ padding: 8 }}>
//                 <input
//                   type="number"
//                   step="0.01"
//                   min="0"
//                   value={it.price}
//                   onChange={(e) => onChange(i, "price", Number(e.target.value))}
//                 />
//               </td>
//               <td style={{ padding: 8 }}>
//                 <input
//                   type="number"
//                   step="1"
//                   min="1"
//                   value={it.quantity}
//                   onChange={(e) =>
//                     onChange(i, "quantity", Math.max(1, Math.trunc(Number(e.target.value))))
//                   }
//                 />
//               </td>
//               <td style={{ padding: 8 }}>
//                 {items.length > 1 && <button onClick={() => onDelete(i)}>Delete</button>}
//               </td>
//             </tr>
//           ))}
//         </tbody>
//       </table>

//       <div style={{ display: "flex", gap: 12, marginTop: 12 }}>
//         <button onClick={onAdd}>Add Row</button>
//       </div>
//     </>
//   );
// }

// src/components/ItemsTable.tsx
import React from "react";
import type { Item } from "./types";
import { isValid } from "./validators";

type Props = {
  items: Item[];
  onChange: (index: number, key: keyof Item, value: Item[keyof Item]) => void;
  onAdd: () => void;
  onDelete: (index: number) => void;
};

export default function ItemsTable({ items, onChange, onAdd, onDelete }: Props) {
  return (
    <>
      <div className="overflow-x-auto">
        <table className="w-full border-separate border-spacing-0">
          <thead>
            <tr className="text-left text-xs uppercase tracking-wide text-white/60">
              <th className="border-b border-white/10 px-3 py-2">Name</th>
              <th className="border-b border-white/10 px-3 py-2">Unit Price (LKR)</th>
              <th className="border-b border-white/10 px-3 py-2">Quantity</th>
              <th className="border-b border-white/10 px-3 py-2 w-24"></th>
            </tr>
          </thead>
          <tbody>
            {items.map((it, i) => {
              const valid = isValid(it);
              return (
                <tr key={i} className="align-top">
                  <td className="px-3 py-2">
                    <input
                      className={`w-full rounded-lg bg-white/[.06] px-3 py-2 outline-none
                        placeholder:text-white/40 border
                        ${it.name.trim().length === 0 ? "border-red-400/40" : "border-white/10"}
                        focus:border-[#00E5FF]/50`}
                      value={it.name}
                      placeholder="e.g., Notebook"
                      onChange={(e) => onChange(i, "name", e.target.value)}
                    />
                  </td>
                  <td className="px-3 py-2">
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      className={`w-40 rounded-lg bg-white/[.06] px-3 py-2 outline-none border
                        ${!(Number.isFinite(it.price) && it.price >= 0) ? "border-red-400/40" : "border-white/10"}
                        focus:border-[#00E5FF]/50`}
                      value={it.price}
                      onChange={(e) =>
                        onChange(i, "price", Number(e.target.value))
                      }
                    />
                  </td>
                  <td className="px-3 py-2">
                    <input
                      type="number"
                      step="1"
                      min="1"
                      className={`w-28 rounded-lg bg-white/[.06] px-3 py-2 outline-none border
                        ${it.quantity < 1 ? "border-red-400/40" : "border-white/10"}
                        focus:border-[#00E5FF]/50`}
                      value={it.quantity}
                      onChange={(e) =>
                        onChange(
                          i,
                          "quantity",
                          Math.max(1, Math.trunc(Number(e.target.value)))
                        )
                      }
                    />
                    {!valid && (
                      <div className="mt-1 text-[10px] text-red-300/80">
                        Name required, price ≥ 0, quantity ≥ 1.
                      </div>
                    )}
                  </td>
                  <td className="px-3 py-2">
                    {items.length > 1 && (
                      <button
                        onClick={() => onDelete(i)}
                        className="rounded-lg border border-white/15 px-3 py-2 text-xs text-white/80 hover:bg-white/5"
                      >
                        Delete
                      </button>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <div className="mt-4">
        <button
          onClick={onAdd}
          className="rounded-lg border border-white/15 px-3 py-2 text-sm text-white/80 hover:bg-white/5"
        >
          Add Row
        </button>
      </div>
    </>
  );
}
