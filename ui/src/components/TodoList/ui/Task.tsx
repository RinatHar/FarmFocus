import { Ripple } from "@maxhub/max-ui";
import clsx from "clsx";
import { Calendar, Check } from "lucide-react";
import { useFarmStore } from "../../../stores/useFarmStore";
import { useEditTaskModalStore } from "../../../stores/editTaskModal";
import { difficultyColor } from "../../../utils/difficultyColor";
import dayjs from "dayjs";
import 'dayjs/locale/ru';
import type { ITask } from "../../../types/farm";

dayjs.locale('ru');

export const Task = ({ task }: { task: ITask }) => {
  const { toggleTask } = useFarmStore();
  const { open } = useEditTaskModalStore();

  const handleToggleTask = (id: number) => {
    setTimeout(() => toggleTask(id), 100);
  };

  return (
    <div
      className={clsx("flex w-full min-h-[60px]", {
        "opacity-50 grayscale": task.done,
      })}
    >
      <button
        onClick={() => handleToggleTask(task.id)}
        className={clsx(
          "relative w-11 p-1 overflow-hidden",
          difficultyColor(task.difficulty)
        )}
      >
        <Ripple className="absolute pointer-events-none" />

        <div className="absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2 w-6 h-6 rounded-full bg-gray-100/60" />

        {task.done && (
          <Check
            strokeWidth={4}
            className={clsx(
              "absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2 size-5",
              difficultyColor(task.difficulty, "text")
            )}
          />
        )}
      </button>

      <button
        onClick={() => open(task.id)}
        className="flex flex-col justify-center w-full bg-base-200 px-3 py-2 text-left"
      >
        <span className={clsx({ "line-through": task.done })}>
          {task.title}
        </span>

        {task.description && (
          <span
            className={clsx(
              "text-sm opacity-50",
              { "line-through": task.done }
            )}
          >
            {task.description}
          </span>
        )}

        <div className="flex justify-between">   
          {task.tag && (
            <div className="mt-1">
              <span
                className="badge"
                style={{ backgroundColor: task.tag.color, color: "#fff" }}
                >
                {task.tag.name}
              </span>
            </div>
          )}

          {task.date && (
            <div className="ml-auto flex items-center gap-1">
              <Calendar className="w-4 h-4 text-gray-400" />
              <span className="text-sm text-gray-400">
                {dayjs(task.date).format("DD MMM YYYY")}
              </span>
            </div>
          )}
        </div>

      </button>
    </div>
  );
};
