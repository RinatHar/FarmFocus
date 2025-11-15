import { create } from "zustand";

interface EditHabitModalStore {
  idHabit: number;
  isOpen: boolean;
  open: (id: number) => void;
  close: () => void;
  toggle: () => void;
}

export const useEditHabitModalStore = create<EditHabitModalStore>((set) => ({
  idHabit: 0,
  isOpen: false,
  open: (id) => set({ isOpen: true, idHabit: id }),
  close: () => set({ isOpen: false }),
  toggle: () => set((state) => ({ isOpen: !state.isOpen })),
}));
