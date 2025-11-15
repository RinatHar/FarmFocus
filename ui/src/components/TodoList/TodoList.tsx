import { Ripple } from "@maxhub/max-ui";
import { m, AnimatePresence } from "framer-motion";
import { Task } from "./ui/Task";
import { useMemo } from "react";
import { useTaskFiltersStore } from "../../stores/useTaskFiltersStore";
import { useFarmStore } from "../../stores/useFarmStore";
import { FilterButton } from "./ui/FilterButton";
import { LayoutList } from "lucide-react";
import clsx from "clsx";

export const TodoList = () => {
  const { tasks, isDrought } = useFarmStore();
  const { selectedDifficulties, completionFilter } = useTaskFiltersStore();

  const filteredTasks = useMemo(() => {
    return tasks.filter((task) => {
      if (selectedDifficulties.length > 0 && !selectedDifficulties.includes(task.difficulty)) {
        return false;
      }

      if (completionFilter === "done" && !task.done) return false;
      if (completionFilter === "undone" && task.done) return false;

      return true;
    });
  }, [tasks, selectedDifficulties, completionFilter]);

  return (
    <>
      <div className="sticky top-0 bg-base-100 z-10 p-4 pb-2 border-b border-base-300">
        <div className="flex items-center justify-between">
          <h1 className={clsx(
            "text-lg font-bold flex items-center gap-1.5",
            { 
              "text-emerald-600": !isDrought,
              "text-amber-600": isDrought,
            },
          )}>
            <LayoutList className="w-5 h-5" />
            Задачи
          </h1>
          <FilterButton />
        </div>
      </div>

      <div className="p-4 pb-20 overflow-y-auto scrollbar-hide h-screen">
        <ul className="flex flex-col gap-2.5">
          <AnimatePresence>
            {filteredTasks.length === 0 ? (
              <m.p
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="text-center text-base-content/50 py-8"
              >
                {selectedDifficulties.length > 0 || completionFilter !== "all"
                  ? "Задачи не найдены"
                  : "Нет задач"}
              </m.p>
            ) : (
              filteredTasks.map((task) => (
                <m.li
                  key={task.id}
                  className="relative w-full shadow-lg rounded-md overflow-hidden"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 100 }}
                  transition={{ duration: 0.4 }}
                  layout
                >
                  <Ripple className="absolute pointer-events-none" />
                  <Task task={task} />
                </m.li>
              ))
            )}
          </AnimatePresence>
        </ul>
      </div>
    </>
  );
};