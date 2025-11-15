import { Sprout } from "lucide-react";
import clsx from "clsx";
import { useMemo, useCallback, memo } from "react";
import type { DifficultyTask } from "../../types/farm";

const difficulties: { key: DifficultyTask; label: string; sprouts: number }[] = [
  { key: "trifle", label: "Просто", sprouts: 1 },
  { key: "easy", label: "Легко", sprouts: 2 },
  { key: "normal", label: "Нормально", sprouts: 3 },
  { key: "hard", label: "Сложно", sprouts: 4 },
];

const difficultyColors = {
  trifle: { bg: "bg-slate-50", text: "text-slate-800", sprout: "text-slate-800", label: "text-slate-500", labelActive: "font-bold text-slate-500", ring: "ring-slate-500" },
  easy: { bg: "bg-emerald-50", text: "text-emerald-800", sprout: "text-emerald-800", label: "text-emerald-500", labelActive: "font-bold text-emerald-500", ring: "ring-emerald-500" },
  normal: { bg: "bg-sky-50", text: "text-sky-800", sprout: "text-sky-800", label: "text-sky-500", labelActive: "font-bold text-sky-500", ring: "ring-sky-500" },
  hard: { bg: "bg-orange-50", text: "text-orange-800", sprout: "text-orange-800", label: "text-orange-500", labelActive: "font-bold text-orange-500", ring: "ring-orange-500" },
  default: { bg: "bg-emerald-50", text: "text-emerald-800", sprout: "text-emerald-800", label: "text-emerald-800", labelActive: "font-bold text-emerald-500", ring: "ring-emerald-500" },
};

type Props = {
  value: DifficultyTask;
  onChange: (value: DifficultyTask) => void;
};


const DifficultyButton = memo(
  ({ keyName, label, sprouts, isActive, colors, onSelect }: {
    keyName: DifficultyTask;
    label: string;
    sprouts: number;
    isActive: boolean;
    colors: typeof difficultyColors.trifle;
    onSelect: (key: DifficultyTask) => void;
  }) => {
    const handleClick = useCallback(() => onSelect(keyName), [onSelect, keyName]);
    
    const sproutIcons = useMemo(() => 
      Array.from({ length: sprouts }).map((_, i) => <Sprout key={i} className={colors.sprout} />),
      [sprouts, colors.sprout]
    );

    return (
      <div className="flex flex-col items-center gap-1">
        <button
          type="button"
          onClick={handleClick}
          className={clsx(
            "flex flex-wrap items-center justify-center w-16 h-16 rounded-lg transition-all duration-200",
            colors.bg,
            colors.text,
            { ["ring-2 " + colors.ring + " scale-105"]: isActive }
          )}
        >
          {sproutIcons}
        </button>
        <span className={clsx("font-mono text-sm transition-colors", isActive ? colors.labelActive : colors.label)}>
          {label}
        </span>
      </div>
    );
  }
);

export const DifficultySelector = memo(({ value, onChange }: Props) => {
  return (
    <div className="flex justify-around gap-2 w-full">
      {difficulties.map(({ key, label, sprouts }) => {
        const isActive = value === key;
        const colors = difficultyColors[key] || difficultyColors.default;

        return (
          <DifficultyButton
            key={key}
            keyName={key}
            label={label}
            sprouts={sprouts}
            isActive={isActive}
            colors={colors}
            onSelect={onChange}
          />
        );
      })}
    </div>
  );
});
