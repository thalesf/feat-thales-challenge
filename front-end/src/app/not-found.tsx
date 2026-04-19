import StatusMessage from "@/components/StatusMessage/StatusMessage";

export default function NotFound() {
  return (
    <StatusMessage
      align="center"
      title="No results found"
      message="We couldn't find any college matching your search. Try a different name."
    />
  );
}
