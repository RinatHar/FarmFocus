import { useCallback, useMemo, useState } from "react";
import { m, AnimatePresence } from "framer-motion";
import { Plant } from "./Plant";
import { usePlantSelectModal } from "../../stores/usePlantSelectModal";
import clsx from "clsx";
import { Lock } from "lucide-react";
import { useFarmStore } from "../../stores/useFarmStore";
import type { IBed } from "../../types/farm";

type BedProps = IBed;

export const Bed = ({ id, plant, isLock }: BedProps) => {
  const { open } = usePlantSelectModal();
  const { harvestPlant } = useFarmStore();
  const [isHarvesting, setIsHarvesting] = useState(false);
  const [isClicked, setIsClicked] = useState(false);

const progress = useMemo(() => {
  if (!plant) return 0;
  return Math.min((plant.currentGrowth / plant.targetGrowth) * 100, 100);
}, [plant]);

const handleClick = useCallback(() => {
  if (isHarvesting || isLock) return;

  setIsClicked(true);
  setTimeout(() => setIsClicked(false), 200);

  if (plant && progress === 100) {
    setIsHarvesting(true);

    harvestPlant(id);

    setTimeout(() => setIsHarvesting(false), 800);
  } else if (!plant) {
    open(id);
  }
}, [isHarvesting, isLock, plant, progress, harvestPlant, id, open]);

  return (
    <m.div
      className={clsx(
        "relative flex flex-col items-center justify-center cursor-pointer select-none",
        { "cursor-not-allowed opacity-50": isLock }
      )}
      onClick={handleClick}
      whileTap={{ scale: isLock ? 1 : 0.9 }}
      animate={isClicked ? { scale: isLock ? 1 : 0.9 } : { scale: 1 }}
      transition={{
        type: "spring",
        stiffness: 400,
        damping: 15,
        duration: 0.15,
      }}
    >
      {isLock && (
        <div className="absolute inset-0 flex items-center justify-center z-20">
          <Lock />
        </div>
      )}

      <AnimatePresence>
        {isHarvesting && (
          <m.div
            key="pulse"
            initial={{ scale: 1, opacity: 0.3 }}
            animate={{ scale: 1.3, opacity: 0 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.6, ease: "easeOut" }}
            className="absolute inset-0 rounded-full bg-emerald-300/30 z-0"
          />
        )}
      </AnimatePresence>

      {/* Прогресс-бар */}
      <AnimatePresence>
        {plant && !isHarvesting && (
          <m.div
            initial={{ opacity: 0, y: -5 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.8 }}
            transition={{ duration: 0.3 }}
            className="absolute bottom-6 left-1/2 transform -translate-x-1/2 w-16 h-1.5 bg-gray-700/15 rounded-full overflow-hidden shadow-sm z-10"
          >
            <m.div
              className={clsx(
                "h-full bg-emerald-400/60 transition-all duration-500",
                { "bg-emerald-400/60": progress < 100 },
                { "animate-pulse bg-emerald-500": progress === 100 }
              )}
              style={{ width: `${progress}%` }}
              exit={{ scaleX: 0, originX: 0 }}
              transition={{ duration: 0.3 }}
            />
          </m.div>
        )}
      </AnimatePresence>

      {/* Грядка */}
      <m.img
        src="/assets/farm/bed.png"
        alt={`bed-${id}`}
        className="w-[90%] h-[90%] object-contain"
        style={{ imageRendering: "pixelated" }}
        animate={
          isHarvesting
            ? { scale: [1, 1.05, 1], filter: ["brightness(1)", "brightness(1.3)", "brightness(1)"] }
            : { scale: 1, filter: "brightness(1)" }
        }
        transition={{ duration: 0.6 }}
      />

      {/* Растение */}
      <AnimatePresence mode="wait">
        {plant && !isHarvesting ? (
          <m.div
            key="plant"
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{
              rotate: 90,
              y: 30,
              opacity: 0,
            }}
            transition={{
              duration: 0.6,
              ease: "easeIn",
            }}
            className="absolute inset-0 flex items-center justify-center"
          >
            <Plant plant={plant} />
          </m.div>
        ) : null}
      </AnimatePresence>
    </m.div>
  );
};
