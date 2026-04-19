import Search from "@/components/Search";
import styles from "./Header.module.css";

export default function Header() {
  return (
    <header className={styles.header}>
      <p className={styles.eyebrow}>Vol. 01 — Student Reviews from Current and Former Students</p>
      <h1 className={styles.title}>College Reviews</h1>
      <p className={styles.subtitle}>
        Search for a college to see real student reviews.
      </p>
      <Search />
    </header>
  );
}
