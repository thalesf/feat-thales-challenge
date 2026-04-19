"use client";

import styles from "./StatusMessage.module.css";

type Props = {
  title: string;
  message: string;
  align?: "left" | "center";
  action?: {
    label: string;
    onClick: () => void;
  };
};

export default function StatusMessage({
  title,
  message,
  align = "left",
  action,
}: Props) {
  return (
    <section
      className={`${styles.statusMessage} ${align === "center" ? styles.centered : ""}`}
      role={action ? "alert" : undefined}
    >
      <h2 className={styles.title}>{title}</h2>
      <p className={`${styles.message} ${action ? styles.withAction : ""}`}>
        {message}
      </p>
      {action ? (
        <button type="button" className={styles.action} onClick={action.onClick}>
          {action.label}
        </button>
      ) : null}
    </section>
  );
}
