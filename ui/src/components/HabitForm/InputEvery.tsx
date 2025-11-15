import type { UseFormRegister } from "react-hook-form";
import clsx from "clsx";
import type { DifficultyTask } from "../../types/farm";
import type { NewHabitFormValues } from "../Modal/Habit/AddHabitModal";

type Props = {
  register: UseFormRegister<NewHabitFormValues>;
  difficulty?: DifficultyTask;
};

export const InputEvery = ({ register, difficulty }: Props) => {
  const getFocusBorderColor = () => {
    if (!difficulty) return "focus:border-emerald-500";
    switch (difficulty) {
      case "trifle":
        return "focus:border-slate-500";
      case "easy":
        return "focus:border-emerald-500";
      case "normal":
        return "focus:border-sky-500";
      case "hard":
        return "focus:border-orange-500";
      default:
        return "focus:border-emerald-500";
    }
  };

  const getLabelColor = () => {
    if (!difficulty) return "text-emerald-700";
    switch (difficulty) {
      case "trifle":
        return "text-slate-500";
      case "easy":
        return "text-emerald-600";
      case "normal":
        return "text-sky-600";
      case "hard":
        return "text-orange-600";
      default:
        return "text-emerald-700";
    }
  };

  return (
    <div className="relative w-36 group">
      <input
        {...register("every", { valueAsNumber: true })}
        type="number"
        id="input"
        placeholder=" "
        className={clsx(
          "peer input w-full h-11 pl-3 pr-8 pt-5 pb-2 text-sm rounded-lg shadow-sm",
          "bg-base-100",
          "focus:outline-none focus:ring-0 transition-all duration-300",
          getFocusBorderColor(),
          "placeholder-transparent",
          "[appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
        )}
        max={100}
        min={1}
      />
      <label
        htmlFor="input"
        className={clsx(
          "absolute left-3 top-1 text-xs font-mono pointer-events-none select-none z-10",
          "transition-all duration-200 ease-in-out",
          "peer-focus:text-xs peer-focus:top-1",
          getLabelColor()
        )}
      >
        Каждый
      </label>
    </div>
  );
};