import clsx from "clsx";
import { Check } from "lucide-react";
import type { DifficultyTask } from "../../../types/farm";

interface DifficultyOption {
  value: DifficultyTask;
  label: string;
  color: string;
  bg: string;
}

const options: DifficultyOption[] = [
  { value: "trifle", label: "Просто", color: "bg-gray-500", bg: "bg-gray-100" },
  { value: "easy", label: "Легко", color: "bg-green-500", bg: "bg-green-100" },
  { value: "normal", label: "Нормально", color: "bg-blue-500", bg: "bg-blue-100" },
  { value: "hard", label: "Сложно", color: "bg-orange-500", bg: "bg-orange-100" },
];

interface Props {
  selected: DifficultyTask[];
  onChange: (values: DifficultyTask[]) => void;
}

export const DifficultyFilter = ({ selected, onChange }: Props) => {
  const toggle = (value: DifficultyTask) => {
    const newValues = selected.includes(value)
      ? selected.filter((v) => v !== value)
      : [...selected, value];
    onChange(newValues);
  };

  return (
    <div className="grid grid-cols-2 gap-3">
      {options.map((opt) => {
        const isSelected = selected.includes(opt.value);
        return (
          <button
            key={opt.value}
            onClick={() => toggle(opt.value)}
            className={clsx(
              "flex items-center gap-2 px-4 py-2.5 rounded-full font-medium text-sm transition-all",
              {
                "text-base-content ring-1 ring-current": isSelected,
                "bg-base-200 text-base-content/70 hover:bg-base-300": !isSelected,
              }
            )}
          >
            <div className={clsx("w-3 h-3 rounded-full", opt.color)} />
            {opt.label}
            {isSelected && <Check className="w-4 h-4" />}
          </button>
        );
      })}
    </div>
  );
};