import { Flame } from "lucide-react";
import clsx from "clsx";
import { useFarmStore } from "../../../stores/useFarmStore";

export const StrickCounter = () => {
  const { strick, didTaskToday } = useFarmStore();



  return (
    <div className="flex items-center gap-1.5">
      <p className={clsx(
          "font-mono",
          {
            "text-orange-800 dark:text-orange-400": didTaskToday,
            "text-gray-500 dark:text-gray-400": !didTaskToday,
          }
        )}
      >
        {strick}
      </p>

      <Flame
        className={clsx(
          "w-4 h-4",
          {
            "text-orange-600 dark:text-orange-400": didTaskToday,
            "text-gray-400 dark:text-gray-500": !didTaskToday,
          }
        )}
        fill={didTaskToday ? "currentColor" : "none"}
        strokeWidth={didTaskToday ? 2.5 : 2}
      />

    </div>
  );
};
