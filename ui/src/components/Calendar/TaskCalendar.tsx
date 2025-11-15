import { useState, useMemo } from "react";
import { ChevronLeft, ChevronRight, Calendar, CalendarRange } from "lucide-react";
import dayjs from "dayjs";
import weekday from "dayjs/plugin/weekday";
import weekOfYear from "dayjs/plugin/weekOfYear";
import { m, AnimatePresence } from "framer-motion";
import { Ripple } from "@maxhub/max-ui";
import clsx from "clsx";
import { useFarmStore } from "../../stores/useFarmStore";
import { Task } from "../TodoList/ui/Task";

dayjs.extend(weekday);
dayjs.extend(weekOfYear);
dayjs.locale('ru');

export const TaskCalendar = () => {
  const { tasks, isDrought } = useFarmStore();
  const [currentWeekStart, setCurrentWeekStart] = useState(dayjs().startOf('week'));
  const [selectedDate, setSelectedDate] = useState(dayjs().startOf('day'));

  const weekDays = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => {
      const date = currentWeekStart.add(i, 'day');
      return {
        date,
        dayName: date.format('dd'),
        dayNumber: date.format('D'),
        isToday: date.isSame(dayjs(), 'day'),
        isSelected: date.isSame(selectedDate, 'day'),
        hasTasks: tasks.some(t => t.date && dayjs(t.date).isSame(date, 'day')),
      };
    });
  }, [currentWeekStart, selectedDate, tasks]);

  const tasksForSelectedDay = useMemo(() => {
    return tasks
      .filter(task => task.date && dayjs(task.date).isSame(selectedDate, 'day'))
      .sort((a, b) => (a.done === b.done ? 0 : a.done ? 1 : -1));
  }, [tasks, selectedDate]);

  const goToPrevWeek = () => setCurrentWeekStart(p => p.subtract(7, 'day'));
  const goToNextWeek = () => setCurrentWeekStart(p => p.add(7, 'day'));
  const goToToday = () => {
    const today = dayjs().startOf('day');
    setCurrentWeekStart(today.startOf('week'));
    setSelectedDate(today);
  };

  return (
    <div className="flex flex-col h-screen bg-base-100 touch-pan-y">
      <header className="sticky top-0 bg-base-100 z-20 border-b border-base-300 px-4 pt-4 pb-2">
        <div className="flex items-center justify-between">
          <h1 className={clsx(
            "text-lg font-bold flex items-center gap-1.5",
            { 
              "text-emerald-600": !isDrought,
              "text-amber-600": isDrought,
            },
          )}>
            <Calendar className="w-5 h-5" />
            Календарь
          </h1>
          <button
            onClick={goToToday}
            className={clsx(
              "text-xs px-3 py-1.5 rounded-full transition-all font-medium",
              selectedDate.isSame(dayjs(), 'day')
                ? isDrought
                  ? "bg-amber-600 text-white"
                  : "bg-emerald-600 text-white"
                : isDrought
                  ? "bg-base-200 text-amber-600"
                  : "bg-base-200 text-emerald-600"
            )}
          >
            Сегодня
          </button>
        </div>
      </header>

      <div className="px-4 py-3">
        <div className="flex items-center justify-between mb-2">
          <button
            onClick={goToPrevWeek}
            className="w-9 h-9 rounded-full flex items-center justify-center bg-base-200 hover:bg-emerald-50 dark:hover:bg-emerald-800 transition-colors"
            aria-label="Предыдущая неделя"
          >
            <ChevronLeft className={clsx("w-5 h-5", isDrought ? "text-amber-700" : "text-emerald-700")} />
          </button>

          <div className="text-center">
            <p className={clsx("text-xs font-medium", isDrought ? "text-amber-700" : "text-emerald-700")}>
              {currentWeekStart.format('D MMM')} — {currentWeekStart.add(6, 'day').format('D MMM')}
            </p>
            <p className={clsx("text-xs", isDrought ? "text-amber-600" : "text-emerald-600")}>{currentWeekStart.format('YYYY')}</p>
          </div>

          <button
            onClick={goToNextWeek}
            className="w-9 h-9 rounded-full flex items-center justify-center bg-base-200 hover:bg-emerald-50 dark:hover:bg-emerald-800 transition-colors"
            aria-label="Следующая неделя"
          >
            <ChevronRight className={clsx("w-5 h-5", isDrought ? "text-amber-700" : "text-emerald-700")} />
          </button>
        </div>

        <div className="grid grid-cols-7 gap-1">
          {weekDays.map((day) => (
            <button
              key={day.date.toString()}
              onClick={() => setSelectedDate(day.date)}
              className={clsx(
                "relative flex flex-col items-center justify-center py-2.5 rounded-xl transition-all duration-200 min-h-16 tap-highlight-transparent",
                "active:scale-95",
                {
                  "bg-emerald-600 text-white shadow-md": day.isSelected && !isDrought,
                  "bg-amber-600 text-white shadow-md": day.isSelected && isDrought,
                  "ring-2 ring-emerald-500 ring-offset-2 ring-offset-base-100": day.isToday && !day.isSelected && !isDrought,
                  "ring-2 ring-amber-500 ring-offset-2 ring-offset-base-100": day.isToday && !day.isSelected && isDrought,
                  "bg-base-200 hover:bg-emerald-50 dark:hover:bg-emerald-900": !day.isSelected && !day.isToday && !isDrought,
                  "bg-base-200 hover:bg-amber-50 dark:hover:bg-amber-900": !day.isSelected && !day.isToday && isDrought,
                }
              )}
            >
              <span className="text-xs opacity-70 font-medium">{day.dayName}</span>
              <span className={clsx("text-base font-bold", {
                "text-white": day.isSelected,
                [isDrought ? "text-amber-800" : "text-emerald-800"]: !day.isSelected,
              })}>
                {day.dayNumber}
              </span>
              {day.hasTasks && !day.isSelected && (
                <div className={clsx(
                  "absolute bottom-1 w-1.5 h-1.5 rounded-full shadow-sm",
                  isDrought ? "bg-amber-500" : "bg-emerald-500"
                )} />
              )}
            </button>
          ))}
        </div>
      </div>

      <div className="px-4 pb-2">
        <div className={clsx(
          "bg-base-200 rounded-xl p-3 text-center border",
          isDrought ? "border-amber-200" : "border-emerald-200",
          isDrought ? "dark:border-amber-600" : "dark:border-emerald-700"
        )}>
          <h3 className={clsx(
            "text-sm font-bold",
            isDrought ? "text-amber-800 dark:text-amber-500" : "text-emerald-800 dark:text-emerald-400"
          )}>
            {selectedDate.format('D MMMM YYYY')}
          </h3>
          <p className={clsx("text-xs mt-0.5", isDrought ? "text-amber-600" : "text-emerald-600")}>
            {tasksForSelectedDay.length} {tasksForSelectedDay.length === 1 ? 'задача' : 'задачи'}
          </p>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-4 pb-24 scrollbar-hide">
        <AnimatePresence mode="wait">
          {tasksForSelectedDay.length === 0 ? (
            <m.div
              key="empty"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="text-center py-16"
            >
              <div className={clsx(
                "w-16 h-16 mx-auto mb-3 rounded-full flex items-center justify-center",
                isDrought ? "bg-amber-100 dark:bg-amber-900/30" : "bg-emerald-100 dark:bg-emerald-900/30"
              )}>
                <CalendarRange className={clsx("w-8 h-8", isDrought ? "text-amber-500" : "text-emerald-500")} />
              </div>
              <p className={clsx("text-sm font-medium", isDrought ? "text-amber-600" : "text-emerald-600")}>Нет задач</p>
            </m.div>
          ) : (
            <m.ul
              key="tasks"
              className="space-y-2.5"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              <AnimatePresence>
                {tasksForSelectedDay.map((task, index) => (
                  <m.li
                    key={task.id}
                    initial={{ opacity: 0, x: -20, scale: 0.98 }}
                    animate={{ opacity: 1, x: 0, scale: 1 }}
                    exit={{ opacity: 0, x: 50, scale: 0.95 }}
                    transition={{ duration: 0.25, delay: index * 0.04 }}
                    className="relative rounded-xl overflow-hidden shadow-sm"
                  >
                    <Ripple className="inset-0 pointer-events-none" />
                    <Task task={task} />
                  </m.li>
                ))}
              </AnimatePresence>
            </m.ul>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
};
