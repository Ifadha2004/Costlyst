/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx,js,jsx}"],
  theme: {
    extend: {
      colors: {
        // Dark canvas
        ink900: "#0b0f14",
        ink950: "#070a0e",

        // Neon accents
        neon: {
          cyan: "#00e5ff",
          violet: "#7a5cff",
          magenta: "#ff47a1",
          lime: "#b0ff72",
        },

        // Copy
        text: {
          main: "#e6f0ff",
          dim: "#8ea3b0",
        },
      },
      fontFamily: {
        cursive: ['"Great Vibes"', "cursive"], // wordmark
        sans: ["Inter", "system-ui", "sans-serif"],
      },
      keyframes: {
        typewriter: {
          "0%": { width: "0ch" },
          "100%": { width: "6ch" }, // C-o-s-t-l-y
        },
        caret: {
          "0%, 49%": { opacity: "1" },
          "50%, 100%": { opacity: "0" },
        },
        float: {
          "0%": { transform: "translateY(0) scale(1)" },
          "50%": { transform: "translateY(-6px) scale(1.02)" },
          "100%": { transform: "translateY(0) scale(1)" },
        },
        sweep: {
          "0%": { transform: "translateX(-30%)" },
          "100%": { transform: "translateX(130%)" },
        },
        pulseGlow: {
          "0%, 100%": { boxShadow: "0 0 24px rgba(122,92,255,.35), 0 0 48px rgba(0,229,255,.25)" },
          "50%": { boxShadow: "0 0 36px rgba(122,92,255,.55), 0 0 64px rgba(0,229,255,.4)" },
        },
      },
      animation: {
        typewriter: "typewriter 1.6s steps(6) .15s forwards",
        caret: "caret 750ms steps(1) infinite",
        float: "float 6s ease-in-out infinite",
        sweep: "sweep 6s linear infinite",
        pulseGlow: "pulseGlow 2.6s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};
