<script setup lang="ts">
import { SelectContent as SelectContentPrimitive, SelectPortal, SelectViewport } from 'reka-ui'
import { cn } from '@/components/brutx/shared/utils'
import { selectContentVariants } from './select-variants'

interface SelectContentProps {
    class?: string
}

const props = withDefaults(defineProps<SelectContentProps>(), {
    class: undefined,
})
</script>

<template>
    <SelectPortal>
        <SelectContentPrimitive :class="cn(selectContentVariants(), props.class)">
            <SelectViewport class="p-1">
                <slot />
            </SelectViewport>
        </SelectContentPrimitive>
    </SelectPortal>
</template>
