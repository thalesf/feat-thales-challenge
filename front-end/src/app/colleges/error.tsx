"use client";

import StatusMessage from "@/components/StatusMessage/StatusMessage";
import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <StatusMessage
      title="Something went wrong"
      message="We couldn't load reviews for this college. The service may be temporarily unavailable."
      action={{ label: "Try again", onClick: reset }}
    />
  );
}
