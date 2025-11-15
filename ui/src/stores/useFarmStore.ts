import { create } from "zustand";
import type { DifficultyTask, IBed, IGood, IHabit, ISeed, ITag, ITask, Period, TypeGood } from "../types/farm";
import type { ServerFarmData } from "../types/api";
import { LevelCalculator } from "../utils/levelCalculator";
import { fetchAddTasks, fetchDeleteTasks, fetchDoneTask, fetchEditTasks, fetchUndoneTask } from "../api/tasks";
import { showRewardToast } from "../utils/showRewardToast";
import { showErrorToast } from "../utils/showErrorToast";
import { fetchDeleteTag } from "../api/tags";
import { fetchAddHabit, fetchDeleteHabit, fetchDoneHabit, fetchEditHabit, fetchUndoneHabit } from "../api/habit";
import { fetchHarvestPlant, fetchSetPlant } from "../api/planet";


export interface FarmState {

  setFromServer: (data: ServerFarmData) => void;
  
  // --- Игрок ---
  userId: number;
  currentXp: number;
  currentLevelXp: number;
  countXpLvLUp: number;
  currentLevel: number;
  coins: number;
  strick: number;
  didTaskToday: boolean;
  isDrought: boolean;

  setUserId: (id: number) => void;

  setDrought: (isDrought: boolean) => void;

  setCoins: (coins: number) => void;
  addCoins: (coins: number) => void;
  removeCoins: (coins: number) => void;

  setXP: (xp: number) => void;
  addXP: (xp: number) => void;
  removeXP: (xp: number) => void;
  
  // --- Задачи ---
  tasks: ITask[];
  toggleTask: (id: number) => void;
  addTask: (title: string, description: string, difficulty: DifficultyTask, startDate: Date | null, tag?: ITag | null) => void;
  editTask: (id: number, title: string, description: string, difficulty: DifficultyTask, startDate: Date | null, tag?: ITag | null) => void;
  removeTask: (id: number) => void;

  // --- Дела ---
  habits: IHabit[];
  toggleHabit: (id: number) => void;
  addHabit: (title: string, description: string, difficulty: DifficultyTask, period: Period, every: number, startDate: Date, tag?: ITag | null) => void;
  editHabit: (id: number, title: string, description: string, difficulty: DifficultyTask, period: Period, every: number, startDate: Date, tag?: ITag | null) => void;
  removeHabit: (id: number) => void;

  // --- Теги ---
  tags: ITag[];
  createTag: (name: string, color: string) => void;
  updateTag: (id: number, name: string, color: string) => void;
  deleteTag: (id: number) => void;

  // --- Ферма ---
  rows: number;
  cols: number;
  field: IBed[];
  unlockNextBed: () => { success: boolean; message?: string };

  // --- Растения ---
  setPlantInBed: (bedId: number, plantId: number) => Promise<void>;
  harvestPlant: (bedId: number) => void;
  advanceGrowth: () => void;

  // --- Семена (шаблоны + инвентарь) ---
  inventorySeeds: ISeed[];
  plantSeed: (bedId: number, seedId: number) => { success: boolean; message?: string };
  addSeedToInventory: (seed: ISeed) => void;

  // --- Магазин ---
  shopItem: IGood[];
  decreaseShopItem: (type: TypeGood, idGood: number) => { success: boolean; message?: string };
  removeShopItem: (type: TypeGood, idGood: number) => { success: boolean; message?: string };
}

const ROWS = 3;
const COLS = 3;

// =============================================
// ИНИЦИАЛИЗАЦИЯ СТОРА
// =============================================

