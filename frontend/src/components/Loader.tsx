import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";

/**
 * Loader – dark neon hero with typewriter “Costly”, glass card and glowing CTA.
 */
export default function Loader() {
  const nav = useNavigate();

  // Allow Enter to start
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Enter") nav("/items");
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [nav]);

  return (
    <div className="relative min-h-dvh overflow-hidden bg-ink900 text-text-main font-sans">
      {/* Layer 0 — subtle vignette */}
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(70%_60%_at_50%_0%,rgba(122,92,255,.10),transparent),radial-gradient(60%_50%_at_50%_100%,rgba(0,229,255,.08),transparent)]" />

      {/* Layer 1 — neon mist blobs */}
      <div className="pointer-events-none absolute -top-28 -left-24 h-80 w-80 rounded-full bg-neon-cyan/25 blur-3xl" />
      <div className="pointer-events-none absolute -bottom-28 -right-24 h-96 w-96 rounded-full bg-neon-violet/25 blur-3xl" />
      <div className="pointer-events-none absolute top-1/3 right-24 h-40 w-40 rounded-full bg-neon-magenta/25 blur-2xl animate-float" />

      {/* Layer 2 — aurora sweep */}
      <div className="pointer-events-none absolute inset-x-0 top-1/3 h-24 -skew-y-6 opacity-20">
        <div className="h-full w-[40%] bg-gradient-to-r from-transparent via-neon-cyan/50 to-transparent blur-xl animate-sweep" />
      </div>

      {/* Hero */}
      <main className="relative z-10 mx-auto flex max-w-3xl flex-col items-center px-6 pt-28 md:pt-36">
        <div className="w-full rounded-3xl border border-white/10 bg-white/[.055] p-8 text-center backdrop-blur-xl shadow-[0_0_40px_rgba(122,92,255,.18),0_0_60px_rgba(0,229,255,.12)]">
          {/* Wordmark */}
          <div className="mb-4 inline-flex items-end justify-center">
            <span
              className="font-cursive text-6xl sm:text-7xl md:text-8xl tracking-wide inline-block overflow-hidden whitespace-nowrap pr-1
                         [text-shadow:0_0_18px_rgba(0,229,255,.55),0_0_28px_rgba(122,92,255,.45)] border-r-2 border-neon-magenta/80 animate-typewriter"
              style={{ width: "0ch" }}
              aria-label="Costly"
            >
              Costly
            </span>
            <span className="ml-0.5 mb-2 h-9 w-[2px] rounded-sm bg-neon-magenta/80 animate-caret" aria-hidden />
          </div>

          {/* Tagline */}
          <p className="mx-auto max-w-xl text-text-dim">
            A minimal cost calculator with a <span className="text-text-main/90">luxurious</span> glow.
            Enter items, preview live totals, then save to see global stats.
          </p>

          {/* CTA */}
          <div className="mt-8 flex justify-center">
            <button
              onClick={() => nav("/items")}
              className="group relative inline-flex items-center gap-2 rounded-full px-7 py-3 font-medium
                         text-ink950 bg-gradient-to-r from-neon-cyan to-neon-violet
                         shadow-[0_0_24px_rgba(122,92,255,.35),0_0_48px_rgba(0,229,255,.25)] animate-pulseGlow
                         hover:brightness-110 focus:outline-none focus:ring-2 focus:ring-neon-cyan/40 transition"
            >
              <Sparkle />
              <span>Start</span>
            </button>
          </div>

          {/* Hint */}
          <div className="mt-6 text-xs text-text-dim">Press <kbd className="rounded bg-white/10 px-1">Enter</kbd> to start</div>
        </div>
      </main>
    </div>
  );
}

/** Simple sparkle icon for the CTA */
function Sparkle() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5 text-ink950" fill="currentColor" aria-hidden>
      <path d="M12 2l2.1 4.2L18.5 8 14.1 9.9 12 14l-2.1-4.1L5.5 8l4.4-1.8L12 2zM6.3 16.6l.7 1.4 1.4.7-1.4.7-.7 1.4-.7-1.4-1.4-.7 1.4-.7.7-1.4zm9.6-.6l.8 1.6 1.6.8-1.6.8-.8 1.6-.8-1.6-1.6-.8 1.6-.8.8-1.6z" />
    </svg>
  );
}
