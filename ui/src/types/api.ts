import type { IBed, IGood, IHabit, ISeed, ITag, ITask, SeedStorage } from "./farm";

export type ServerFarmData = {
  currentXp: number;
  coins: number;
  strick: number;
  didTaskToday: boolean;
  isDrought: boolean;

  tasks: ITask[];
  habits: IHabit[];
  tags: ITag[];
  field: IBed[];

  seeds: SeedStorage[];

  inventorySeeds: ISeed[];
  shopItem: IGood[];
};

export type DoneTaskResponse = {
  xpEarned: number;
  plantsGrown: number;
};

export type DoneHabitsResponse = {
  xpEarned: number;
  plantsGrown: number;
};

export type HarvestPlantResponse = {
  xpEarned: number;
  task_id: number;
  goldEarned: number;
};
