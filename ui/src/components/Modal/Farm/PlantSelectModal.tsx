import { m, AnimatePresence } from "framer-motion";
import { ArrowLeft, Flower2 } from "lucide-react";
import { usePlantSelectModal } from "../../../stores/usePlantSelectModal";
import { SeedCard } from "./SeedCard";
import { useFarmStore } from "../../../stores/useFarmStore";
import { memo } from "react";

export const PlantSelectModal = memo(() => {
  const { isOpen, close, selectedBedId } = usePlantSelectModal();
  const { inventorySeeds, plantSeed } = useFarmStore();

  const handleSelect = (seedId: number) => {
    if (!selectedBedId) return;
    plantSeed(selectedBedId, seedId);
    close();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <m.div
          initial={{ y: "100%" }}
          animate={{ y: 0 }}
          exit={{ y: "100%" }}
          transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
          className="fixed inset-0 z-60 flex flex-col bg-base-100"
        >
          <header className="sticky top-0 z-10 flex items-center justify-between gap-2 p-4 bg-emerald-500 text-emerald-50 rounded-t-2xl shadow-md">
            <button
              type="button"
              onClick={close}
              className="rounded-full p-1 hover:bg-emerald-600 transition-colors"
            >
              <ArrowLeft className="w-6 h-6" />
            </button>

            <h1 className="font-semibold text-xl">Выбор семян</h1>

            <div className="w-6" />
          </header>

          <main className="flex-1 overflow-y-auto p-4 scrollbar-hide">
            {inventorySeeds.length > 0 ? (
              <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
                {inventorySeeds.map((seed) => (
                  <SeedCard
                    key={seed.id}
                    seed={seed}
                    onSelect={() => handleSelect(seed.seedId)}
                  />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-center space-y-3">
                <Flower2 strokeWidth={1} size={64} className="text-base-content/80" />
                <p className="text-lg font-medium text-base-content/60">
                  У вас пока нет семян
                </p>
                <p className="text-sm text-gray-500">
                  Посетите магазин, чтобы купить новые семена!
                </p>
              </div>
            )}
          </main>

          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-8 bg-liner-to-t from-base-100 to-transparent" />
        </m.div>
      )}
    </AnimatePresence>
  );
});