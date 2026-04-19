import type { Metadata } from "next";
import StatusMessage from "@/components/StatusMessage/StatusMessage";

export const metadata: Metadata = {
  title: "College not found",
};

export default function NotFound() {
  return (
    <StatusMessage
      align="center"
      title="College not found"
      message="We couldn't find that college. Try searching again above."
    />
  );
}
