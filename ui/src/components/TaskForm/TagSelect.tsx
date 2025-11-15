import { useState, useMemo, useRef, useEffect, useCallback, memo } from "react";
import { m, AnimatePresence } from "framer-motion";
import { Plus, Edit2, Trash2, Check, X, ChevronDown } from "lucide-react";
import clsx from "clsx";
import { useFarmStore } from "../../stores/useFarmStore";
import type { DifficultyTask, ITag } from "../../types/farm";
import { fetchCreatedTag, fetchUpdateTag } from "../../api/tags";

// Цвета для тегов
const COLORS = ["#ef4444", "#ec4899", "#8b5cf6", "#6366f1", "#3b82f6", "#10b981", "#f59e0b"] as const;

type TagSelectProps = {
  selectedTag?: ITag | null;
  difficulty?: DifficultyTask;
  onSelect: (tag: ITag | undefined) => void;
};

// Мемоизированный компонент одного тега
const TagItem = memo(
  ({
    tag,
    isEditing,
    editName,
    editColor,
    onSelect,
    onEdit,
    onDelete,
    onChangeEditName,
    onChangeEditColor,
    onSaveEdit,
    onCancelEdit,
  }: {
    tag: ITag;
    isEditing: boolean;
    editName: string;
    editColor: string;
    onSelect: (tag: ITag) => void;
    onEdit: (tag: ITag) => void;
    onDelete: (id: number) => void;
    onChangeEditName: (name: string) => void;
    onChangeEditColor: (color: string) => void;
    onSaveEdit: () => void;
    onCancelEdit: () => void;
  }) => {
    return (
      <div
        className={clsx(
          "flex items-center justify-between py-3 px-3 rounded-xl transition hover:bg-base-200",
          isEditing && "bg-base-200"
        )}
      >
        {isEditing ? (
          <div className="flex flex-col gap-3 w-full">
            <input
              value={editName}
              onChange={(e) => onChangeEditName(e.target.value)}
              className="input input-bordered w-full h-12 rounded-xl text-base focus:outline-none"
              autoFocus
            />
            <div className="flex items-center justify-between">
              <div className="flex gap-2 flex-wrap" onClick={(e) => e.stopPropagation()}>
                {COLORS.map((c) => (
                  <button
                    key={c}
                    type="button"
                    onClick={() => onChangeEditColor(c)}
                    style={{ backgroundColor: c }}
                    className={clsx(
                      "w-10 h-10 rounded-full ring-2 transition",
                      editColor === c ? "ring-base-content/80 shadow-lg scale-90" : "ring-transparent"
                    )}
                  />
                ))}
              </div>
              <div className="flex gap-2">
                <button onClick={onSaveEdit} className="btn btn-success btn-sm h-10">
                  <Check className="w-4 h-4" />
                </button>
                <button onClick={onCancelEdit} className="btn btn-error btn-sm h-10">
                  <X className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div
            className="flex justify-between items-center w-full cursor-pointer"
            onClick={() => onSelect(tag)}
          >
            <div className="flex items-center gap-3">
              <div className="w-6 h-6 rounded-full border border-base-300" style={{ backgroundColor: tag.color }} />
              <span className="font-medium text-base">{tag.name}</span>
            </div>
            <div className="flex gap-2" onClick={(e) => e.stopPropagation()}>
              <button type="button" className="btn btn-ghost btn-sm h-9 w-9 p-0" onClick={() => onEdit(tag)}>
                <Edit2 className="w-4 h-4" />
              </button>
              <button
                type="button"
                className="btn btn-ghost btn-sm h-9 w-9 p-0 text-error"
                onClick={() => onDelete(tag.id)}
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }
);

export const TagSelect = ({ selectedTag, difficulty, onSelect }: TagSelectProps) => {
  const { tags, createTag, updateTag, deleteTag } = useFarmStore();
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");
  const [editColor, setEditColor] = useState<string>(COLORS[0]);
  const [isCreating, setIsCreating] = useState(false);

  const containerRef = useRef<HTMLDivElement>(null);
  const createInputRef = useRef<HTMLInputElement>(null);

  const filteredTags = useMemo(
    () => tags.filter((t) => t.name.toLowerCase().includes(search.toLowerCase())),
    [tags, search]
  );

  const closeModal = useCallback(() => {
    setIsOpen(false);
    setIsCreating(false);
    setEditingId(null);
    setEditName("");
    setEditColor(COLORS[0]);
  }, []);

  const handleCreate = useCallback(() => {
    if (!editName.trim()) return;
    createTag(editName, editColor);
    setEditName("");
    setEditColor(COLORS[0]);
    setIsCreating(false);

    fetchCreatedTag({name: editName, color: editColor})
  }, [createTag, editName, editColor]);

  const startEdit = useCallback((tag: ITag) => {
    setEditingId(tag.id);
    setEditName(tag.name);
    setEditColor(tag.color);
    setIsCreating(false);
  }, []);

  const saveEdit = useCallback(() => {
    if (!editName.trim() || editingId === null) return;
    updateTag(editingId, editName, editColor);
    setEditingId(null);
    setEditName("");
    setEditColor(COLORS[0]);

    fetchUpdateTag({id: editingId, name: editName, color: editColor})
  }, [editName, editColor, editingId, updateTag]);

  const cancelEdit = useCallback(() => {
    setEditingId(null);
    setEditName("");
    setEditColor(COLORS[0]);
  }, []);

  const handleSelect = useCallback(
    (tag: ITag) => {
      onSelect(tag);
      closeModal();
    },
    [onSelect, closeModal]
  );

  useEffect(() => {
    if (!isOpen) return;
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) closeModal();
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen, closeModal]);

  useEffect(() => {
    if (isCreating && createInputRef.current) createInputRef.current.focus();
  }, [isCreating]);

  return (
    <div ref={containerRef} className="relative w-full">
      {/* Кнопка открытия */}
      <button
        type="button"
        onClick={() => setIsOpen(true)}
        className={clsx(
          "input w-full flex justify-between items-center rounded-lg text-left h-12 px-3",
          "focus:outline-none focus:ring-0 focus:ring-offset-0",
          {
            "focus:border-emerald-400 dark:focus:border-emerald-500": !difficulty,
            "focus:border-slate-300 dark:focus:border-slate-400": difficulty === "trifle",
            "focus:border-emerald-300 dark:focus:border-emerald-400": difficulty === "easy",
            "focus:border-sky-300 dark:focus:border-sky-400": difficulty === "normal",
            "focus:border-orange-300 dark:focus:border-orange-400": difficulty === "hard",
          }
        )}
      >
        <div className="flex items-center gap-2 flex-1 min-w-0">
          {selectedTag ? (
            <>
              <div className="w-5 h-5 rounded-full border border-base-300 shrink-0" style={{ backgroundColor: selectedTag.color }} />
              <span className="truncate text-sm">{selectedTag.name}</span>
            </>
          ) : (
            <span className="text-sm">– Выберите тег –</span>
          )}
        </div>
        {selectedTag && (
          <div
            role="button"
            onClick={(e) => {
              e.stopPropagation();
              onSelect(undefined);
            }}
            className="btn btn-ghost btn-xs rounded-full p-0 w-6 h-6 ml-2"
          >
            <X className="w-4 h-4" />
          </div>
        )}
        <ChevronDown className="w-5 h-5 opacity-70 shrink-0" />
      </button>

      <AnimatePresence>
        {isOpen && (
          <>
            <m.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }} className="fixed inset-0 bg-black/50 z-40" onClick={closeModal} />

            <m.div
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "100%" }}
              transition={{ type: "spring", damping: 35, stiffness: 400 }}
              className="fixed bottom-0 left-0 right-0 z-50 bg-base-100 rounded-t-2xl shadow-2xl max-h-[85vh] flex flex-col"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Заголовок */}
              <div className="flex items-center justify-between p-4 border-b border-base-300">
                <h3 className="font-medium text-lg">Теги</h3>
                <button onClick={closeModal} className="btn btn-ghost btn-sm">
                  <X className="w-5 h-5" />
                </button>
              </div>

              {/* Поиск */}
              <div className="p-4">
                <input
                  ref={createInputRef}
                  type="text"
                  placeholder="Поиск тегов..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="input input-bordered w-full h-12 text-base rounded-xl focus:outline-none"
                />
              </div>

              {/* Список тегов */}
              <div className="flex-1 overflow-y-auto px-4 pb-2">
                {filteredTags.length === 0 ? (
                  <p className="text-center text-sm text-base-500 py-8">{search ? "Теги не найдены" : "Нет тегов"}</p>
                ) : (
                  <div className="space-y-1">
                    {filteredTags.map((tag) => (
                      <TagItem
                        key={tag.id}
                        tag={tag}
                        isEditing={editingId === tag.id}
                        editName={editName}
                        editColor={editColor}
                        onSelect={handleSelect}
                        onEdit={startEdit}
                        onDelete={(id) => {
                          deleteTag(id);
                          if (selectedTag?.id === id) onSelect(undefined);
                        }}
                        onChangeEditName={setEditName}
                        onChangeEditColor={setEditColor}
                        onSaveEdit={saveEdit}
                        onCancelEdit={cancelEdit}
                      />
                    ))}
                  </div>
                )}
              </div>

              {/* Создание нового тега */}
              <div className="border-t border-base-300 bg-base-100 p-4">
                {isCreating ? (
                  <div className="space-y-3">
                    <input
                      ref={createInputRef}
                      type="text"
                      placeholder="Название тега..."
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      onKeyDown={(e) => e.key === "Enter" && handleCreate()}
                      className="input input-bordered w-full h-12 rounded-xl text-base focus:outline-none"
                      autoFocus
                    />
                    <div className="flex items-center justify-between">
                      <div className="flex gap-2 flex-wrap" onClick={(e) => e.stopPropagation()}>
                        {COLORS.map((c) => (
                          <button
                            key={c}
                            type="button"
                            onClick={() => setEditColor(c)}
                            style={{ backgroundColor: c }}
                            className={clsx(
                              "w-10 h-10 rounded-full ring-2 transition",
                              editColor === c ? "ring-base-content/80 shadow-lg scale-90" : "ring-transparent"
                            )}
                          />
                        ))}
                      </div>
                      <div className="flex gap-2">
                        <button onClick={handleCreate} disabled={!editName.trim()} className="btn btn-success btn-sm h-10">
                          <Check className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => {
                            setIsCreating(false);
                            setEditName("");
                            setEditColor(COLORS[0]);
                          }}
                          className="btn btn-error btn-sm h-10"
                        >
                          <X className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                ) : (
                  <button
                    type="button"
                    onClick={() => setIsCreating(true)}
                    className="btn btn-outline w-full h-12 flex items-center justify-center gap-2 text-base font-medium"
                  >
                    <Plus className="w-5 h-5" />
                    Добавить тег
                  </button>
                )}
              </div>
            </m.div>
          </>
        )}
      </AnimatePresence>
    </div>
  );
};
