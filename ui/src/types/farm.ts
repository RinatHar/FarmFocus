export type DifficultyTask = "trifle" | "easy" | "normal" | "hard";
export type RaritySeed = "common" | "uncommon" | "rare" | "legendary" | "Unique";
export type TypeGood = "seed" | "bed";
export type Period = "day" | "week" | "month" | "year";

export interface ITag {
  id: number;
  name: string;
  color: string;
}

export interface ITask {
  id: number;
  title: string;
  description?: string;
  difficulty: DifficultyTask;
  done: boolean;
  date?: Date | null;
  tag?: ITag | null;
}

export interface IHabit {
  id: number;
  title: string;
  description?: string;
  difficulty: DifficultyTask;
  done: boolean;
  count: number;
  period: Period;
  every: number;
  startDate: Date;
  tag?: ITag | null;
}

export interface ISeed {
  id: number;
  name: string;
  icon: string;
  targetGrowth: number;
  rarity: RaritySeed;
  quantity: number;
  seedId: number;
}

export type SeedStorage = Omit<ISeed, "quantity">;

export interface IPlant {
  id: number;
  name: string;
  currentGrowth: number;
  targetGrowth: number;
  imgPath: string;
}

export interface IBed {
  id: number;
  plant: IPlant | null;
  isLock: boolean;
}

export interface IGood {
  type: TypeGood;
  id: number;
  quantity: number;
  cost: number;
  item: ISeed;
}
