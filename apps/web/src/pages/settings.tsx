import { useEffect, useState } from "react";
import { type ApiToken, api, type Me } from "../api";
import { timeAgo } from "../components/bits";

const ALL_SCOPES = [
  { id: "items:read", label: "items:read", hint: "read library + content" },
  { id: "items:write", label: "items:write", hint: "save urls, re-capture" },
  { id: "search:read", label: "search:read", hint: "search the archive" },
  { id: "notes:write", label: "notes:write", hint: "write notes" },
  { id: "shares:write", label: "shares:write", hint: "create share links" },
] as const;

export function SettingsPage() {
  const [me, setMe] = useState<Me | null>(null);
  const [tokens, setTokens] = useState<ApiToken[]>([]);
  const [name, setName] = useState("");

  const [scopes, setScopes] = useState<Set<string>>(
    new Set(["items:read", "items:write", "search:read"]),
  );

  const [fresh, setFresh] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = () => {
    api
      .me()
      .then(setMe)
      .catch((e) => setError(e.message));
    api
      .listTokens()
      .then((r) => setTokens(r.tokens))
      .catch(() => {});
  };

  useEffect(load, []);

  const create = async () => {
    if (!name.trim() || scopes.size === 0) return;

    try {
      const r = await api.createToken(name.trim(), [...scopes], null);
      setFresh(r.token);
      setName("");
      api.listTokens().then((r2) => setTokens(r2.tokens));
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  return (
    <div className="mx-auto max-w-2xl px-6 py-8">
      <h1 className="mb-6 font-serif text-2xl">Settings</h1>
      {error && <p className="mb-4 text-sm text-danger">{error}</p>}

      <h2 className="mb-2 font-mono text-[11px] font-semibold tracking-[0.12em] text-ink-3 uppercase">
        Identity
      </h2>
      <div className="mb-8 rounded-xl border border-line bg-surface p-4 text-[13px]">
        {me ? (
          <dl className="grid grid-cols-[110px_1fr] gap-y-1.5">
            <dt className="text-ink-3">library</dt>
            <dd>{me.library_name}</dd>
            <dt className="text-ink-3">signed in as</dt>
            <dd>{me.email ?? me.actor}</dd>
            <dt className="text-ink-3">auth via</dt>
            <dd className="font-mono text-xs">{me.via}</dd>
          </dl>
        ) : (
          <p className="text-ink-3">loading…</p>
        )}
      </div>

      <h2 className="mb-2 font-mono text-[11px] font-semibold tracking-[0.12em] text-ink-3 uppercase">
        API keys — for agents
      </h2>
      <p className="mb-3 text-[12.5px] text-ink-2">
        Create a key, then point your agent at{" "}
        <a href="/skill.md" className="text-accent underline decoration-line underline-offset-2">
          /skill.md
        </a>{" "}
        — it describes the endpoints and auth.
      </p>

      <div className="mb-4 rounded-xl border border-line bg-surface p-4">
        <input
          className="mb-3 w-full rounded-md border border-line bg-paper px-3 py-1.5 text-[13px] outline-none focus:border-accent"
          placeholder="Key name (e.g. claude-code)"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <div className="mb-3 flex flex-wrap gap-2">
          {ALL_SCOPES.map((s) => (
            <button
              key={s.id}
              type="button"
              title={s.hint}
              onClick={() => {
                const next = new Set(scopes);

                if (next.has(s.id)) next.delete(s.id);
                else next.add(s.id);
                setScopes(next);
              }}
              className={`rounded-full border px-3 py-1 font-mono text-[11px] ${
                scopes.has(s.id)
                  ? "border-accent bg-accent-soft font-semibold text-accent"
                  : "border-line text-ink-3 hover:border-ink-3"
              }`}
            >
              {s.label}
            </button>
          ))}
        </div>
        <button
          type="button"
          className="rounded-lg bg-ink px-4 py-2 text-[13px] font-semibold text-white disabled:opacity-40"
          disabled={!name.trim() || scopes.size === 0}
          onClick={() => void create()}
        >
          Create key
        </button>
        {fresh && (
          <div className="mt-4 rounded-lg border border-accent bg-accent-soft p-3">
            <p className="mb-1 text-[11px] font-semibold text-accent uppercase">
              copy now — shown once
            </p>
            <code className="block overflow-x-auto font-mono text-[12px] break-all select-all">
              {fresh}
            </code>
          </div>
        )}
      </div>

      <ul className="divide-y divide-line border-y border-line">
        {tokens.map((t) => (
          <li key={t.id} className="flex items-center gap-3 py-3">
            <div className="min-w-0 flex-1">
              <p className="text-[13px] font-medium">{t.name}</p>
              <p className="font-mono text-[11px] text-ink-3">
                {(JSON.parse(t.scopes) as string[]).join(" ")}
              </p>
            </div>
            <span className="text-[11px] text-ink-3">
              {t.last_used_at ? `used ${timeAgo(t.last_used_at)}` : "never used"}
            </span>
            <button
              type="button"
              className="rounded-md border border-line px-2.5 py-1 text-[11.5px] text-ink-2 hover:border-danger hover:text-danger"
              onClick={() =>
                void api
                  .revokeToken(t.id)
                  .then(() => api.listTokens())
                  .then((r) => setTokens(r.tokens))
              }
            >
              revoke
            </button>
          </li>
        ))}
        {tokens.length === 0 && (
          <li className="py-8 text-center text-sm text-ink-3">No API keys yet.</li>
        )}
      </ul>
    </div>
  );
}
