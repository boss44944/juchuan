import { cva } from 'class-variance-authority'
import { brutalHoverLift, brutalPress, brutalHighlightLift } from '@/components/brutx/shared/lib/brutal-interaction-variants'

export const selectTriggerVariants = cva(
    [
        'flex w-full items-center justify-between gap-2 px-4 py-2',
        'bg-brutal-bg text-brutal-fg border-3 border-brutal rounded-brutal',
        'font-medium placeholder:text-brutal-placeholder',
        'shadow-brutal transition-all duration-150',
        brutalHoverLift,
        brutalPress,
        'focus:outline-none focus:shadow-brutal-lg focus:-translate-x-0.5 focus:-translate-y-0.5',
        'disabled:pointer-events-none disabled:opacity-50',
        '[&>span]:line-clamp-1',
    ],
    {
        variants: {
            size: {
                sm: 'h-9 text-sm',
                default: 'h-11 text-base',
                lg: 'h-14 text-lg',
            },
            variant: {
                default: '',
                error: 'border-brutal-destructive focus:ring-brutal-destructive',
                success: 'border-brutal-success focus:ring-brutal-success',
            },
        },
        defaultVariants: {
            size: 'default',
            variant: 'default',
        },
    }
)

export const selectContentVariants = cva(
    [
        'relative z-50 max-h-96 min-w-[8rem] overflow-hidden',
        'bg-brutal-bg text-brutal-fg border-3 border-brutal shadow-brutal rounded-brutal',
    ]
)

export const selectItemVariants = cva(
    [
        'relative flex w-full cursor-pointer select-none items-center py-2 pl-8 pr-3',
        'font-medium outline-none',
        brutalHighlightLift,
        brutalPress,
        'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
    ],
    {
        variants: {
            variant: {
                default: 'focus:bg-brutal-accent focus:text-brutal-accent-foreground',
                primary: 'focus:bg-brutal-primary focus:text-brutal-primary-foreground',
                secondary: 'focus:bg-brutal-secondary focus:text-brutal-secondary-foreground',
            },
        },
        defaultVariants: {
            variant: 'default',
        },
    }
)
