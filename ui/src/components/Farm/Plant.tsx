import { m } from "framer-motion";
import type { IPlant } from "../../types/farm";

interface PlantProps {
  plant?: IPlant | null;
}

export const Plant = ({ plant }: PlantProps) => {
  if (!plant) return null;

  let imageNumber: number;

  if (plant.currentGrowth >= plant.targetGrowth) {
    imageNumber = 4;
  } else {
    const progress = plant.currentGrowth / plant.targetGrowth;
    const stageIndex = Math.floor(progress * 3);
    imageNumber = stageIndex + 1;
  }

  const imgSrc = `${plant.imgPath.replace(/\/$/, '')}/state${imageNumber}.png`;

  return (
    <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
      <m.div
        animate={{
          rotate: [1, -1, 1, -1, 1],
        }}
        transition={{
          duration: 5,
          repeat: Infinity,
          ease: "easeInOut",
        }}
        style={{ originY: "90%", originX: "50%" }}
        className="w-full h-full mb-8 flex items-center justify-center"
      >
        <m.img
          src={imgSrc}
          alt={`${plant.name} - stage ${imageNumber}`}
          className="w-[70%] h-[70%] object-contain"
          style={{ imageRendering: "pixelated" }}
          key={plant.id + plant.currentGrowth}
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.4 }}
        />
      </m.div>
    </div>
  );
};