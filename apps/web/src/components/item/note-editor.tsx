import { useEffect, useState } from "react";
import { api, type Note } from "../../api";

export function NoteEditor({
  itemId,
  notes,
  onSaved,
  flash,
}: {
  itemId: string;
  notes: Note[];
  onSaved: () => Promise<void>;
  flash: (m: string) => void;
}) {
  const note = notes[0];
  const [body, setBody] = useState(note?.body ?? "");
  const [dirty, setDirty] = useState(false);
  useEffect(() => {
    if (!dirty) setBody(note?.body ?? "");
  }, [note?.body, dirty]);

  const save = async () => {
    if (note) await api.patchNote(note.id, body, note.version);
    else if (body.trim()) await api.addNote(itemId, body);
    else return;
    setDirty(false);
    flash("note saved");
    await onSaved();
  };

  return (
    <div>
      <textarea
        className="w-full rounded-lg border border-line bg-surface p-3 text-[12.5px] leading-relaxed text-ink-2 outline-none focus:border-accent"
        rows={3}
        placeholder="Why did you save this?"
        value={body}
        onChange={(e) => {
          setBody(e.target.value);
          setDirty(true);
        }}
      />
      {dirty && (
        <button
          type="button"
          className="mt-1.5 rounded-md bg-accent px-3 py-1 text-xs font-semibold text-white"
          onClick={() => void save()}
        >
          save note
        </button>
      )}
    </div>
  );
}
