/**
 * Модуль для расчёта уровня и прогресса в игре
 * - Базовый рост: +5%
 * - Каждые 10 уровней: +1% к росту (5% → 6% → ... → 15% max)
 */

export class LevelCalculator {
  private static readonly START_THRESHOLD = 50; // exp для 2-го уровня
  private static readonly BASE_GROWTH = 1.05;   // +5%
  private static readonly MAX_BONUS = 0.10;      // +10% → итого 15%

  /**
   * Рассчитывает текущий уровень по опыту
   */
  static calculateLevel(exp: number): number {
    if (exp <= 0) return 1;

    let level = 1;
    let threshold = this.START_THRESHOLD;
    let remaining = Number(exp);

    while (remaining >= threshold && level < 9999) {
      remaining -= threshold;
      level++;

      const bonus = Math.min((level - 1) / 10 * 0.01, this.MAX_BONUS);
      const growth = this.BASE_GROWTH + bonus;
      threshold *= growth;

      if (threshold > 1e18) break;
    }

    return Math.min(level, 9999);
  }

  /**
   * Общий опыт, нужный для достижения уровня (ДО уровня)
   */
  static experienceForLevel(level: number): number {
    if (level <= 1) return 0;

    let total = 0;
    let threshold = this.START_THRESHOLD;

    for (let i = 1; i < level; i++) {
      total += Math.round(threshold);
      const bonus = Math.min(i / 10 * 0.01, this.MAX_BONUS);
      const growth = this.BASE_GROWTH + bonus;
      threshold *= growth;

      if (total > Number.MAX_SAFE_INTEGER) {
        return Number.MAX_SAFE_INTEGER;
      }
    }

    return total;
  }

  /**
   * Сколько опыта осталось до следующего уровня
   */
  static experienceToNextLevel(currentLevel: number, currentExp: number): number {
    if (currentLevel >= 9999) return 0;

    const expForCurrent = this.experienceForLevel(currentLevel);
    const expForNext = this.experienceForLevel(currentLevel + 1);
    const needed = expForNext - expForCurrent;
    const have = currentExp - expForCurrent;

    if (have >= needed) return 0;
    return needed - have;
  }

  /**
   * Прогресс к следующему уровню (0.0 - 1.0)
   */
  static progressToNext(currentLevel: number, currentExp: number): number {
    if (currentLevel >= 9999) return 1.0;

    const expForCurrent = this.experienceForLevel(currentLevel);
    const expForNext = this.experienceForLevel(currentLevel + 1);
    const have = Number(currentExp - expForCurrent);
    const need = Number(expForNext - expForCurrent);

    if (have >= need) return 1.0;
    return Math.max(have / need, 0);
  }

  static totalXpForNextLevel(currentLevel: number): number {
    const expForCurrent = this.experienceForLevel(currentLevel);
    const expForNext = this.experienceForLevel(currentLevel + 1);
    return expForNext - expForCurrent;
  }

}