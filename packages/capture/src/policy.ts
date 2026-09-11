import { AppError } from "@recally/domain";

// Network policy per plan §14.2. Applied to the initial URL, every redirect
// hop, and every asset request — not just the first target.
//
// DNS checks alone are not network isolation; where a resolved address cannot
// be verified the policy denies rather than widens (plan §14.2).

const BLOCKED_HOSTNAMES = new Set(["localhost", "metadata.google.internal", "instance-data"]);

const BLOCKED_SUFFIXES = [".localhost", ".internal", ".localdomain", ".home.arpa"];

function isPrivateIPv4(ip: string): boolean {
  const parts = ip.split(".").map(Number);
  if (parts.length !== 4 || parts.some((n) => Number.isNaN(n) || n < 0 || n > 255)) return false;
  const [a, b] = parts as [number, number, number, number];
  return (
    a === 10 ||
    a === 127 ||
    (a === 172 && b >= 16 && b <= 31) ||
    (a === 192 && b === 168) ||
    (a === 169 && b === 254) ||
    a === 0 ||
    (a === 100 && b >= 64 && b <= 127) || // CGNAT
    a >= 224 // multicast/reserved/broadcast
  );
}

function isBlockedIPv6(ip: string): boolean {
  const normalized = ip.toLowerCase().replace(/^\[|\]$/g, "");
  return (
    normalized === "::1" ||
    normalized === "::" ||
    normalized.startsWith("fe80") || // link-local
    normalized.startsWith("fc") ||
    normalized.startsWith("fd") || // ULA
    (normalized.startsWith("::ffff:") && isPrivateIPv4(normalized.slice(7))) // v4-mapped
  );
}

const IPV4_RE = /^\d{1,3}(\.\d{1,3}){3}$/;

export function checkUrlTarget(rawUrl: string): URL {
  let url: URL;
  try {
    url = new URL(rawUrl);
  } catch {
    throw new AppError("invalid_input", `unparseable url: ${rawUrl.slice(0, 200)}`);
  }

  if (url.protocol !== "https:" && url.protocol !== "http:") {
    throw new AppError("policy_denied", `scheme ${url.protocol} not allowed`);
  }
  if (url.username || url.password) {
    throw new AppError("policy_denied", "url userinfo not allowed");
  }
  const port = url.port ? Number(url.port) : url.protocol === "https:" ? 443 : 80;
  if (![80, 443, 8080, 8443].includes(port)) {
    throw new AppError("policy_denied", `port ${port} not allowed`);
  }

  const host = url.hostname.toLowerCase();
  if (BLOCKED_HOSTNAMES.has(host) || BLOCKED_SUFFIXES.some((s) => host.endsWith(s))) {
    throw new AppError("policy_denied", `host ${host} not allowed`);
  }
  if (IPV4_RE.test(host) && isPrivateIPv4(host)) {
    throw new AppError("policy_denied", `private address ${host} not allowed`);
  }
  if (host.includes(":") && isBlockedIPv6(host)) {
    throw new AppError("policy_denied", `private address ${host} not allowed`);
  }
  return url;
}

// Actions the browser stage may take must serve "read this content" only.
// Login, posting, purchases, subscriptions, messages, and form fills are never
// approved regardless of what the model asks for (plan §5.6).
const ALLOWED_BROWSER_ACTIONS = new Set([
  "observe",
  "wait_for",
  "scroll",
  "expand",
  "navigate_same_article",
  "capture",
]);

export function checkBrowserAction(action: string): void {
  if (!ALLOWED_BROWSER_ACTIONS.has(action)) {
    throw new AppError("policy_denied", `browser action ${action} not allowed`);
  }
}

export const VALID_BROWSER_UPGRADE_REASONS = new Set([
  "js_shell",
  "summary_only",
  "expand_control",
  "lazy_load_gap",
  "multi_page",
  "content_contradiction",
]);

export function checkBrowserUpgrade(reasonCode: string, evidenceIds: string[]): void {
  if (!VALID_BROWSER_UPGRADE_REASONS.has(reasonCode)) {
    throw new AppError(
      "policy_denied",
      `browser upgrade reason ${reasonCode} is not evidence-backed`,
    );
  }
  if (evidenceIds.length === 0) {
    throw new AppError("policy_denied", "browser upgrade requires cited observations");
  }
}
