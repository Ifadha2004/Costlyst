import React, { useEffect } from "react";
import { motion } from "framer-motion";
import { useNavigate } from "react-router-dom";

const WRITE_MS = 2800;
const HOLD_MS = 600;

export default function Loader() {
  const nav = useNavigate();

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === "Enter" && nav("/items");
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [nav]);

  // Natural handwriting pacing
  const times = [0, 0.15, 0.32, 0.5, 0.67, 0.83, 0.93, 1];
  const penStops = ["0%", "14%", "30%", "47%", "64%", "80%", "90%", "100%"];

  // IMPORTANT: give extra vertical space in the clip-path mask so glow & tails aren’t cropped
  const TB = "-12%"; // top/bottom expansion
  const clipFrames = [
    `inset(${TB} 100% ${TB} 0%)`,
    `inset(${TB} 86%  ${TB} 0%)`,
    `inset(${TB} 70%  ${TB} 0%)`,
    `inset(${TB} 53%  ${TB} 0%)`,
    `inset(${TB} 36%  ${TB} 0%)`,
    `inset(${TB} 20%  ${TB} 0%)`,
    `inset(${TB} 10%  ${TB} 0%)`,
    `inset(${TB} 0%   ${TB} 0%)`,
  ];

  return (
    <div className="relative min-h-dvh flex items-center justify-center bg-[#06060e] overflow-hidden font-sans text-white/90">
      {/* background glow */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(122,92,255,0.18),transparent),radial-gradient(circle_at_70%_80%,rgba(0,229,255,0.18),transparent)]" />
      <div className="pointer-events-none absolute -top-32 -left-32 h-96 w-96 rounded-full bg-[#00E5FF]/15 blur-3xl" />
      <div className="pointer-events-none absolute bottom-0 right-0 h-[28rem] w-[28rem] rounded-full bg-[#7A5CFF]/15 blur-3xl" />

      {/* card */}
      <div className="relative text-center px-6 py-10 rounded-3xl border border-white/10 bg-white/[.05] backdrop-blur-2xl shadow-[0_0_60px_rgba(122,92,255,.15)]">
        <div className="relative inline-block py-3"> {/* extra vertical padding for safety */}
          {/* invisible outline to stabilize layout */}
          <h1
            aria-hidden
            className="select-none font-[Great_Vibes]"
            style={{
              fontSize: "13vw",          // ↓ slightly smaller than before
              // desktop cap (matches the animated one below)
              // fallback for large screens via clamp:
              // @ts-ignore
              // eslint-disable-next-line
              fontSizeAdjust: undefined,
              lineHeight: 1.12,          // ↑ a bit taller to fit C/y
              color: "transparent",
              WebkitTextStroke: "0.7px rgba(255,255,255,0.08)",
              letterSpacing: "0.015em",
            }}
          >
            Costly
          </h1>

          {/* animated reveal */}
          <motion.div
            className="absolute inset-0"
            initial={{ clipPath: clipFrames[0] }}
            animate={{ clipPath: clipFrames }}
            transition={{ duration: WRITE_MS / 1000, ease: "easeInOut", times }}
            style={{ willChange: "clip-path" }}
          >
            <h1
              className="select-none font-[Great_Vibes]"
              style={{
                // match size/leading exactly
                fontSize: "13vw",
                lineHeight: 1.12,
                background:
                  "linear-gradient(90deg, #7A5CFF 0%, #00E5FF 50%, #7A5CFF 100%)",
                WebkitBackgroundClip: "text",
                backgroundClip: "text",
                color: "transparent",
                textShadow:
                  "0 0 15px rgba(122,92,255,.6), 0 0 30px rgba(0,229,255,.4)",
              }}
            >
              Costly
            </h1>
          </motion.div>

          {/* glow “pen” bead */}
          <motion.span
            className="absolute top-1/2 -translate-y-1/2 h-3 w-3 rounded-full"
            style={{
              backgroundColor: "#00E5FF",
              boxShadow: "0 0 10px #00E5FF, 0 0 22px rgba(0,229,255,0.6)",
            }}
            initial={{ left: penStops[0], opacity: 1 }}
            animate={{ left: penStops }}
            transition={{ duration: WRITE_MS / 1000, ease: "easeInOut", times }}
          />
        </div>

        {/* tagline */}
        <motion.p
          className="mt-4 text-sm text-gray-300 max-w-md mx-auto tracking-wide"
          initial={{ opacity: 0, y: 6 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 2.3, duration: 0.6 }}
        >
          An efficient system for recording costs, computing totals, and generating insightful summaries.
        </motion.p>

        {/* CTA */}
        <motion.button
          onClick={() => nav("/items")}
          className="mt-7 px-7 py-3 rounded-full text-black font-medium bg-gradient-to-r from-[#00E5FF] to-[#7A5CFF]
                     shadow-[0_0_30px_rgba(122,92,255,0.4),0_0_40px_rgba(0,229,255,0.3)]
                     hover:brightness-110 transition-all focus:outline-none"
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 2.8, duration: 0.5 }}
        >
          ✨ Start
        </motion.button>

        {/* hint */}
        <motion.p
          className="mt-3 text-xs text-gray-400"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 3.2, duration: 0.4 }}
        >
          Press <kbd className="px-1 bg-white/10 rounded">Enter</kbd> to start
        </motion.p>
      </div>
    </div>
  );
}
