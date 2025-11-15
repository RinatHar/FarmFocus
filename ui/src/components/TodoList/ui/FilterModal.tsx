import { AnimatePresence, m } from "framer-motion";
import { X } from "lucide-react";
import { useMemo, useState } from "react";
import { useFilterModalStore } from "../../../stores/useFilterModalStore";
import { useTaskFiltersStore, type CompletionFilter } from "../../../stores/useTaskFiltersStore";
import { useFarmStore } from "../../../stores/useFarmStore";
import { DifficultyFilter } from "./DifficultyFilter";
import type { DifficultyTask } from "../../../types/farm";

export const FilterModal = () => {
  const { isOpen, close } = useFilterModalStore();
  const { selectedDifficulties, completionFilter, setDifficulties, setCompletionFilter, clear } = useTaskFiltersStore();
  const { tasks } = useFarmStore();

  const [localDifficulties, setLocalDifficulties] = useState<DifficultyTask[]>(selectedDifficulties);
  const [localCompletion, setLocalCompletion] = useState(completionFilter);

  const filteredCount = useMemo(() => {
    return tasks.filter((t) => {
      if (localDifficulties.length > 0 && !localDifficulties.includes(t.difficulty)) return false;
      if (localCompletion === "done" && !t.done) return false;
      if (localCompletion === "undone" && t.done) return false;
      return true;
    }).length;
  }, [tasks, localDifficulties, localCompletion]);

  const hasChanges =
    JSON.stringify(localDifficulties) !== JSON.stringify(selectedDifficulties) ||
    localCompletion !== completionFilter;

  const handleApply = () => {
    setDifficulties(localDifficulties);
    setCompletionFilter(localCompletion);
    close();
  };

  const handleClear = () => {
    clear();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <m.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-40"
            onClick={close}
          />

          <m.div
            initial={{ y: "100%" }}
            animate={{ y: 0 }}
            exit={{ y: "100%" }}
            transition={{ type: "spring", damping: 35, stiffness: 400 }}
            className="fixed bottom-0 left-0 right-0 z-50 bg-base-100 rounded-t-2xl shadow-2xl max-h-[85vh] flex flex-col"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-base-300">
              <h3 className="font-medium text-lg">Фильтры</h3>
              <div className="flex items-center gap-2">
                {(localDifficulties.length > 0 || localCompletion !== "all") && (
                  <button onClick={handleClear} className="btn btn-ghost btn-sm text-error">
                    Сбросить
                  </button>
                )}
                <button onClick={close} className="btn btn-ghost btn-sm">
                  <X className="w-5 h-5" />
                </button>
              </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto p-4 space-y-6">
              {/* Сложность */}
              <div>
                <label className="block text-sm font-medium text-base-content/70 mb-2">
                  Сложность
                </label>
                <DifficultyFilter
                  selected={localDifficulties}
                  onChange={setLocalDifficulties}
                />
              </div>

              {/* Статус */}
              <div>
                <label className="block text-sm font-medium text-base-content/70 mb-2">
                  Статус
                </label>
                <div className="grid grid-cols-3 gap-2">
                  {[
                    { value: "all", label: "Все" },
                    { value: "done", label: "Сделанные" },
                    { value: "undone", label: "Не сделанные" },
                  ].map(({ value, label }) => (
                    <button
                      key={value}
                      onClick={() => setLocalCompletion(value as CompletionFilter)}
                      className={`btn btn-sm h-10 ${
                        localCompletion === value ? "btn-primary" : "btn-ghost"
                      }`}
                    >
                      {label}
                    </button>
                  ))}
                </div>
              </div>
            </div>

            {/* Footer */}
            <div className="p-4 border-t border-base-300">
              <button
                onClick={handleApply}
                disabled={!hasChanges}
                className="btn btn-primary w-full h-12"
              >
                Показать {filteredCount} задач
              </button>
            </div>
          </m.div>
        </>
      )}
    </AnimatePresence>
  );
};