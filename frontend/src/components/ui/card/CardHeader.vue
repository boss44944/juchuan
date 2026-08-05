<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/components/brutx/shared/utils'

interface CardHeaderProps {
    class?: string
}

const props = defineProps<CardHeaderProps>()

const classes = computed(() =>
    cn('flex flex-col space-y-1.5 pb-4', props.class)
)
</script>

<template>
    <div :class="classes">
        <slot />
    </div>
</template>