export const useFarmStore = create<FarmState>((set, get) => ({
  setFromServer: (data: ServerFarmData) => set((state) => {
  const {
    currentXp = 0,
    coins = 0,
    strick = 0,
    didTaskToday = false,
    isDrought = false,
  } = data;

  const currentLevel = LevelCalculator.calculateLevel(data.currentXp);
  const countXpLvLUp = LevelCalculator.totalXpForNextLevel(currentLevel);
  const currentLevelXp = currentXp - LevelCalculator.experienceForLevel(currentLevel);

  const tasks: ITask[] = (data.tasks || [])
    .map((t) => ({
      id: t.id,
      title: t.title,
      description: t.description,
      difficulty: t.difficulty,
      done: t.done,
      date: t.date ? new Date(t.date) : null,
      tag: t.tag,
    }));


  const habits: IHabit[] = (data.habits || []).map((h) => ({
    id: h.id,
    title: h.title,
    description: h.description,
    difficulty: h.difficulty,
    done: h.done,
    count: h.count || 0,
    period: h.period as Period,
    every: h.every,
    startDate: new Date(h.startDate),
    tag: h.tag,
  }));

  const tags: ITag[] = data.tags || [];

  const field: IBed[] = data.field || state.field;

  const inventorySeeds = data.inventorySeeds?.length ? data.inventorySeeds : state.inventorySeeds;
  const shopItem = data.shopItem?.length ? data.shopItem : state.shopItem;

  return {
    currentXp,
    countXpLvLUp,
    currentLevel,
    currentLevelXp,
    coins,
    strick,
    didTaskToday,
    isDrought,
    tasks,
    habits,
    tags,
    field,
    inventorySeeds,
    shopItem,
  };
}),


  // ────────────────────────────────────────
  // ИГРОК
  // ────────────────────────────────────────
  userId: 0,
  currentXp: 0,
  currentLevelXp: 0,
  countXpLvLUp: 0,
  currentLevel: 1,
  coins: 0,
  strick: 0,
  didTaskToday: false,
  isDrought: false,

  setDrought(isDrought) {
    set({ isDrought });
  },

  setUserId: (id) => set({ userId: id }),

  setCoins: (coins) => set({ coins }),
  addCoins: (coins) => set({ coins: get().coins + coins }),
  removeCoins: (coins) => {
    const current = get().coins;
    const newCoins = current > coins ? current - coins : 0; // не ниже 0
    get().setCoins(newCoins);
  },

  setXP: (xp) => set((state) => {
      const currentLevel = LevelCalculator.calculateLevel(xp);
      const currentLevelXp = xp - LevelCalculator.experienceForLevel(currentLevel);
      
      if (state.currentLevel != currentLevel){
        const countXpLvLUp = LevelCalculator.totalXpForNextLevel(currentLevel);
        return({...state, currentXp: xp, currentLevel, countXpLvLUp, currentLevelXp})
      }

      return({...state, currentXp: xp, currentLevel, currentLevelXp})
    }),
  addXP: (xp) => {
    const newXp = get().currentXp + xp;
    get().setXP(newXp);
  },
  removeXP: (xp) => {
    const current = get().currentXp;
    const newXp = current > xp ? current - xp : 0;
    get().setXP(newXp);
  },

  // ────────────────────────────────────────
  // ЗАДАЧИ
  // ────────────────────────────────────────
  tasks: [],

  toggleTask: (id) => {
    const state = get();
    const task = state.tasks.find((t) => t.id === id);
    if (!task) return;

    state.setDrought(false);

    const wasDone = task.done;
    const newDone = !wasDone;

    const newTasks = state.tasks.map((t) =>
      t.id === id ? { ...t, done: newDone } : t
    );

    const newStrick = state.didTaskToday ? state.strick : state.strick + 1;

    set({
      tasks: newTasks,
      strick: newStrick,
      didTaskToday: true,
    });

    const request = newDone ? fetchDoneTask(id) : fetchUndoneTask(id);

    request.then((result) => {
      const { xpEarned, plantsGrown } = result;

      get().addXP(xpEarned);
      showRewardToast(xpEarned);

      if (plantsGrown > 0) {
        get().advanceGrowth();
      }

    }).catch(() => {
      set({
        tasks: state.tasks.map((t) =>
          t.id === id ? { ...t, done: wasDone } : t
        ),
        strick: state.strick,
        didTaskToday: state.didTaskToday,
      });
      showErrorToast("Не удалось выполнить задачу");
    });
  },

  addTask: (title, description, difficulty, date, tag) => {
    const state = get();

    const tempId = Date.now();
    const tempTask: ITask = {
      id: tempId,
      title,
      description,
      difficulty,
      date,
      tag: tag ?? undefined,
      done: false,
    };

    set({
      tasks: [...state.tasks, tempTask],
    });

    const taskNoId = { title, description, difficulty, date, tag };
    
    fetchAddTasks(taskNoId)
    .then((serverTask) => {
      if (!serverTask) throw new Error("No server response");

      set((current) => ({
        tasks: current.tasks.map((t) =>
          t.id === tempId
            ? { ...t, id: serverTask.id }
            : t
        ),
      }));
    })
    .catch(() => {
      set((current) => ({
        tasks: current.tasks.filter((t) => t.id !== tempId),
      }));
      showErrorToast("Не удалось добавить задачу");
    });
  },

  editTask: (id: number, title: string, description: string, difficulty: DifficultyTask, date: Date | null, tag?: ITag | null) => {
    const state = get();

    set((s) => ({
      tasks: s.tasks.map((t) =>
        t.id === id
          ? { ...t, title, description, difficulty, date, tag }
          : t
      ),
    }));

    const updatedTask: ITask = {
      ...state.tasks.find(t => t.id === id)!,
      title,
      description,
      difficulty,
      date,
      tag: tag ?? undefined,
    };

    fetchEditTasks(updatedTask)
      .catch(() => {
        const originalTask = state.tasks.find((t) => t.id === id);
        if (originalTask) {
          set((s) => ({
            tasks: s.tasks.map((t) =>
              t.id === id ? originalTask : t
            ),
          }));
        }
        showErrorToast("Не удалось сохранить изменения");
      });
  },

  removeTask: (id: number) => {
    const state = get();

    set((s) => ({
      tasks: s.tasks.filter((t) => t.id !== id),
    }));


    fetchDeleteTasks(id).catch(() => {
      const task = state.tasks.find((t) => t.id === id);
      if (task) {
        set((s) => ({
          tasks: [...s.tasks, task].sort((a, b) => a.id - b.id),
        }));
      }
      showErrorToast("Не удалось удалить задачу");
    });
  },

  // ────────────────────────────────────────
  // Дела
  // ────────────────────────────────────────

  habits: [],

  toggleHabit: (id) => {
    const state = get();
    const habit = state.habits.find((h) => h.id === id);
    if (!habit) return;

    const wasDone = habit.done;
    const newDone = !wasDone;

    state.setDrought(false);

    const newHabits = state.habits.map((t) =>
      t.id === id ? { ...t, done: newDone, count: newDone ? t.count + 1 : t.count - 1 } : t
    );

    const newStrick = state.didTaskToday ? state.strick : state.strick + 1;

    set({
      habits: newHabits,
      strick: newStrick,
      didTaskToday: true,
    });

    const request = newDone ? fetchDoneHabit(id) : fetchUndoneHabit(id);

    request.then((result) => {
      const { xpEarned, plantsGrown } = result;

      get().addXP(xpEarned);
      showRewardToast(xpEarned);

      if (plantsGrown > 0) {
        get().advanceGrowth();
      }

    }).catch(() => {
      set({
        habits: state.habits.map((h) =>
          h.id === id ? { ...h, done: wasDone } : h
        ),
        strick: state.strick,
        didTaskToday: state.didTaskToday,
      });
      showErrorToast("Не удалось выполнить задачу");
    });
  },

  addHabit: (title, description, difficulty, period, every, startDate, tag) => {
    const state = get();

    const tempId = Date.now();
    const tempHabit: IHabit = {
      id: tempId,
      title,
      description,
      difficulty,
      period,
      every,
      startDate,
      tag: tag ?? undefined,
      done: false,
      count: 0,
    };

    set({
      habits: [...state.habits, tempHabit],
    });

    const habitNoId = { title, description, difficulty, startDate, tag, every, period };
    
    fetchAddHabit(habitNoId)
    .then((serverTask) => {
      if (!serverTask) throw new Error("No server response");

      set((current) => ({
        habits: current.habits.map((h) =>
          h.id === tempId
            ? { ...h, id: serverTask.id }
            : h
        ),
      }));
    })
    .catch(() => {
      set((current) => ({
        habits: current.habits.filter((h) => h.id !== tempId),
      }));
      showErrorToast("Не удалось добавить задачу");
    });
  },

  editHabit: (id, title, description, difficulty, period, every, startDate, tag) => {
    const state = get();

    set((s) => ({
      habits: s.habits.map((h) =>
        h.id === id
          ? { ...h, title, description, difficulty, period, every, startDate, tag }
          : h
      ),
    }));

    const updatedHabit: IHabit = {
      ...state.habits.find(h => h.id === id)!,
      title,
      description,
      difficulty,
      startDate,
      tag: tag ?? undefined,
      period,
      every,
    };

    fetchEditHabit(updatedHabit)
      .catch(() => {
        const originalHabit = state.habits.find((h) => h.id === id);
        if (originalHabit) {
          set((s) => ({
            tasks: s.habits.map((h) =>
              h.id === id ? originalHabit : h
            ),
          }));
        }
        showErrorToast("Не удалось сохранить изменения");
      });
  },

  removeHabit: (id) => {
    const state = get();

    set((s) => ({
      habits: s.habits.filter((h) => h.id !== id),
    }));


    fetchDeleteHabit(id).catch(() => {
      const habit = state.habits.find((t) => t.id === id);
      if (habit) {
        set((s) => ({
          habits: [...s.habits, habit].sort((a, b) => a.id - b.id),
        }));
      }
      showErrorToast("Не удалось удалить задачу");
    });
  },

  // ────────────────────────────────────────
  // ТЕГИ
  // ────────────────────────────────────────
  tags: [
    { id: 1, name: "Дом", color: "#10b981" },
    { id: 2, name: "Работа", color: "#3b82f6" },
    { id: 3, name: "Учёба", color: "#ec4899" },
  ],

  createTag: (name, color) =>
    set((state) => {
      const newId = Math.max(...state.tags.map((t) => t.id), 0) + 1;
      return { tags: [...state.tags, { id: newId, name, color }] };
    }),

  updateTag: (id, name, color) =>
    set((state) => ({
      tags: state.tags.map((t) => (t.id === id ? { ...t, name, color } : t)),
    })),

  deleteTag: (id: number) => {
    const state = get();

    set((s) => ({
      tags: s.tags.filter((t) => t.id !== id),
      tasks: s.tasks.map((task) =>
        task.tag?.id === id ? { ...task, tag: undefined } : task
      ),
    }));


    fetchDeleteTag(id)
      .catch(() => {
        const deletedTag = state.tags.find((t) => t.id === id);
        if (deletedTag) {
          set((s) => ({
            tags: [...s.tags, deletedTag].sort((a, b) => a.id - b.id),
            tasks: s.tasks.map((task) => {
              const originalTask = state.tasks.find((t) => t.id === task.id);
              return originalTask?.tag?.id === id
                ? { ...task, tag: deletedTag }
                : task;
            }),
          }));
        }
        showErrorToast("Не удалось удалить тег");
      });
  },

  // ────────────────────────────────────────
  // ФЕРМА (ГРЯДКИ)
  // ────────────────────────────────────────
  rows: ROWS,
  cols: COLS,

  field: Array.from({ length: ROWS * COLS }, (_, i) => ({
    id: i + 1,
    plant: null,
    isLock: i > 0, // Только первая грядка открыта
  })),

  unlockNextBed: () => {
    const state = get();
    const lockedBedIndex = state.field.findIndex((bed) => bed.isLock);

    if (lockedBedIndex === -1) {
      return { success: false, message: "Все грядки уже открыты" };
    }

    const updatedField = [...state.field];
    updatedField[lockedBedIndex] = {
      ...updatedField[lockedBedIndex],
      isLock: false,
    };

    set({
      field: updatedField,
    });

    return { success: true };
  },

  // ────────────────────────────────────────
  // РАСТЕНИЯ
  // ────────────────────────────────────────

  setPlantInBed: async (bedId, plantId) => {
    try {
      const serverPlant = await fetchSetPlant(bedId, plantId);
      if (!serverPlant) throw new Error();

      set((s) => ({
        field: s.field.map(b =>
          b.id === bedId ? { ...b, plant: serverPlant } : b
        ),
      }));
    } catch {
      showErrorToast("Не удалось посадить");
    }
  },

  harvestPlant: async (bedId) => {
    const state = get();

    const bed = state.field.find((b) => b.id === bedId);
    if (!bed?.plant) return { success: false };

    const isGrown = bed.plant.currentGrowth >= bed.plant.targetGrowth;
    if (!isGrown) return { success: false };

    const oldField = state.field;
    set({
      field: state.field.map((b) =>
        b.id === bedId ? { ...b, plant: null } : b
      ),
    });

    try {
      const { xpEarned, goldEarned } = await fetchHarvestPlant(bed.plant.id);

      get().addXP(xpEarned);
      get().addCoins(goldEarned);

      showRewardToast(xpEarned, goldEarned);

    } catch {
      set({ field: oldField });
      showErrorToast("Не удалось собрать урожай");
    }
  },

  advanceGrowth: () => {
    set((state) => {
      if (state.isDrought) {
        return state;
      }

      return {
        field: state.field.map((bed) => {
          if (!bed.plant || bed.plant.currentGrowth >= bed.plant.targetGrowth) {
            return bed;
          }

          return {
            ...bed,
            plant: {
              ...bed.plant,
              currentGrowth: bed.plant.currentGrowth + 1,
            },
          };
        }),
      };
    });
  },

  // ────────────────────────────────────────
  // СЕМЕНА
  // ────────────────────────────────────────

  inventorySeeds: [
    { id: 1, name: "Пшеница", icon: "/assets/seeds/wheat.png", targetGrowth: 4, rarity: "common", seedId: 1, quantity: 1  },
    { id: 2, name: "Баклажан", icon: "/assets/seeds/aubergine.png", targetGrowth: 8, rarity: "uncommon", seedId: 2, quantity: 1 },
  ],

  plantSeed: (bedId, seedId) => {
    const state = get();

    const seed = state.inventorySeeds.find(s => s.seedId === seedId);
    const bed = state.field.find(b => b.id === bedId);

    if (!seed || seed.quantity <= 0) return { success: false, message: "Нет семян" };
    if (!bed || bed.plant || bed.isLock) return { success: false, message: "Грядка занята" };

    get().setPlantInBed(bedId, seedId);

    set((s) => ({
      inventorySeeds: s.inventorySeeds
        .map(s => s.seedId === seedId ? { ...s, quantity: s.quantity - 1 } : s)
        .filter(s => s.quantity > 0),
    }));

    return { success: true };
  },

  addSeedToInventory: (seed) => {
    set((state) => {
      const existingSeedIndex = state.inventorySeeds.findIndex(
        (item) => item.id === seed.id
      );

      if (existingSeedIndex !== -1) {
        const updatedSeeds = [...state.inventorySeeds];
        updatedSeeds[existingSeedIndex] = {
          ...updatedSeeds[existingSeedIndex],
          quantity: updatedSeeds[existingSeedIndex].quantity + 1,
        };
        
        return {
          inventorySeeds: updatedSeeds,
        };
      } else {

        return {
          inventorySeeds: [...state.inventorySeeds, seed],
        };
      }
    });
  },

  // ────────────────────────────────────────
  // МАГАЗИН
  // ────────────────────────────────────────
  shopItem: [],

  decreaseShopItem: (type: TypeGood, idGood: number) => {
    const state = get();
    const itemIndex = state.shopItem.findIndex(
      (item) => item.type === type && item.id === idGood
    );

    if (itemIndex === -1) {
      return { success: false, message: "Товар не найден в магазине" };
    }

    const item = state.shopItem[itemIndex];

    if (item.quantity <= 1) {
      set({
        shopItem: state.shopItem.filter((_, index) => index !== itemIndex),
      });
    } else {
      const updatedShop = [...state.shopItem];
      updatedShop[itemIndex] = {
        ...item,
        quantity: item.quantity - 1,
      };

      set({ shopItem: updatedShop });
    }

    return { success: true };
  },

  removeShopItem: (type: TypeGood, idGood: number) => {
    const state = get();
    const exists = state.shopItem.some(
      (item) => item.type === type && item.id === idGood
    );

    if (!exists) {
      return { success: false, message: "Товар не найден" };
    }

    set({
      shopItem: state.shopItem.filter(
        (item) => !(item.type === type && item.id === idGood)
      ),
    });

    return { success: true };
  },

}));