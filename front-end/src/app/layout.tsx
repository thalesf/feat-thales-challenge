import type { Metadata } from "next";
import { Inter, Fraunces, JetBrains_Mono } from "next/font/google";
import Header from "@/components/Header/Header";
import styles from "./layout.module.css";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-sans-next",
});

const fraunces = Fraunces({
  subsets: ["latin"],
  axes: ["SOFT", "opsz"],
  style: ["normal", "italic"],
  display: "swap",
  variable: "--font-serif-next",
});

const jetBrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-mono-next",
});

export const metadata: Metadata = {
  title: "College Reviews",
  description: "Niche Full-Stack Coding Exercise",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${inter.variable} ${fraunces.variable} ${jetBrainsMono.variable}`}
    >
      <body>
        <main className={styles.main}>
          <Header />
          {children}
        </main>
      </body>
    </html>
  );
}
