import { Ripple } from "@maxhub/max-ui";
import clsx from "clsx";
import { Check, ClockCheck } from "lucide-react";
import { useFarmStore } from "../../../stores/useFarmStore";
import { difficultyColor } from "../../../utils/difficultyColor";
import type { IHabit } from "../../../types/farm";
import { useEditHabitModalStore } from "../../../stores/editHabitModal";



export const Habit = ({ habit }: { habit: IHabit }) => {
  const { toggleHabit } = useFarmStore();
  const { open } = useEditHabitModalStore();

  const handleToggle = (id: number) => {
    toggleHabit(id);
  };

  return (
    <div
      className={clsx("flex w-full min-h-[60px]", {
        "opacity-50 grayscale": habit.done,
      })}
    >
      <button
        onClick={() => handleToggle(habit.id)}
        className={clsx(
          "relative w-11 p-1 overflow-hidden",
          difficultyColor(habit.difficulty)
        )}
      >
        <Ripple className="absolute pointer-events-none" />

        <div className="absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2 w-6 h-6 rounded-lg bg-gray-100/60" />

        {habit.done && (
          <Check
            strokeWidth={4}
            className={clsx(
              "absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2 size-5",
              difficultyColor(habit.difficulty, "text")
            )}
          />
        )}
      </button>

      <button
        onClick={() => open(habit.id)}
        className="flex flex-col justify-center w-full bg-base-200 px-3 py-2 text-left"
      >
        <span className={clsx({ "line-through": habit.done })}>
          {habit.title}
        </span>

        {habit.description && (
          <span
            className={clsx(
              "text-sm opacity-50",
              { "line-through": habit.done }
            )}
          >
            {habit.description}
          </span>
        )}

        <div className="flex justify-between">   
          {habit.tag && (
            <div className="mt-1">
              <span
                className="badge"
                style={{ backgroundColor: habit.tag.color, color: "#fff" }}
                >
                {habit.tag.name}
              </span>
            </div>
          )}


          <div className="ml-auto flex items-center gap-1">
            <ClockCheck className="w-4 h-4 text-gray-400" />
            <span className="text-sm text-gray-400">
              {habit.count}
            </span>
          </div>

        </div>

      </button>
    </div>
  );
};
