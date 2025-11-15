import { memo } from "react";
import type { Period } from "../../types/farm";

interface PeriodSummaryProps {
  period: Period;
  every: number;
}

const pluralRules = new Intl.PluralRules("ru-RU");

const getPluralForm = (number: number, forms: [string, string, string]): string => {
  const rule = pluralRules.select(number);
  if (rule === "one") return forms[0];
  if (rule === "few") return forms[1];
  return forms[2];
};

export const PeriodSummary = memo(({ period, every }: PeriodSummaryProps) => {
  if (every <= 0) return null;

const periodNames: Record<
  Period,
  [string, string, string]
> = {
  day: ["день", "дня", "дней"],
  week: ["неделю", "недели", "недель"],
  month: ["месяц", "месяца", "месяцев"],
  year: ["год", "года", "лет"],
};

  const periodForms = periodNames[period];
  const everyText = every === 1 ? "каждый" : "каждые";

  const periodText = getPluralForm(every, periodForms);

  return (
    <p className="text-sm text-base-content/70 mt-1 px-4">
      Повторять {everyText} {every > 1 && `${every} `}{periodText}
    </p>
  );
});

PeriodSummary.displayName = "PeriodSummary";