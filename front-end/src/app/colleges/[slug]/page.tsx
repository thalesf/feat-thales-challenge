import type { Metadata } from "next";
import { notFound } from "next/navigation";
import ReviewList from "@/components/ReviewList/ReviewList";
import { fetchReviews, NotFoundError } from "@/lib/api/colleges";

type Params = Promise<{ slug: string }>;

export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  const { slug } = await params;
  try {
    const { college } = await fetchReviews(slug);
    return { title: `${college.name} — College Reviews` };
  } catch {
    return {};
  }
}

export default async function CollegePage({ params }: { params: Params }) {
  const { slug } = await params;
  try {
    const { college, reviews } = await fetchReviews(slug);
    return <ReviewList college={college} reviews={reviews} />;
  } catch (err) {
    if (err instanceof NotFoundError) notFound();
    throw err;
  }
}
