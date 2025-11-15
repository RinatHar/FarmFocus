import { ru } from "date-fns/locale";
import { DayPicker } from "react-day-picker";
import clsx from "clsx";
import { useMemo, useCallback, memo } from "react";
import type { DifficultyTask } from "../../types/farm";

type Props = {
  value?: Date | null;
  difficulty?: DifficultyTask;
  onChange: (date: Date | undefined) => void;
};

export const DatePicker = memo(({ value, difficulty, onChange }: Props) => {
  const buttonClass = useMemo(() => {
    const base = "input outline-none w-full rounded-lg";
    const focusClass = {
      undefined: "focus:border-emerald-400 dark:focus:border-emerald-500",
      trifle: "focus:border-slate-200 dark:focus:border-slate-300",
      easy: "focus:border-emerald-200 dark:focus:border-emerald-300",
      normal: "focus:border-sky-200 dark:focus:border-sky-300",
      hard: "focus:border-orange-200 dark:focus:orange-emerald-300",
    };
    return clsx(base, focusClass[difficulty ?? "undefined"]);
  }, [difficulty]);

  const handleSelect = useCallback(
    (date?: Date) => {
      onChange(date);
    },
    [onChange]
  );

  return (
    <>
      <button
        type="button"
        popoverTarget="rdp-popover"
        className={buttonClass}
        style={{ anchorName: "--rdp" } as React.CSSProperties}
      >
        {value ? value.toLocaleDateString() : "- Выберите дату -"}
      </button>

      <div popover="auto" id="rdp-popover" className="dropdown dropdown-top" style={{ positionAnchor: "--rdp" } as React.CSSProperties}>
        <DayPicker
          className="react-day-picker"
          mode="single"
          selected={value || undefined}
          onSelect={handleSelect}
          locale={ru}
        />
      </div>
    </>
  );
});


export default DatePicker;