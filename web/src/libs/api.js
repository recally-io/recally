import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient();

export async function request(path, method = "GET", params = {}, body = {}) {
  let options = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
  };

  if (!!params) {
    options["query"] = params;
  }
  if (!!body) {
    options["body"] = JSON.stringify(body);
  }

  const res = await fetch(path, options);
  if (!res.ok) {
    throw new Error(
      `Failed to fetch ${path}, response status: ${res.status} ${res.statusText}`,
    );
  }
  return res.json();
}
