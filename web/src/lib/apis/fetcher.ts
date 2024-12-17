type FetcherError = Error & { status?: number; info?: unknown };

const fetcher = async <T>(url: string, init?: RequestInit): Promise<T> => {
  const response = await fetch(url, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  if (!response.ok) {
    const error = new Error(
      "An error occurred while fetching the data.",
    ) as FetcherError;
    error.status = response.status;

    try {
      error.info = await response.json();
    } catch {
      error.info = await response.text();
    }

    throw error;
  }
  const data = await response.json();
  return data.data;
};

export default fetcher;
