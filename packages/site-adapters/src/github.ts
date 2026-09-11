import { z } from "zod";
import type { AdapterFetch, AdapterResult, SiteAdapter } from "./types";

// github.com/{owner}/{repo} → repo record via the public REST API.
const inputSchema = z.object({});

export const githubRepoAdapter: SiteAdapter = {
  descriptor: {
    id: "github/repo",
    domains: ["github.com"],
    summary: "GitHub repository record: metadata, description, README",
    produces: "repo json + readme markdown",
  },
  match(url: URL): boolean {
    if (url.hostname !== "github.com") return false;
    const parts = url.pathname.split("/").filter(Boolean);
    // owner/repo only — issues/pulls/actions are different record types.
    return parts.length === 2 && !parts[1]!.includes(".");
  },
  inputSchema: () => inputSchema,
  async run(input, _ctx, fetchJson: AdapterFetch): Promise<AdapterResult> {
    const parts = new URL(input.url).pathname.split("/").filter(Boolean);
    const [owner, repo] = parts as [string, string];
    const repoUrl = `https://api.github.com/repos/${owner}/${repo}`;
    const readmeUrl = `https://api.github.com/repos/${owner}/${repo}/readme`;
    const repoData = (await fetchJson(repoUrl)) as Record<string, unknown>;

    let readme: string | null = null;
    try {
      const readmeData = (await fetchJson(
        `${readmeUrl}?${new URLSearchParams({ ref: "HEAD" })}`,
      )) as {
        content?: string;
        encoding?: string;
      };
      if (readmeData.content && readmeData.encoding === "base64") {
        readme = new TextDecoder().decode(
          Uint8Array.from(atob(readmeData.content.replace(/\s/g, "")), (c) => c.charCodeAt(0)),
        );
      }
    } catch {
      readme = null; // missing README is a gap, not a failure
    }

    return {
      record: {
        kind: "github_repo",
        full_name: repoData.full_name,
        description: repoData.description ?? null,
        homepage: repoData.homepage ?? null,
        language: repoData.language ?? null,
        stars: repoData.stargazers_count ?? null,
        topics: repoData.topics ?? [],
        license: (repoData.license as { spdx_id?: string } | null)?.spdx_id ?? null,
        default_branch: repoData.default_branch ?? null,
        created_at: repoData.created_at ?? null,
        pushed_at: repoData.pushed_at ?? null,
        readme_markdown: readme,
      },
      provenance: {
        endpointUrls: [repoUrl, readmeUrl],
        fetchedAt: new Date().toISOString(),
      },
      mediaRefs: (repoData.owner as { avatar_url?: string } | undefined)?.avatar_url
        ? [
            {
              url: (repoData.owner as { avatar_url: string }).avatar_url,
              kind: "image",
              role: "optional",
            },
          ]
        : [],
      canonicalUrl: (repoData.html_url as string | undefined) ?? input.url,
    };
  },
};
