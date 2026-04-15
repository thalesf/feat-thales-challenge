import type { Metadata } from "next";
import "./globals.css";

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
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
