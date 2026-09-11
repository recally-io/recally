import { z } from "zod";
import type { AdapterFetch, AdapterResult, SiteAdapter } from "./types";

// news.ycombinator.com/item?id=N → item + comment tree via the Algolia HN API.
const inputSchema = z.object({
  max_comments: z.number().int().positive().max(500).default(200),
});

export const hackernewsItemAdapter: SiteAdapter = {
  descriptor: {
    id: "hackernews/item",
    domains: ["news.ycombinator.com"],
    summary: "Hacker News item record: story + comment tree",
    produces: "item json + comments",
  },
  match(url: URL): boolean {
    return (
      url.hostname === "news.ycombinator.com" &&
      url.pathname === "/item" &&
      url.searchParams.has("id")
    );
  },
  inputSchema: () => inputSchema,
  async run(input, _ctx, fetchJson: AdapterFetch): Promise<AdapterResult> {
    const args = inputSchema.parse(input.args);
    const id = new URL(input.url).searchParams.get("id");
    const endpoint = `https://hn.algolia.com/api/v1/items/${id}`;
    const data = (await fetchJson(endpoint)) as {
      id: number;
      title?: string;
      url?: string;
      text?: string;
      author?: string;
      points?: number;
      created_at?: string;
      type?: string;
      children?: unknown[];
    };

    let comments = 0;
    const countComments = (nodes: unknown[] | undefined) => {
      for (const n of nodes ?? []) {
        comments++;
        countComments((n as { children?: unknown[] }).children);
      }
    };
    countComments(data.children);
    const truncated = comments > args.max_comments;

    return {
      record: {
        kind: "hn_item",
        id: data.id,
        type: data.type ?? null,
        title: data.title ?? null,
        url: data.url ?? null,
        text: data.text ?? null,
        author: data.author ?? null,
        points: data.points ?? null,
        created_at: data.created_at ?? null,
        comment_count_fetched: Math.min(comments, args.max_comments),
        comments_truncated: truncated,
        children: data.children ?? [],
      },
      provenance: {
        endpointUrls: [endpoint],
        fetchedAt: new Date().toISOString(),
      },
      mediaRefs: [],
      canonicalUrl: `https://news.ycombinator.com/item?id=${data.id}`,
    };
  },
};
