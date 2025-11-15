import type { UseFormRegister } from "react-hook-form";
import clsx from "clsx";
import type { DifficultyTask } from "../../types/farm";
import type { NewHabitFormValues } from "../Modal/Habit/AddHabitModal";

type Props = {
  register: UseFormRegister<NewHabitFormValues>;
  autoFocus?: boolean;
  difficulty?: DifficultyTask;
};

export const InputTitle = ({ register, autoFocus=false, difficulty }: Props) => {
  return (
    <div className="relative w-full group">
      <input
        {...register("title")}
        autoFocus={autoFocus}
        type="text"
        id="input"
        placeholder=" "
        className={clsx(
          "peer w-full pl-3 pr-4 pt-4 pb-2 text-sm",
          "rounded-lg shadow-md focus:outline-none transition-all duration-300",
          {
            "bg-emerald-100 text-emerald-900": !difficulty,
            "bg-slate-50 text-slate-800": difficulty === "trifle",
            "bg-emerald-50 text-emerald-800": difficulty === "easy",
            "bg-sky-50 text-sky-800": difficulty === "normal",
            "bg-orange-50 text-orange-800": difficulty === "hard",
          }
        )}
      />
      <label
        htmlFor="input"
        className={clsx(
          "absolute font-mono left-3 top-1 text-xs transition-all duration-200 ease-in-out",
                    "peer-placeholder-shown:top-2.5",
                    "peer-placeholder-shown:text-base",
                    "peer-focus:top-1",
                    "peer-focus:text-xs" ,
                    "cursor-text select-none pointer-events-none",
          {
            "peer-placeholder-shown:text-emerald-600 peer-focus:text-emerald-600 text-emerald-800": !difficulty,
            "peer-placeholder-shown:text-slate-400 peer-focus:text-slate-500 text-slate-500": difficulty === "trifle",
            "peer-placeholder-shown:text-emerald-400 peer-focus:text-emerald-500 text-emerald-500": difficulty === "easy",
            "peer-placeholder-shown:text-sky-400 peer-focus:text-sky-500 text-sky-500": difficulty === "normal",
            "peer-placeholder-shown:text-orange-400 peer-focus:text-orange-500 text-orange-500": difficulty === "hard",
          }
        )}
      >
        Название задачи
      </label>
    </div>
  );
};