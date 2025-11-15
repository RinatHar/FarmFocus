import { useEffect, useState } from "react";

export const ThemeToggle = () => {
  const [theme, setTheme] = useState(localStorage.getItem("theme") || "light");

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
  }, [theme]);

  return (
    <label className="swap swap-rotate">
      {/* Светлая иконка */}
      <input
        type="checkbox"
        onChange={() => setTheme(theme === "light" ? "dark" : "light")}
        checked={theme === "dark"}
      />
      {/* ☀️ и 🌙 из DaisyUI */}
      <svg
        className="swap-on fill-current w-8 h-8"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
      >
        <path d="M5.64 17.657l1.414-1.414A7.974 7.974 0 0016 18a8 8 0 00-8-8c-.546 0-1.08.056-1.6.162l1.24-1.24a10 10 0 0114.142 14.142l-1.414 1.414A10 10 0 015.64 17.657z" />
      </svg>
      <svg
        className="swap-off fill-current w-8 h-8"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
      >
        <path d="M21.64 13.354A9 9 0 1110.646 2.36 7 7 0 0021.64 13.354z" />
      </svg>
    </label>
  );
};
