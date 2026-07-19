import { cn } from "@/lib/utils";
import { HTMLAttributes, forwardRef } from "react";

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "elevated" | "bordered";
}

const variantStyles = {
  default: "bg-white dark:bg-surface-card",
  elevated:
    "bg-white dark:bg-surface-card shadow-lg shadow-black/5 dark:shadow-black/20",
  bordered:
    "bg-white dark:bg-surface-card border border-surface-border",
};

export const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = "default", ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(
          "rounded-xl p-6",
          variantStyles[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Card.displayName = "Card";
