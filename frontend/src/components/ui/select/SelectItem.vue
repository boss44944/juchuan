<script setup lang="ts">
import { computed } from 'vue'
import { type VariantProps } from 'class-variance-authority'
import {
    SelectItem as SelectItemPrimitive,
    SelectItemIndicator,
    SelectItemText,
} from 'reka-ui'
import { Check } from '@lucide/vue'
import { cn } from '@/components/brutx/shared/utils'
import { selectItemVariants } from './select-variants'

type ItemVariantProps = VariantProps<typeof selectItemVariants>

interface SelectItemProps {
    value: string
    disabled?: boolean
    variant?: NonNullable<ItemVariantProps['variant']>
    class?: string
}

const props = withDefaults(defineProps<SelectItemProps>(), {
    disabled: false,
    variant: 'default',
    class: undefined,
})

const classes = computed(() =>
    cn(selectItemVariants({ variant: props.variant }), props.class)
)
</script>

<template>
    <SelectItemPrimitive :value="value" :disabled="disabled" :class="classes">
        <span class="absolute left-2 flex h-4 w-4 items-center justify-center">
            <SelectItemIndicator>
                <Check class="h-4 w-4 stroke-[3]" />
            </SelectItemIndicator>
        </span>
        <SelectItemText><slot /></SelectItemText>
    </SelectItemPrimitive>
</template>
