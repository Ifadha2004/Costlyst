import React, { useMemo, useState, useEffect } from "react";
import ItemsTable from "./ItemsTable";
import type { Item, Stats } from "./types";
import { normalize, isValid } from "./validators";
import { API_URL } from "../lib/http";
import { saveAndCalculate } from "./api"; // clearAll removed

// ---- Toast model (with id so each new toast resets the timer) ----
type Toast = { id: number; kind: "ok" | "err"; text: string };

// FE helper to compute preview stats
function computeStats(rows: Item[]): Stats {
  const valid = rows.filter(isValid);
  const lineItemCount = valid.length;
  const totalQuantity = valid.reduce((s, r) => s + r.quantity, 0);
  const totalCost = Number(
    valid.reduce((s, r) => s + r.quantity * r.price, 0).toFixed(2)
  );
  const avgUnitPrice =
    totalQuantity > 0 ? Number((totalCost / totalQuantity).toFixed(2)) : 0;
  const avgLineCost =
    lineItemCount > 0 ? Number((totalCost / lineItemCount).toFixed(2)) : 0;
  return { lineItemCount, totalQuantity, totalCost, avgUnitPrice, avgLineCost };
}

export default function ItemsPage() {
  const [items, setItems]   = useState<Item[]>([{ name: "", price: 0, quantity: 1 }]);
  const [batch, setBatch]   = useState<Stats | null>(null);
  const [global, setGlobal] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(false);
  const [toast, setToast]   = useState<Toast | null>(null);

  // Auto-dismiss toast ~2.6s after it's shown, and reset timer for each new toast
  useEffect(() => {
    if (!toast) return;
    const t = setTimeout(() => setToast(null), 2600);
    return () => clearTimeout(t);
  }, [toast?.id]);

  const addRow = () => setItems((p) => [...p, { name: "", price: 0, quantity: 1 }]);
  const delRow = (i: number) => setItems((p) => p.filter((_, idx) => idx !== i));
  const update = (i: number, k: keyof Item, v: Item[keyof Item]) =>
    setItems((p) => p.map((it, idx) => (idx === i ? { ...it, [k]: v } : it)));
  const clearDraft = () => {
    setItems([{ name: "", price: 0, quantity: 1 }]);
    setBatch(null);
    setGlobal(null);
    setToast(null);
    // keep `global` (it reflects saved DB). If you want to hide global too, uncomment:
    // setGlobal(null);
  };

  // Live preview (client computed)
  const preview = useMemo(() => computeStats(items.map(normalize)), [items]);

  const submit = async () => {
    setToast(null);
    setLoading(true);
    try {
      const payload = items.map(normalize).filter(isValid);
      if (payload.length === 0) {
        throw new Error("Add at least one valid item (name, price ≥ 0, quantity ≥ 1).");
      }
      const resp = await saveAndCalculate(payload);
      console.log("bulk resp", resp);
      setBatch(resp.batch);
      setGlobal(resp.global);
      setToast({ id: Date.now(), kind: "ok", text: "Items saved. Stats updated." });
    } catch (e: any) {
      setToast({ id: Date.now(), kind: "err", text: e?.message ?? "Something went wrong." });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-dvh bg-[#06060e] text-white/90">
      <div className="mx-auto max-w-5xl px-6 py-10">
        <header className="mb-6">
          <h1 className="text-2xl font-semibold tracking-tight">Items Cost Calculator</h1>
          <p className="mt-1 text-sm text-white/60">
            Enter name, unit price, and quantity. Live preview updates as you type.
            Save to persist and see global stats.
          </p>
        </header>

        {/* Draft table */}
        <div className="rounded-2xl border border-white/10 bg-white/[.04] backdrop-blur-xl p-4 md:p-5 shadow-[0_0_40px_rgba(122,92,255,.12)]">
          <ItemsTable items={items} onChange={update} onAdd={addRow} onDelete={delRow} />

          <div className="mt-5 flex flex-wrap items-center gap-3">
            <button
              onClick={submit}
              disabled={loading}
              className="inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-black font-medium
                         bg-gradient-to-r from-[#00E5FF] to-[#7A5CFF]
                         shadow-[0_0_20px_rgba(122,92,255,.35),0_0_30px_rgba(0,229,255,.25)]
                         hover:brightness-110 focus:outline-none focus:ring-2 focus:ring-[#00E5FF]/40 disabled:opacity-50"
            >
              {loading ? "Saving…" : "Save & Calculate"}
            </button>

            {/* moved right + red style */}
            <button
              onClick={clearDraft}
              className="ml-auto rounded-full border border-red-400/40 text-red-200/90 px-4 py-2 text-sm
                         hover:bg-red-500/10 transition"
              title="Clear all draft rows and reset preview (does not touch DB)"
            >
              Clear draft
            </button>
          </div>
        </div>

        {/* Stats */}
        <section className="mt-8 grid gap-6 md:grid-cols-2">
          <div className="rounded-2xl border border-white/10 bg-white/[.04] p-4 md:p-5">
            <h3 className="mb-3 text-sm font-medium text-white/70">Preview (unsaved)</h3>
            <div className="grid grid-cols-2 gap-3">
              <StatCard label="Line Items" value={preview.lineItemCount} />
              <StatCard label="Total Quantity" value={preview.totalQuantity} />
              <StatCard label="Total Cost (LKR)" value={preview.totalCost.toFixed(2)} />
              <StatCard label="Avg Unit Price (LKR)" value={preview.avgUnitPrice.toFixed(2)} />
              <StatCard label="Avg Line Cost (LKR)" value={preview.avgLineCost.toFixed(2)} />
            </div>
          </div>

          <div className="rounded-2xl border border-white/10 bg-white/[.04] p-4 md:p-5">
            <h3 className="mb-3 text-sm font-medium text-white/70">
              {global ? "Global (saved)" : "Global (saved) — no data yet"}
            </h3>
            <div className="grid grid-cols-2 gap-3">
              <StatCard label="Line Items" value={global?.lineItemCount ?? 0} />
              <StatCard label="Total Quantity" value={global?.totalQuantity ?? 0} />
              <StatCard label="Total Cost (LKR)" value={(global?.totalCost ?? 0).toFixed(2)} />
              <StatCard label="Avg Unit Price (LKR)" value={(global?.avgUnitPrice ?? 0).toFixed(2)} />
              <StatCard label="Avg Line Cost (LKR)" value={(global?.avgLineCost ?? 0).toFixed(2)} />
            </div>

            {batch && (
              <p className="mt-3 text-xs text-white/60">
                Last save batch: {batch.lineItemCount} line(s), LKR {batch.totalCost.toFixed(2)}.
              </p>
            )}
          </div>
        </section>

        {/* API footer */}
        <p className="mt-8 text-xs text-white/50">API: {API_URL}</p>
      </div>

      {/* Toast */}
      {toast && (
        <div
          className={`fixed right-4 top-4 z-50 rounded-xl px-4 py-3 text-sm backdrop-blur
            ${toast.kind === "ok"
              ? "bg-emerald-400/20 text-emerald-200 border border-emerald-300/30"
              : "bg-red-400/20 text-red-200 border border-red-300/30"}`}
          role="status" aria-live="polite"
        >
          {toast.text}
        </div>
      )}
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-white/10 bg-white/[.03] p-4">
      <div className="text-[11px] uppercase tracking-wide text-white/50">{label}</div>
      <div className="mt-1 text-xl font-semibold">{value}</div>
    </div>
  );
}
