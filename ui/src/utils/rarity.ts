import type { RaritySeed } from "../types/farm";






export const rarityConfig: Record<
  RaritySeed,
  { label: string; light: string; dark: string; bgLight: string; bgDark: string }
> = {
  common: {
    label: "Обычное",
    light: "text-gray-600",
    dark: "text-gray-400",
    bgLight: "bg-gray-200",
    bgDark: "bg-gray-700",
  },
  uncommon: {
    label: "Необычное",
    light: "text-green-700",
    dark: "text-green-400",
    bgLight: "bg-green-200",
    bgDark: "bg-green-900",
  },
  rare: {
    label: "Редкое",
    light: "text-blue-700",
    dark: "text-blue-400",
    bgLight: "bg-blue-200",
    bgDark: "bg-blue-900",
  },
  legendary: {
    label: "Легендарное",
    light: "text-purple-700",
    dark: "text-purple-400",
    bgLight: "bg-purple-200",
    bgDark: "bg-purple-900",
  },
  Unique: {
    label: "Уникальное",
    light: "text-yellow-700",
    dark: "text-yellow-400",
    bgLight: "bg-gradient-to-r from-yellow-200 to-orange-300",
    bgDark: "bg-gradient-to-r from-yellow-200 to-orange-300",
  },
};