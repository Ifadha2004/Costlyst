import React, { useState } from "react";
import ItemsTable from "./ItemsTable";
import { saveAndCalculate } from "./api";
import type { Item, Stats } from "./types";
import { normalize, isValid } from "./validators";
import { API_URL } from "../lib/http";

export default function ItemsPage() {
  const [items, setItems] = useState<Item[]>([{ name: "", price: 0, quantity: 1 }]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [ok, setOk] = useState<string | null>(null);

  const addRow = () => setItems((p) => [...p, { name: "", price: 0, quantity: 1 }]);
  const delRow = (i: number) => setItems((p) => p.filter((_, idx) => idx !== i));
  const update = (i: number, k: keyof Item, v: Item[keyof Item]) =>
    setItems((p) => p.map((it, idx) => (idx === i ? { ...it, [k]: v } : it)));

  const submit = async () => {
    setError(null);
    setOk(null);
    setLoading(true);
    try {
      const payload = items.map(normalize).filter(isValid);
      if (payload.length === 0) throw new Error("Add at least one valid item (name, price ≥ 0, quantity ≥ 1).");
      const data = await saveAndCalculate(payload);
      setStats(data);
      setOk("Items saved and stats calculated.");
    } catch (e: any) {
      setError(e.message || "Something went wrong");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 960, margin: "40px auto", padding: 24, fontFamily: "system-ui, Arial" }}>
      <h1>Items Cost Calculator</h1>
      <p>Enter name, unit price, and quantity. Backend saves and returns batch totals & averages.</p>

      <ItemsTable items={items} onChange={update} onAdd={addRow} onDelete={delRow} />

      <div style={{ display: "flex", gap: 12, marginTop: 12 }}>
        <button onClick={submit} disabled={loading}>{loading ? "Saving..." : "Save & Calculate"}</button>
      </div>

      {error && <div style={{ color: "crimson", marginTop: 8 }}>{error}</div>}
      {ok && <div style={{ color: "green", marginTop: 8 }}>{ok}</div>}

      {stats && (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(2, 1fr)", gap: 12, marginTop: 16 }}>
          <Card label="Line Items" value={stats.lineItemCount} />
          <Card label="Total Quantity" value={stats.totalQuantity} />
          <Card label="Total Cost (LKR)" value={stats.totalCost.toFixed(2)} />
          <Card label="Avg Unit Price (LKR)" value={stats.avgUnitPrice.toFixed(2)} />
          <Card label="Avg Line Cost (LKR)" value={stats.avgLineCost.toFixed(2)} />
        </div>
      )}

      <p style={{ color: "#777", marginTop: 16 }}>API: {API_URL}</p>
    </div>
  );
}

function Card({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div style={{ border: "1px solid #ddd", borderRadius: 8, padding: 12 }}>
      <div style={{ fontSize: 12, color: "#666" }}>{label}</div>
      <div style={{ fontSize: 22, fontWeight: 600 }}>{value}</div>
    </div>
  );
}
