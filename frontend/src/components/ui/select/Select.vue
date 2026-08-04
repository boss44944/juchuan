<script setup lang="ts">
import { computed, useSlots } from 'vue'
import { SelectRoot, SelectValue } from 'reka-ui'
import SelectTrigger from './SelectTrigger.vue'
import SelectContent from './SelectContent.vue'
import SelectItem from './SelectItem.vue'
import { cn } from '@/components/brutx/shared/utils'

export interface SelectOption {
    label: string
    value: string
    disabled?: boolean
    [key: string]: unknown
}

interface SelectProps {
    options?: SelectOption[]
    placeholder?: string
    disabled?: boolean
    required?: boolean
    name?: string
    id?: string
    size?: 'sm' | 'default' | 'lg'
    variant?: 'default' | 'error' | 'success'
    class?: string
    contentClass?: string
}

const props = withDefaults(defineProps<SelectProps>(), {
    options: () => [],
    placeholder: 'Select an option',
    disabled: false,
    required: false,
    name: undefined,
    id: undefined,
    size: 'default',
    variant: 'default',
    class: undefined,
    contentClass: undefined,
})

const modelValue = defineModel<string>()
const slots = useSlots()
const hasDefaultSlot = computed(() => !!slots.default)
</script>

<template>
    <SelectRoot v-model="modelValue" :name="name" :disabled="disabled" :required="required">
        <!-- 完全自定义模式：使用默认插槽接管全部内容 -->
        <slot v-if="hasDefaultSlot" />
        <template v-else>
            <SelectTrigger
                :id="id"
                :size="size"
                :variant="variant"
                :disabled="disabled"
                :class="cn(props.class)"
            >
                <SelectValue :placeholder="placeholder" />
            </SelectTrigger>
            <SelectContent :class="contentClass">
                <SelectItem
                    v-for="opt in options"
                    :key="opt.value"
                    :value="opt.value"
                    :disabled="opt.disabled"
                >
                    {{ opt.label }}
                </SelectItem>
            </SelectContent>
        </template>
    </SelectRoot>
</template>
