"use client";

import { useEffect, useRef, useState } from "react";
import Autocomplete from "@/components/Autocomplete/Autocomplete";
import { useDebouncedValue } from "@/hooks/useDebounce";
import { fetchAutocomplete } from "@/lib/api/colleges";
import type { College } from "@/lib/api/types";

const DEBOUNCE_MS = 200;

export default function Search() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<College[]>([]);
  const [searchedQuery, setSearchedQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query, DEBOUNCE_MS);
  const latestRequestId = useRef(0);

  useEffect(() => {
    const q = debouncedQuery.trim();
    if (!q) {
      setResults([]);
      setSearchedQuery("");
      return;
    }

    const requestId = latestRequestId.current;
    const controller = new AbortController();
    fetchAutocomplete(q, controller.signal)
      .then((items) => {
        if (requestId !== latestRequestId.current) return;
        setResults(items);
        setSearchedQuery(q);
      })
      .catch((err) => {
        if (err instanceof DOMException && err.name === "AbortError") return;
        console.error(err);
      });

    return () => controller.abort();
  }, [debouncedQuery]);

  const isSettled = query.trim() === searchedQuery && searchedQuery.length > 0;

  return (
    <Autocomplete
      query={query}
      onQueryChange={setQuery}
      items={results}
      isSettled={isSettled}
    />
  );
}
