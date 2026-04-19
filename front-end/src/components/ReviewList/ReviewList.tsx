import type { College } from "@/lib/api/types";
import styles from "./ReviewList.module.css";

type Props = { college: College; reviews: string[] };

export default function ReviewList({ college, reviews }: Props) {
  return (
    <section className={styles.reviewList} aria-labelledby="reviews-heading">
      <header className={styles.header}>
        <h2 id="reviews-heading" className={styles.title}>
          {college.name}
        </h2>
        {reviews.length > 0 && (
          <span className={styles.count}>
            {reviews.length} {reviews.length === 1 ? "review" : "reviews"}
          </span>
        )}
      </header>
      {reviews.length === 0 ? (
        <p className={styles.emptyState}>No reviews yet.</p>
      ) : (
        <ul className={styles.list}>
          {reviews.map((text, i) => (
            <li key={`${college.url}#${i}`} className={styles.item}>
              {text}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
