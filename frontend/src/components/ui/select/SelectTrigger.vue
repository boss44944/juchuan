<script setup lang="ts">
import { computed } from 'vue'
import { type VariantProps } from 'class-variance-authority'
import { SelectTrigger as SelectTriggerPrimitive, SelectIcon } from 'reka-ui'
import { ChevronDown } from '@lucide/vue'
import { cn } from '@/components/brutx/shared/utils'
import { selectTriggerVariants } from './select-variants'
import { iconSizeVariants } from '@/components/brutx/shared/lib/icon-size-variants'

type TriggerVariantProps = VariantProps<typeof selectTriggerVariants>

interface SelectTriggerProps {
    size?: NonNullable<TriggerVariantProps['size']>
    variant?: NonNullable<TriggerVariantProps['variant']>
    disabled?: boolean
    id?: string
    class?: string
}

const props = withDefaults(defineProps<SelectTriggerProps>(), {
    size: 'default',
    variant: 'default',
    disabled: false,
    id: undefined,
    class: undefined,
})

const SIZE_TO_ICON: Record<NonNullable<TriggerVariantProps['size']>, 'sm' | 'default' | 'lg'> = {
    sm: 'sm',
    default: 'default',
    lg: 'lg',
}

const classes = computed(() =>
    cn(selectTriggerVariants({ size: props.size, variant: props.variant }), props.class)
)

const iconClasses = computed(() =>
    cn(iconSizeVariants({ size: SIZE_TO_ICON[props.size ?? 'default'] }))
)
</script>

<template>
    <SelectTriggerPrimitive :id="id" :class="classes" :disabled="disabled" aria-haspopup="listbox">
        <slot />
        <SelectIcon as-child>
            <ChevronDown :class="iconClasses" />
        </SelectIcon>
    </SelectTriggerPrimitive>
</template>
