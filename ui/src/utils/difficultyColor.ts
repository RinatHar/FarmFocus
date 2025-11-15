


export const difficultyColor = (value?: string, type: "text" | "bg" = "bg") => {
  if (type === "text") {
    switch(value) { 
      case "trifle": return "text-slate-500";
      case "easy": return "text-emerald-500";
      case "normal": return "text-sky-500";
      case "hard": return "text-orange-500";
      default: return "text-emerald-500";
    }
  } else {
    switch(value) {
      case "trifle": return "bg-slate-500";
      case "easy": return "bg-emerald-500";
      case "normal": return "bg-sky-500";
      case "hard": return "bg-orange-500";
      default: return "bg-emerald-500";
    }
  }
};
