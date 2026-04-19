import type { College, ReviewsResponse } from "./types";

export class NotFoundError extends Error {
  constructor(message = "Not found") {
    super(message);
    this.name = "NotFoundError";
  }
}

const CLIENT_BASE =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

const SERVER_BASE =
  process.env.API_BASE ?? CLIENT_BASE;

export async function fetchAutocomplete(q: string, signal?: AbortSignal): Promise<College[]> {
  const res = await fetch(
    `${CLIENT_BASE}/autocomplete?q=${encodeURIComponent(q)}`,
    { signal },
  );
  if (!res.ok) throw new Error(`autocomplete failed: ${res.status}`)
  return res.json()
}

export async function fetchReviews(url: string): Promise<ReviewsResponse>  {
  const res = await fetch(`${SERVER_BASE}/reviews?url=${encodeURIComponent(url)}`, {
    next: { revalidate: 60 },
  });
  if (res.status === 404) throw new NotFoundError(`college not found: ${url}`);
  if (!res.ok) throw new Error(`reviews failed: ${res.status}`);
  return res.json()
}
