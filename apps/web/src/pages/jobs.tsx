import { useState } from "react";
import { api } from "../api";

// Minimal job inspector — paste a job id, see status/result. The full job
// list + controls land with the M1 UI work.
export function JobsPage() {
  const [id, setId] = useState("");
  const [job, setJob] = useState<{
    status: string;
    result: string | null;
    error: string | null;
  } | null>(null);
  const [error, setError] = useState<string | null>(null);

  return (
    <div>
      <form
        className="mb-4 flex gap-2"
        onSubmit={(e) => {
          e.preventDefault();
          setError(null);
          api
            .getJob(id)
            .then(setJob)
            .catch((e) => setError(e.message));
        }}
      >
        <input
          className="flex-1 rounded border px-3 py-2 text-sm"
          placeholder="job id"
          value={id}
          onChange={(e) => setId(e.target.value)}
        />
        <button type="submit" className="rounded bg-black px-4 py-2 text-sm text-white">
          check
        </button>
      </form>
      {error && <p className="text-sm text-red-600">{error}</p>}
      {job && (
        <pre className="rounded bg-neutral-100 p-3 text-xs">{JSON.stringify(job, null, 2)}</pre>
      )}
    </div>
  );
}
