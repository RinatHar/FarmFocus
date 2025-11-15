import { create } from "zustand";
import type { DifficultyTask } from "../types/farm";

export type CompletionFilter = "all" | "done" | "undone";

export const useTaskFiltersStore = create<{
  selectedDifficulties: DifficultyTask[];
  completionFilter: CompletionFilter;
  setDifficulties: (values: DifficultyTask[]) => void;
  setCompletionFilter: (filter: CompletionFilter) => void;
  clear: () => void;
}>((set) => ({
  selectedDifficulties: [],
  completionFilter: "undone",
  setDifficulties: (values) => set({ selectedDifficulties: values }),
  setCompletionFilter: (filter) => set({ completionFilter: filter }),
  clear: () => set({ selectedDifficulties: [], completionFilter: "undone" }),
}));