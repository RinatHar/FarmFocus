import { create } from "zustand";

interface PlantSelectModalStore {
  isOpen: boolean;
  selectedBedId: number | null;
  open: (id: number) => void;
  close: () => void;
}

export const usePlantSelectModal = create<PlantSelectModalStore>((set) => ({
  isOpen: false,
  selectedBedId: null,
  open: (id) => set({ isOpen: true, selectedBedId: id }),
  close: () => set({ isOpen: false, selectedBedId: null }),
}));
