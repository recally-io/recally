export async function sha256Hex(data: string | ArrayBuffer | Uint8Array): Promise<string> {
  const bytes =
    typeof data === "string"
      ? new TextEncoder().encode(data)
      : data instanceof Uint8Array
        ? data
        : new Uint8Array(data);

  const digest = await crypto.subtle.digest("SHA-256", bytes as unknown as ArrayBuffer);

  return [...new Uint8Array(digest)].map((b) => b.toString(16).padStart(2, "0")).join("");
}
