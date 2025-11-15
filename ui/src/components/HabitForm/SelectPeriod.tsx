import type { UseFormRegister } from "react-hook-form";
import clsx from "clsx";
import type { DifficultyTask } from "../../types/farm";
import type { NewHabitFormValues } from "../Modal/Habit/AddHabitModal";

type Props = {
  register: UseFormRegister<NewHabitFormValues>;
  difficulty?: DifficultyTask;
};

export const SelectPeriod = ({ register, difficulty }: Props) => {
const getFocusBorderColor = () => {
    if (!difficulty) return "focus-within:border-emerald-500";
    switch (difficulty) {
      case "trifle": return "focus-within:border-slate-500";
      case "easy": return "focus-within:border-emerald-500";
      case "normal": return "focus-within:border-sky-500";
      case "hard": return "focus-within:border-orange-500";
      default: return "focus-within:border-emerald-500";
    }
  };

  const getLabelColor = () => {
    if (!difficulty) return "text-emerald-700";
    switch (difficulty) {
      case "trifle": return "text-slate-500";
      case "easy": return "text-emerald-600";
      case "normal": return "text-sky-600";
      case "hard": return "text-orange-600";
      default: return "text-emerald-700";
    }
  };

  return (
    <div className="relative w-full group">
      <select
        {...register("period")}
        id="period"
        className={clsx(
          "peer select w-full h-11 pl-3 pr-8 pt-5 pb-2 text-sm rounded-lg shadow-sm",
          "bg-base-100",
          "outline-none ring-0 transition-colors duration-300",
          getFocusBorderColor(),
        )}
        defaultValue="day"
      >
        <option value="day">Ежедневный</option>
        <option value="week">Еженедельный</option>
        <option value="month">Ежемесячный</option>
        <option value="year">Ежегодный</option>
      </select>

      <label
        htmlFor="period"
        className={clsx(
          "absolute left-3 top-1 text-xs font-mono pointer-events-none select-none z-10",
          getLabelColor()
        )}
      >
        Период
      </label>
    </div>
  );
};