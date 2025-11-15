import { Ripple } from "@maxhub/max-ui";
import { m, AnimatePresence } from "framer-motion";
import { useFarmStore } from "../../stores/useFarmStore";
import { Habit } from "./ui/Habit";
import { BriefcaseBusiness } from "lucide-react";
import clsx from "clsx";

export const HabitList = () => {
  const { habits, isDrought } = useFarmStore();


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
            <BriefcaseBusiness className="w-5 h-5" />
            Дела
          </h1>
        </div>
      </div>

      <div className="p-4 pb-20 overflow-y-auto scrollbar-hide h-screen">
        <ul className="flex flex-col gap-2.5">
          <AnimatePresence>
            {
              habits.map((habit) => {

                return (
                  <m.li
                    key={habit.id}
                    className="relative w-full shadow-lg rounded-md overflow-hidden"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: 100 }}
                    transition={{ duration: 0.4 }}
                    layout
                  >
                    <Ripple className="absolute pointer-events-none" />
                    <Habit habit={habit} />
                  </m.li>
                );
              })
            }
          </AnimatePresence>
        </ul>
      </div>
    </>
  );
};