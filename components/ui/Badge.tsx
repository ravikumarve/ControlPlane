import { cn } from "@/lib/utils";

type BadgeVariant = "default" | "success" | "warning" | "danger" | "accent";

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  className?: string;
}

const variantStyles: Record<BadgeVariant, string> = {
  default: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300",
  success: "bg-success/10 text-success-900 dark:bg-success/20 dark:text-success",
  warning: "bg-warning/10 text-warning-900 dark:bg-warning/20 dark:text-warning",
  danger: "bg-danger/10 text-danger-900 dark:bg-danger/20 dark:text-danger",
  accent: "bg-accent/10 text-accent-900 dark:bg-accent/20 dark:text-accent",
};

export function Badge({ children, variant = "default", className }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
        variantStyles[variant],
        className
      )}
    >
      {children}
    </span>
  );
}
