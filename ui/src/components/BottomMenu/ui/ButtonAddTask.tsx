import { Plus } from "lucide-react";
import { m } from "framer-motion";
import { useAddTaskModalStore } from "../../../stores/addTaskModal";
import type { SwiperPage } from "../../../App";
import { useAddHabitModalStore } from "../../../stores/addHabitModal";
import { useFarmStore } from "../../../stores/useFarmStore";

type Props = {
  currentPage: SwiperPage;
};

export const ButtonAddTask = ({ currentPage }: Props) => {
  const { isDrought } = useFarmStore();
  const { open: openTaskAdd } = useAddTaskModalStore();
  const { open: openHabitAdd } = useAddHabitModalStore();

  const handleOpenModalAdd = () => {
    if (currentPage === "habit-list") {
      openHabitAdd();
    } else {
      openTaskAdd();
    }
  };

  const iconColor = isDrought ? "text-amber-500" : "text-lime-500";
  const gradientFrom = isDrought ? "from-amber-200" : "from-emerald-200";
  const gradientVia = isDrought ? "via-orange-200" : "via-green-200";
  const gradientTo = isDrought ? "to-yellow-200" : "to-lime-200";
  const glowColor = isDrought ? "bg-amber-300/50" : "bg-emerald-300/50";

  return (
    <m.button
      onClick={handleOpenModalAdd}
      className="relative group"
      whileTap={{ scale: 0.92 }}
      whileHover={{ scale: 1.08 }}
    >
      <m.div
        className={`absolute inset-0 rounded-full ${glowColor} blur-xl`}
        animate={{
          scale: [1, 1.3, 1],
          opacity: [0.4, 0.2, 0.4],
        }}
        transition={{
          duration: 2,
          repeat: Infinity,
          ease: "easeInOut",
        }}
      />

      <div
        className={`relative w-12 h-12 sm:text-2xl xs:w-16 xs:h-16 rounded-full shadow-lg overflow-hidden p-1
          bg-gradient-to-br ${gradientFrom} ${gradientVia} ${gradientTo}`}
      >
        <div className="w-full h-full rounded-full bg-base-100/80 flex items-center justify-center">
          <Plus className={iconColor} />
        </div>
      </div>
    </m.button>
  );
};