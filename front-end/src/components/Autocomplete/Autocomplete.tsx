"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  useEffect,
  useId,
  useRef,
  useState,
  type KeyboardEvent,
} from "react";
import type { College } from "@/lib/api/types";
import styles from "./Autocomplete.module.css";

type Props = {
  query: string;
  onQueryChange: (value: string) => void;
  items: College[];
  isSettled: boolean;
};

const hrefFor = (college: College) => `/colleges/${college.url}`;

export default function Autocomplete({
  query,
  onQueryChange,
  items,
  isSettled,
}: Props) {
  const router = useRouter();
  const [activeIndex, setActiveIndex] = useState(-1);
  const [isOpen, setIsOpen] = useState(true);
  const listboxId = useId();
  const wrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setActiveIndex(-1);
  }, [items.length]);

  useEffect(() => {
    const handlePointerDown = (event: MouseEvent) => {
      if (!wrapperRef.current) return;
      if (!wrapperRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setActiveIndex(-1);
      }
    };
    document.addEventListener("mousedown", handlePointerDown);
    return () => document.removeEventListener("mousedown", handlePointerDown);
  }, []);

  const hasItems = items.length > 0;
  const hasQuery = query.trim().length > 0;
  const showListbox = isOpen && hasQuery && hasItems;
  const showEmpty = isOpen && hasQuery && !hasItems && isSettled;

  const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    switch (event.key) {
      case "ArrowDown":
        if (!hasItems) return;
        event.preventDefault();
        setIsOpen(true);
        setActiveIndex((i) => (i + 1) % items.length);
        break;
      case "ArrowUp":
        if (!hasItems) return;
        event.preventDefault();
        setIsOpen(true);
        setActiveIndex((i) => (i <= 0 ? items.length - 1 : i - 1));
        break;
      case "Enter": {
        if (!showListbox) return;
        const index = activeIndex >= 0 ? activeIndex : 0;
        const selected = items[index];
        if (selected) {
          event.preventDefault();
          router.push(hrefFor(selected));
          setIsOpen(false);
        }
        break;
      }
      case "Escape":
        event.preventDefault();
        setIsOpen(false);
        setActiveIndex(-1);
        break;
    }
  };

  return (
    <div ref={wrapperRef} className={styles.autocomplete}>
      <input
        className={styles.input}
        type="search"
        value={query}
        placeholder="Search colleges, e.g. Harvard..."
        role="combobox"
        aria-expanded={showListbox}
        aria-controls={listboxId}
        aria-autocomplete="list"
        aria-activedescendant={
          showListbox && activeIndex >= 0
            ? `${listboxId}-option-${activeIndex}`
            : undefined
        }
        onChange={(e) => {
          setActiveIndex(-1);
          setIsOpen(true);
          onQueryChange(e.target.value);
        }}
        onFocus={() => setIsOpen(true)}
        onKeyDown={handleKeyDown}
      />

      {showListbox ? (
        <ul id={listboxId} className={styles.listbox} role="listbox">
          {items.map((college, index) => (
            <li
              key={college.url}
              id={`${listboxId}-option-${index}`}
              role="option"
              aria-selected={index === activeIndex}
              className={`${styles.option} ${
                index === activeIndex ? styles.optionActive : ""
              }`}
              onMouseEnter={() => setActiveIndex(index)}
            >
              <Link
                href={hrefFor(college)}
                className={styles.optionLink}
                onClick={() => setIsOpen(false)}
              >
                {college.name}
              </Link>
            </li>
          ))}
        </ul>
      ) : showEmpty ? (
        <div className={styles.emptyState} role="status">
          No matches for &ldquo;{query.trim()}&rdquo;
        </div>
      ) : null}
    </div>
  );
}
