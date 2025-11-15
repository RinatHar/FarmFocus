import { create } from "zustand";

interface EditTaskModalStore {
  idTask: number;
  isOpen: boolean;
  open: (id: number) => void;
  close: () => void;
  toggle: () => void;
}

export const useEditTaskModalStore = create<EditTaskModalStore>((set) => ({
  idTask: 0,
  isOpen: false,
  open: (id) => set({ isOpen: true, idTask: id }),
  close: () => set({ isOpen: false }),
  toggle: () => set((state) => ({ isOpen: !state.isOpen })),
}));
