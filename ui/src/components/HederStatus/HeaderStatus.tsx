import { CoinsCounter } from "./ui/CoinsCounter"
import { StrickCounter } from "./ui/StrickCounter"
import { XPBar } from "./ui/XP"









export const HeaderStatus = () => {
  return (
    <div className="flex justify-between items-center px-4">
      <XPBar />
      <div className="flex gap-2">
        <CoinsCounter />
        <StrickCounter />
      </div>
    </div>
  )}