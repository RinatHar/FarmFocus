import clsx from "clsx";
import type { LucideIcon } from "lucide-react";

type BottomMenuButtonProps = {
  icon: LucideIcon;
  label: string;
  page: string;
  currentPage: string;
  isDrought: boolean;
  onClick?: () => void;
};

export const BottomMenuButton = ({
  icon: Icon,
  label,
  page,
  currentPage,
  isDrought,
  onClick,
}: BottomMenuButtonProps) => {
  const isActive = currentPage === page;

  return (
    <button
      className="flex flex-col gap-1 justify-center items-center"
      onClick={onClick}
    >
      <Icon
        className={clsx({
          "text-base-content/80": !isActive,
          "text-lime-500 scale-110": isActive && !isDrought,
          "text-amber-500 scale-110": isActive && isDrought,
        })}
        strokeWidth={isActive ? 2.5 : 2}
      />
      <span
        className={clsx("text-xs font-mono", {
          "text-base-content/80": !isActive,
          "text-lime-500": isActive && !isDrought,
          "text-amber-500": isActive && isDrought,
        })}
      >
        {label}
      </span>
    </button>
  );
};